// https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

package bridge

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
)

const (
	// StartHTTPPortNumber is the HTTP listening port for the first bridge link
	StartHTTPPortNumber = 8889
)

func transferCloser(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()
	io.Copy(dest, src)
}

// StartBridgeNetwork listens to a port for http traffic and routes it through a link's WebRTC channel
func (l *Bridge) StartBridgeNetwork() {
	util.Info.Println("StartBridgeNetwork", l.LinkParticipants(), "port", l.ListenPort())
	server := &http.Server{
		Addr: ":" + strconv.Itoa(int(l.ListenPort())),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.handleTunnel(w, r)
		}),
		ConnState: l.onConnStateEvent,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	util.Info.Println("StartBridgeNetwork before ListenAndServe", l.LinkParticipants())
	l.HTTPServer = server
	server.ListenAndServe()
	util.Info.Println("StartBridgeNetwork after ListenAndServe IdleTimeout", l.HTTPServer.IdleTimeout)
}

// handleTunnel handles link bridge tunnel connection
func (l *Bridge) handleTunnel(w http.ResponseWriter, r *http.Request) {
	util.Info.Println("handleTunnel", l.LinkParticipants(), r.Host)
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		util.Info.Println("handleTunnel Hijack error", l.LinkParticipants(), r.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transferCloser(destConn, l.ReadWriteCloser)
	go transferCloser(l.ReadWriteCloser, clientConn)
}

// StopBridgeNetwork stops listening to port p for http traffic and routes it through a link
func (l *Bridge) StopBridgeNetwork() error {
	return l.HTTPServer.Shutdown(nil)
}

func (l *Bridge) onConnStateEvent(conn net.Conn, state http.ConnState) {
	util.Info.Println("onChangeConnState", l.LinkParticipants(), "state", state.String())
}
