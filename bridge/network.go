// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

package bridge

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

const (
	// StartHTTPPortNumber is the HTTP listening port for the first bridge link
	StartHTTPPortNumber = 8889
)

// StartBridgeNetwork listens to a port for http traffic and routes it through a link's WebRTC channel
func (l *Bridge) StartBridgeNetwork() {
	util.Info.Println("StartBridgeNetwork", l.LinkParticipants(), "port", l.ListenPort())
	server := &http.Server{
		Addr: ":" + strconv.Itoa(int(l.ListenPort())),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				l.handleTunnel(w, r)
			} else {
				http.Error(w, "HTTP not supported", http.StatusNotImplemented)
				return
			}
		}),
		ConnState: l.onConnStateEvent,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	util.Info.Println("StartBridgeNetwork before ListenAndServe", l.LinkParticipants())
	l.HTTPServer = server
	go server.ListenAndServe()
}

func transferCloser(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()
	io.Copy(dest, src)
}

// handleTunnel handles link bridge tunnel connection
func (l *Bridge) handleTunnel(w http.ResponseWriter, r *http.Request) {
	util.Info.Println("handleTunnel", l.LinkParticipants(), r.Host)
	byteRequest, err := httputil.DumpRequest(r, true)
	if err != nil {
		util.Error.Println("handleTunnel DumpRequest error", l.LinkParticipants(), r.Host, err.Error())
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	util.Info.Println("handleTunnel", l.LinkParticipants(), "byteRequest len", len(byteRequest))
	reqReader := bytes.NewReader(byteRequest)
	reqCloser := ioutil.NopCloser(reqReader)
	// todo: wrap send with standard envelop so receive knows if it is a request or response
	err = l.DataChannel.Send(byteRequest)
	if err != nil {
		util.Error.Println("handleTunnel Send error", l.LinkParticipants(), r.Host, err.Error())
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	// todo: get response from WebRTC messages
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		util.Error.Println("Hijacking not supported", l.LinkParticipants(), r.Host)
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		util.Error.Println("handleTunnel Hijack error", l.LinkParticipants(), r.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transferCloser(clientConn, reqCloser)
}

// StopBridgeNetwork stops listening to port p for http traffic and routes it through a link
func (l *Bridge) StopBridgeNetwork() error {
	return l.HTTPServer.Shutdown(nil)
}

func (l *Bridge) onConnStateEvent(conn net.Conn, state http.ConnState) {
	util.Info.Println("onChangeConnState", l.LinkParticipants(), "state", state.String())
}
