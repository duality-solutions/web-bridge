package goproxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
	"google.golang.org/protobuf/proto"
)

type ConnectActionLiteral int

const (
	ConnectAccept = iota
	ConnectReject
	ConnectMitm
	ConnectHijack
	ConnectHTTPMitm
	ConnectProxyAuthHijack
)

var (
	OkConnect       = &ConnectAction{Action: ConnectAccept, TLSConfig: TLSConfigFromCA(&GoproxyCa)}
	MitmConnect     = &ConnectAction{Action: ConnectMitm, TLSConfig: TLSConfigFromCA(&GoproxyCa)}
	HTTPMitmConnect = &ConnectAction{Action: ConnectHTTPMitm, TLSConfig: TLSConfigFromCA(&GoproxyCa)}
	RejectConnect   = &ConnectAction{Action: ConnectReject, TLSConfig: TLSConfigFromCA(&GoproxyCa)}
	httpsRegexp     = regexp.MustCompile(`^https:\/\/`)
)

// ConnectAction enables the caller to override the standard connect flow.
// When Action is ConnectHijack, it is up to the implementer to send the
// HTTP 200, or any other valid http response back to the client from within the
// Hijack func
type ConnectAction struct {
	Action    ConnectActionLiteral
	Hijack    func(req *http.Request, client net.Conn, ctx *ProxyCtx)
	TLSConfig func(host string, ctx *ProxyCtx) (*tls.Config, error)
}

func stripPort(s string) string {
	ix := strings.IndexRune(s, ':')
	if ix == -1 {
		return s
	}
	return s[:ix]
}

func (proxy *ProxyHTTPServer) dial(network, addr string) (c net.Conn, err error) {
	if proxy.Tr.Dial != nil {
		return proxy.Tr.Dial(network, addr)
	}
	return net.Dial(network, addr)
}

func (proxy *ProxyHTTPServer) connectDial(network, addr string) (c net.Conn, err error) {
	if proxy.ConnectDial == nil {
		return proxy.dial(network, addr)
	}
	return proxy.ConnectDial(network, addr)
}

type halfClosable interface {
	net.Conn
	CloseWrite() error
	CloseRead() error
}

var _ halfClosable = (*net.TCPConn)(nil)
var counter uint64 = 0

func (proxy *ProxyHTTPServer) handleTunnel(w http.ResponseWriter, r *http.Request) {
	ctx := &ProxyCtx{Req: r, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy, certStore: proxy.CertStore}

	hij, ok := w.(http.Hijacker)
	if !ok {
		panic("httpserver does not support hijacking")
	}

	proxyClient, _, e := hij.Hijack()
	if e != nil {
		panic("Cannot hijack connection " + e.Error())
	}

	ctx.Logf("Running %d CONNECT handlers", len(proxy.httpsHandlers))
	todo, host := OkConnect, r.URL.Host
	for i, h := range proxy.httpsHandlers {
		newtodo, newhost := h.HandleConnect(host, ctx)

		// If found a result, break the loop immediately
		if newtodo != nil {
			todo, host = newtodo, newhost
			ctx.Logf("on %dth handler: %v %s", i, todo, host)
			break
		}
	}
	switch todo.Action {
	case ConnectAccept:
		if !hasPort.MatchString(host) {
			host += ":80"
		}
		targetSiteCon, err := proxy.connectDial("tcp", host)
		if err != nil {
			httpError(proxyClient, ctx, err)
			return
		}
		ctx.Logf("Accepting CONNECT to %s", host)
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

		targetTCP, targetOK := targetSiteCon.(halfClosable)
		proxyClientTCP, clientOK := proxyClient.(halfClosable)
		if targetOK && clientOK {
			go copyAndClose(ctx, targetTCP, proxyClientTCP)
			go copyAndClose(ctx, proxyClientTCP, targetTCP)
		} else {
			go func() {
				var wg sync.WaitGroup
				wg.Add(2)
				go copyOrWarn(ctx, targetSiteCon, proxyClient, &wg)
				go copyOrWarn(ctx, proxyClient, targetSiteCon, &wg)
				wg.Wait()
				proxyClient.Close()
				targetSiteCon.Close()

			}()
		}

	case ConnectHijack:
		todo.Hijack(r, proxyClient, ctx)
	case ConnectHTTPMitm:
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		ctx.Logf("Assuming CONNECT is plain HTTP tunneling, mitm proxying it for %v", proxy.BridgeLinkNames)
		targetSiteCon, err := proxy.connectDial("tcp", host)
		if err != nil {
			ctx.Warnf("Error dialing to %s: %s", host, err.Error())
			return
		}
		for {
			client := bufio.NewReader(proxyClient)
			remote := bufio.NewReader(targetSiteCon)
			req, err := http.ReadRequest(client)
			if err != nil && err != io.EOF {
				ctx.Warnf("cannot read request of MITM HTTP client: %+#v", err)
			}
			if err != nil {
				return
			}
			req, resp := proxy.filterRequest(req, ctx)
			if resp == nil {
				if err := req.Write(targetSiteCon); err != nil {
					httpError(proxyClient, ctx, err)
					return
				}
				resp, err = http.ReadResponse(remote, req)
				if err != nil {
					httpError(proxyClient, ctx, err)
					return
				}
				defer resp.Body.Close()
			}
			resp = proxy.filterResponse(resp, ctx)
			if err := resp.Write(proxyClient); err != nil {
				httpError(proxyClient, ctx, err)
				return
			}
		}
	case ConnectMitm:
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		ctx.Logf("Assuming CONNECT is TLS, mitm proxying it")
		// this goes in a separate goroutine, so that the net/http server won't think we're
		// still handling the request even after hijacking the connection. Those HTTP CONNECT
		// request can take forever, and the server will be stuck when "closed".
		// TODO: Allow Server.Close() mechanism to shut down this connection as nicely as possible
		tlsConfig := defaultTLSConfig
		if todo.TLSConfig != nil {
			var err error
			tlsConfig, err = todo.TLSConfig(host, ctx)
			if err != nil {
				httpError(proxyClient, ctx, err)
				return
			}
		}
		go func() {
			//TODO: cache connections to the remote website
			rawClientTLS := tls.Server(proxyClient, tlsConfig)
			if err := rawClientTLS.Handshake(); err != nil {
				ctx.Warnf("Cannot handshake client %v %v", r.Host, err)
				return
			}
			defer rawClientTLS.Close()
			clientTLSReader := bufio.NewReader(rawClientTLS)
			for !isEOF(clientTLSReader) {
				req, err := http.ReadRequest(clientTLSReader)
				var ctx = &ProxyCtx{Req: req, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy, UserData: ctx.UserData}
				if err != nil && err != io.EOF {
					return
				}
				if err != nil {
					ctx.Warnf("Cannot read TLS request from mitm'd client %v %v", r.Host, err)
					return
				}
				req.RemoteAddr = r.RemoteAddr // since we're converting the request, need to carry over the original connecting IP as well
				ctx.Logf("req %v", r.Host)
				ctx.Logf("handleTunnel %v request received %v", proxy.BridgeLinkNames, r.Host)
				if !httpsRegexp.MatchString(req.URL.String()) {
					req.URL, err = url.Parse("https://" + r.Host + req.URL.String())
				}
				byteURL := []byte(req.URL.String())
				var byteBody []byte
				r.Body.Read(byteBody)
				wireHTTPRequest := WireMessage{
					SessionId:  util.UniqueId(byteURL),
					Type:       MessageType_request,
					Method:     req.Method,
					URL:        byteURL,
					Header:     HeaderToWireArray(r.Header),
					Body:       byteBody,
					Size:       uint32(len(byteURL)),
					Oridinal:   0,
					Compressed: false,
				}
				data, err := proto.Marshal(wireHTTPRequest.ProtoReflect().Interface())
				if err != nil {
					ctx.Logf("handleTunnel %v marshaling error: %v", proxy.BridgeLinkNames, err)
					continue
				}
				// Send a serialized  protocol buffer request message using a WebRTC channel
				_, err = proxy.DataChannelWriter.Write(data)
				if err != nil {
					ctx.Logf("handleTunnel %v WebRTC DataChannel writer error: %v", proxy.BridgeLinkNames, err)
					continue
				}
				ctx.Logf("handleTunnel sent protocol buffer request message via WebRTC to %v: %v", r.Host, wireHTTPRequest.GetSessionId()[:9])
				counter++
				timeout := time.Second * 10
				response, headers, err := proxy.waitForWebRTCMessage(wireHTTPRequest.GetSessionId(), timeout)
				if err != nil {
					response = []byte(err.Error())
					ctx.Logf("handleTunnel %v error while waiting for WebRTC response for %v: %v", r.Host, wireHTTPRequest.GetSessionId()[:9], err)
				}
				ctx.Logf("handleTunnel response size %d", len(response))
				// Bug fix which goproxy fails to provide request
				// information URL in the context when does HTTPS MITM
				ctx.Req = req
				resp := http.Response{
					Header: w.Header(),
					Body:   ioutil.NopCloser(bytes.NewBuffer(response)),
				}
				text := resp.Status

				statusCode := strconv.Itoa(200) + " "
				if strings.HasPrefix(text, statusCode) {
					text = text[len(statusCode):]
				}

				// always use 1.1 to support chunked encoding
				if _, err := io.WriteString(rawClientTLS, "HTTP/1.1"+" "+statusCode+text+"\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response HTTP status from mitm'd client: %v", err)
					return
				}
				// Since we don't know the length of resp, return chunked encoded response
				// TODO: use a more reasonable scheme
				resp.Header.Del("Content-Length")
				for _, header := range headers {
					if header.Key != "Content-Length" {
						resp.Header.Add(header.Key, header.Value)
					}
				}
				resp.Header.Set("Transfer-Encoding", "chunked")
				// Force connection close otherwise chrome will keep CONNECT tunnel open forever
				resp.Header.Set("Connection", "close")
				if err := resp.Header.Write(rawClientTLS); err != nil {
					ctx.Warnf("Cannot write TLS response header from mitm'd client: %v", err)
					return
				}
				if _, err = io.WriteString(rawClientTLS, "\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response header end from mitm'd client: %v", err)
					return
				}
				chunked := newChunkedWriter(rawClientTLS)
				if _, err := io.Copy(chunked, resp.Body); err != nil {
					ctx.Warnf("Cannot write TLS response body from mitm'd client: %v", err)
					return
				}
				if err := chunked.Close(); err != nil {
					ctx.Warnf("Cannot write TLS chunked EOF from mitm'd client: %v", err)
					return
				}
				if _, err = io.WriteString(rawClientTLS, "\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response chunked trailer from mitm'd client: %v", err)
					return
				}
				ctx.Logf("Reached end \n")
			}
			ctx.Logf("Exiting on EOF")
		}()
	case ConnectProxyAuthHijack:
		proxyClient.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\n"))
		todo.Hijack(r, proxyClient, ctx)
	case ConnectReject:
		if ctx.Resp != nil {
			if err := ctx.Resp.Write(proxyClient); err != nil {
				ctx.Warnf("Cannot write response that reject http CONNECT: %v", err)
			}
		}
		proxyClient.Close()
	}
}

func (proxy *ProxyHTTPServer) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	ctx := &ProxyCtx{Req: r, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy, certStore: proxy.CertStore}

	hij, ok := w.(http.Hijacker)
	if !ok {
		panic("httpserver does not support hijacking")
	}

	proxyClient, _, e := hij.Hijack()
	if e != nil {
		panic("Cannot hijack connection " + e.Error())
	}

	ctx.Logf("Running %d CONNECT handlers", len(proxy.httpsHandlers))
	todo, host := OkConnect, r.URL.Host
	for i, h := range proxy.httpsHandlers {
		newtodo, newhost := h.HandleConnect(host, ctx)

		// If found a result, break the loop immediately
		if newtodo != nil {
			todo, host = newtodo, newhost
			ctx.Logf("on %dth handler: %v %s", i, todo, host)
			break
		}
	}
	switch todo.Action {
	case ConnectAccept:
		if !hasPort.MatchString(host) {
			host += ":80"
		}
		targetSiteCon, err := proxy.connectDial("tcp", host)
		if err != nil {
			httpError(proxyClient, ctx, err)
			return
		}
		ctx.Logf("Accepting CONNECT to %s", host)
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

		targetTCP, targetOK := targetSiteCon.(halfClosable)
		proxyClientTCP, clientOK := proxyClient.(halfClosable)
		if targetOK && clientOK {
			go copyAndClose(ctx, targetTCP, proxyClientTCP)
			go copyAndClose(ctx, proxyClientTCP, targetTCP)
		} else {
			go func() {
				var wg sync.WaitGroup
				wg.Add(2)
				go copyOrWarn(ctx, targetSiteCon, proxyClient, &wg)
				go copyOrWarn(ctx, proxyClient, targetSiteCon, &wg)
				wg.Wait()
				proxyClient.Close()
				targetSiteCon.Close()

			}()
		}

	case ConnectHijack:
		todo.Hijack(r, proxyClient, ctx)
	case ConnectHTTPMitm:
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		ctx.Logf("Assuming CONNECT is plain HTTP tunneling, mitm proxying it")
		targetSiteCon, err := proxy.connectDial("tcp", host)
		if err != nil {
			ctx.Warnf("Error dialing to %s: %s", host, err.Error())
			return
		}
		for {
			client := bufio.NewReader(proxyClient)
			remote := bufio.NewReader(targetSiteCon)
			req, err := http.ReadRequest(client)
			if err != nil && err != io.EOF {
				ctx.Warnf("cannot read request of MITM HTTP client: %+#v", err)
			}
			if err != nil {
				return
			}
			req, resp := proxy.filterRequest(req, ctx)
			if resp == nil {
				if err := req.Write(targetSiteCon); err != nil {
					httpError(proxyClient, ctx, err)
					return
				}
				resp, err = http.ReadResponse(remote, req)
				if err != nil {
					httpError(proxyClient, ctx, err)
					return
				}
				defer resp.Body.Close()
			}
			resp = proxy.filterResponse(resp, ctx)
			if err := resp.Write(proxyClient); err != nil {
				httpError(proxyClient, ctx, err)
				return
			}
		}
	case ConnectMitm:
		proxyClient.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
		ctx.Logf("Assuming CONNECT is TLS, mitm proxying it")
		// this goes in a separate goroutine, so that the net/http server won't think we're
		// still handling the request even after hijacking the connection. Those HTTP CONNECT
		// request can take forever, and the server will be stuck when "closed".
		// TODO: Allow Server.Close() mechanism to shut down this connection as nicely as possible
		tlsConfig := defaultTLSConfig
		if todo.TLSConfig != nil {
			var err error
			tlsConfig, err = todo.TLSConfig(host, ctx)
			if err != nil {
				httpError(proxyClient, ctx, err)
				return
			}
		}
		go func() {
			//TODO: cache connections to the remote website
			rawClientTLS := tls.Server(proxyClient, tlsConfig)
			if err := rawClientTLS.Handshake(); err != nil {
				ctx.Warnf("Cannot handshake client %v %v", r.Host, err)
				return
			}
			defer rawClientTLS.Close()
			clientTLSReader := bufio.NewReader(rawClientTLS)
			for !isEOF(clientTLSReader) {
				req, err := http.ReadRequest(clientTLSReader)
				var ctx = &ProxyCtx{Req: req, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy, UserData: ctx.UserData}
				if err != nil && err != io.EOF {
					return
				}
				if err != nil {
					ctx.Warnf("Cannot read TLS request from mitm'd client %v %v", r.Host, err)
					return
				}
				req.RemoteAddr = r.RemoteAddr // since we're converting the request, need to carry over the original connecting IP as well
				ctx.Logf("req %v", r.Host)

				if !httpsRegexp.MatchString(req.URL.String()) {
					req.URL, err = url.Parse("https://" + r.Host + req.URL.String())
				}

				// Bug fix which goproxy fails to provide request
				// information URL in the context when does HTTPS MITM
				ctx.Req = req

				req, resp := proxy.filterRequest(req, ctx)
				if resp == nil {
					if isWebSocketRequest(req) {
						ctx.Logf("Request looks like websocket upgrade.")
						proxy.serveWebsocketTLS(ctx, w, req, tlsConfig, rawClientTLS)
						return
					}
					if err != nil {
						ctx.Warnf("Illegal URL %s", "https://"+r.Host+req.URL.Path)
						return
					}
					removeProxyHeaders(ctx, req)
					resp, err = ctx.RoundTrip(req)
					if err != nil {
						ctx.Warnf("Cannot read TLS response from mitm'd server %v", err)
						return
					}
					ctx.Logf("resp %v", resp.Status)
				}
				resp = proxy.filterResponse(resp, ctx)
				defer resp.Body.Close()

				text := resp.Status
				statusCode := strconv.Itoa(resp.StatusCode) + " "
				if strings.HasPrefix(text, statusCode) {
					text = text[len(statusCode):]
				}
				// always use 1.1 to support chunked encoding
				if _, err := io.WriteString(rawClientTLS, "HTTP/1.1"+" "+statusCode+text+"\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response HTTP status from mitm'd client: %v", err)
					return
				}
				// Since we don't know the length of resp, return chunked encoded response
				// TODO: use a more reasonable scheme
				resp.Header.Del("Content-Length")
				resp.Header.Set("Transfer-Encoding", "chunked")
				// Force connection close otherwise chrome will keep CONNECT tunnel open forever
				resp.Header.Set("Connection", "close")
				if err := resp.Header.Write(rawClientTLS); err != nil {
					ctx.Warnf("Cannot write TLS response header from mitm'd client: %v", err)
					return
				}
				if _, err = io.WriteString(rawClientTLS, "\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response header end from mitm'd client: %v", err)
					return
				}
				chunked := newChunkedWriter(rawClientTLS)
				if _, err := io.Copy(chunked, resp.Body); err != nil {
					ctx.Warnf("Cannot write TLS response body from mitm'd client: %v", err)
					return
				}
				if err := chunked.Close(); err != nil {
					ctx.Warnf("Cannot write TLS chunked EOF from mitm'd client: %v", err)
					return
				}
				if _, err = io.WriteString(rawClientTLS, "\r\n"); err != nil {
					ctx.Warnf("Cannot write TLS response chunked trailer from mitm'd client: %v", err)
					return
				}
			}
			ctx.Logf("Exiting on EOF")
		}()
	case ConnectProxyAuthHijack:
		proxyClient.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\n"))
		todo.Hijack(r, proxyClient, ctx)
	case ConnectReject:
		if ctx.Resp != nil {
			if err := ctx.Resp.Write(proxyClient); err != nil {
				ctx.Warnf("Cannot write response that reject http CONNECT: %v", err)
			}
		}
		proxyClient.Close()
	}
}

func httpError(w io.WriteCloser, ctx *ProxyCtx, err error) {
	if _, err := io.WriteString(w, "HTTP/1.1 502 Bad Gateway\r\n\r\n"); err != nil {
		ctx.Warnf("Error responding to client: %s", err)
	}
	if err := w.Close(); err != nil {
		ctx.Warnf("Error closing client connection: %s", err)
	}
}

func copyOrWarn(ctx *ProxyCtx, dst io.Writer, src io.Reader, wg *sync.WaitGroup) {
	if _, err := io.Copy(dst, src); err != nil {
		ctx.Warnf("Error copying to client: %s", err)
	}
	wg.Done()
}

func copyAndClose(ctx *ProxyCtx, dst, src halfClosable) {
	if _, err := io.Copy(dst, src); err != nil {
		ctx.Warnf("Error copying to client: %s", err)
	}

	dst.CloseWrite()
	src.CloseRead()
}

func dialerFromEnv(proxy *ProxyHTTPServer) func(network, addr string) (net.Conn, error) {
	https_proxy := os.Getenv("HTTPS_PROXY")
	if https_proxy == "" {
		https_proxy = os.Getenv("https_proxy")
	}
	if https_proxy == "" {
		return nil
	}
	return proxy.NewConnectDialToProxy(https_proxy)
}

func (proxy *ProxyHTTPServer) NewConnectDialToProxy(https_proxy string) func(network, addr string) (net.Conn, error) {
	return proxy.NewConnectDialToProxyWithHandler(https_proxy, nil)
}

func (proxy *ProxyHTTPServer) NewConnectDialToProxyWithHandler(https_proxy string, connectReqHandler func(req *http.Request)) func(network, addr string) (net.Conn, error) {
	u, err := url.Parse(https_proxy)
	if err != nil {
		return nil
	}
	if u.Scheme == "" || u.Scheme == "http" {
		if strings.IndexRune(u.Host, ':') == -1 {
			u.Host += ":80"
		}
		return func(network, addr string) (net.Conn, error) {
			connectReq := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: addr},
				Host:   addr,
				Header: make(http.Header),
			}
			if connectReqHandler != nil {
				connectReqHandler(connectReq)
			}
			c, err := proxy.dial(network, u.Host)
			if err != nil {
				return nil, err
			}
			connectReq.Write(c)
			// Read response.
			// Okay to use and discard buffered reader here, because
			// TLS server will not speak until spoken to.
			br := bufio.NewReader(c)
			resp, err := http.ReadResponse(br, connectReq)
			if err != nil {
				c.Close()
				return nil, err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				resp, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				c.Close()
				return nil, errors.New("proxy refused connection" + string(resp))
			}
			return c, nil
		}
	}
	if u.Scheme == "https" || u.Scheme == "wss" {
		if strings.IndexRune(u.Host, ':') == -1 {
			u.Host += ":443"
		}
		return func(network, addr string) (net.Conn, error) {
			c, err := proxy.dial(network, u.Host)
			if err != nil {
				return nil, err
			}
			c = tls.Client(c, proxy.Tr.TLSClientConfig)
			connectReq := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: addr},
				Host:   addr,
				Header: make(http.Header),
			}
			if connectReqHandler != nil {
				connectReqHandler(connectReq)
			}
			connectReq.Write(c)
			// Read response.
			// Okay to use and discard buffered reader here, because
			// TLS server will not speak until spoken to.
			br := bufio.NewReader(c)
			resp, err := http.ReadResponse(br, connectReq)
			if err != nil {
				c.Close()
				return nil, err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 500))
				if err != nil {
					return nil, err
				}
				c.Close()
				return nil, errors.New("proxy refused connection" + string(body))
			}
			return c, nil
		}
	}
	return nil
}

func TLSConfigFromCA(ca *tls.Certificate) func(host string, ctx *ProxyCtx) (*tls.Config, error) {
	return func(host string, ctx *ProxyCtx) (*tls.Config, error) {
		var err error
		var cert *tls.Certificate

		hostname := stripPort(host)
		config := defaultTLSConfig.Clone()
		ctx.Logf("signing for %s", stripPort(host))

		genCert := func() (*tls.Certificate, error) {
			return signHost(*ca, []string{hostname})
		}
		if ctx.certStore != nil {
			cert, err = ctx.certStore.Fetch(hostname, genCert)
		} else {
			cert, err = genCert()
		}

		if err != nil {
			ctx.Warnf("Cannot sign host certificate with provided CA: %s", err)
			return nil, err
		}

		config.Certificates = append(config.Certificates, *cert)
		return config, nil
	}
}
