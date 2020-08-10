package bridge

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

	goproxy "github.com/duality-solutions/web-bridge/goproxy"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/inconshreveable/go-vhost"
)

const (
	// StartHTTPPortNumber is the HTTP listening port for the first bridge link
	StartHTTPPortNumber uint16 = 8889
)

// StartBridgeNetwork listens to a port for http traffic and routes it through a link's WebRTC channel
func (b *Bridge) StartBridgeNetwork(reader io.Reader, writer io.Writer) {
	util.Info.Println("StartBridgeNetwork", b.LinkParticipants(), "http port", b.ListenPort(), "https port", b.ListenPort()+1)
	httpAddr := ":" + strconv.Itoa(int(b.ListenPort()))
	httpsAddr := ":" + strconv.Itoa(int(b.ListenPort()+1))
	proxy := goproxy.NewProxyHTTPServer()
	proxy.Verbose = true
	proxy.DataChannelReader = reader
	proxy.DataChannelWriter = writer
	proxy.BridgeID = b.LinkID()
	proxy.BridgeLinkNames = b.LinkParticipants()
	testMessage := []byte("init web-bridge")
	n, err := proxy.DataChannelWriter.Write(testMessage)
	if err != nil {
		util.Error.Println("StartBridgeNetwork", proxy.BridgeLinkNames, "write test message failed.", err)
	}
	if proxy.Verbose {
		log.Printf("Server starting up! - configured to listen on http interface %s and https interface %s", httpAddr, httpsAddr)
	}
	util.Info.Println("StartBridgeNetwork", proxy.BridgeLinkNames, "sent test message with size", n)

	go func() {
		log.Fatalln(http.ListenAndServe(httpAddr, proxy))
	}()

	// listen to the TLS ClientHello but make it a CONNECT request instead
	ln, err := net.Listen("tcp", httpsAddr)
	if err != nil {
		log.Fatalf("Error listening for https connections - %v", err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting new connection - %v", err)
			continue
		}
		go func(c net.Conn) {
			tlsConn, err := vhost.TLS(c)
			if err != nil {
				log.Printf("Error accepting new connection - %v", err)
			}
			if tlsConn.Host() == "" {
				log.Printf("Cannot support non-SNI enabled clients")
				return
			}
			connectReq := &http.Request{
				Method: "CONNECT",
				URL: &url.URL{
					Opaque: tlsConn.Host(),
					Host:   net.JoinHostPort(tlsConn.Host(), "443"),
				},
				Host:       tlsConn.Host(),
				Header:     make(http.Header),
				RemoteAddr: c.RemoteAddr().String(),
			}
			resp := dumbResponseWriter{tlsConn}
			proxy.ServeHTTP(resp, connectReq)
		}(c)
	}
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

type dumbResponseWriter struct {
	net.Conn
}

func (dumb dumbResponseWriter) Header() http.Header {
	panic("Header() should not be called on this ResponseWriter")
}

func (dumb dumbResponseWriter) Write(buf []byte) (int, error) {
	if bytes.Equal(buf, []byte("HTTP/1.0 200 OK\r\n\r\n")) {
		return len(buf), nil // throw away the HTTP OK response from the faux CONNECT request
	}
	return dumb.Write(buf)
}

func (dumb dumbResponseWriter) WriteHeader(code int) {
	panic("WriteHeader() should not be called on this ResponseWriter")
}

func (dumb dumbResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return dumb, bufio.NewReadWriter(bufio.NewReader(dumb), bufio.NewWriter(dumb)), nil
}

// copied/converted from https.go
func dial(proxy *goproxy.ProxyHTTPServer, network, addr string) (c net.Conn, err error) {
	if proxy.Tr.Dial != nil {
		return proxy.Tr.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

// copied/converted from https.go
func connectDial(proxy *goproxy.ProxyHTTPServer, network, addr string) (c net.Conn, err error) {
	if proxy.ConnectDial == nil {
		return dial(proxy, network, addr)
	}
	return proxy.ConnectDial(network, addr)
}
