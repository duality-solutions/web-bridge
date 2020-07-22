package goproxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/duality-solutions/web-bridge/bridge"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"google.golang.org/protobuf/proto"
)

//var channelsMap map[string](chan *TestStruct)

// ProxyHTTPServer is the basic proxy type. Implements http.Handler.
type ProxyHTTPServer struct {
	// session variable must be aligned in i386
	// see http://golang.org/src/pkg/sync/atomic/doc.go#L41
	sess int64
	// KeepDestinationHeaders indicates the proxy should retain any headers present in the http.Response before proxying
	KeepDestinationHeaders bool
	// setting Verbose to true will log information on each request sent to the proxy
	Verbose         bool
	Logger          Logger
	NonProxyHandler http.Handler
	reqHandlers     []ReqHandler
	respHandlers    []RespHandler
	httpsHandlers   []HttpsHandler
	Tr              *http.Transport
	// ConnectDial will be used to create TCP connections for CONNECT requests
	// if nil Tr.Dial will be used
	ConnectDial       func(network string, addr string) (net.Conn, error)
	CertStore         CertStorage
	DataChannelWriter io.Writer
	DataChannelReader io.Reader
	mapWebRTCMessages map[string]*bridge.WireMessage
}

var hasPort = regexp.MustCompile(`:\d+$`)

func copyHeaders(dst, src http.Header, keepDestHeaders bool) {
	if !keepDestHeaders {
		for k := range dst {
			dst.Del(k)
		}
	}
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}

func isEOF(r *bufio.Reader) bool {
	_, err := r.Peek(1)
	if err == io.EOF {
		return true
	}
	return false
}

// readWebRTCMessageLoop creates a process that continues to read data from the WebRTC channel
func (proxy *ProxyHTTPServer) readWebRTCMessageLoop() {
	proxy.mapWebRTCMessages = make(map[string]*bridge.WireMessage, 0)
	for {
		buffer := make([]byte, bridge.MaxTransmissionBytes)
		_, err := proxy.DataChannelReader.Read(buffer)
		if err != nil {
			fmt.Println("readWebRTCMessageLoop Read error:", err)
			return
		}
		buffer = bytes.Trim(buffer, "\x00")
		wr := bridge.WireMessage{}
		err = proto.Unmarshal(buffer, &wr)
		if err != nil {
			log.Fatal("readWebRTCMessageLoop unmarshaling error:", err)
			continue
		} else {
			fmt.Println("readWebRTCMessageLoop data received:", wr.SessionId, len(buffer), wr.Oridinal)
		}
		proxy.mapWebRTCMessages[wr.GetSessionId()] = &wr
		fmt.Println("readWebRTCMessageLoop channel triggered:", wr.SessionId, len(buffer))
	}
}

// waitForWebRTCMessage tries to get a response for the given sessionID before the timeout duration
func (proxy *ProxyHTTPServer) waitForWebRTCMessage(sessionID string, timeout time.Duration) ([]byte, []*bridge.HttpHeader, error) {
	var response []byte
	max := uint32(bridge.MaxTransmissionBytes - 300)
	// TODO: Use a more effient method to wait for a reponse and add a timeout.
	for uint64(len(proxy.mapWebRTCMessages)) < 1 {
		time.Sleep(11 * time.Millisecond)
	}
	wm := proxy.mapWebRTCMessages[sessionID]
	headers := wm.Header
	chunks := uint32((wm.GetSize() / max) + 1)
	if chunks > 1 {
		for i := uint32(0); i < chunks; i++ {
			wm := proxy.mapWebRTCMessages[sessionID]
			response = append(response, wm.GetBody()...)
		}
	} else {
		response = wm.GetBody()
	}
	return response, headers, nil
}

func (proxy *ProxyHTTPServer) filterRequest(r *http.Request, ctx *ProxyCtx) (req *http.Request, resp *http.Response) {
	req = r
	for _, h := range proxy.reqHandlers {
		req, resp = h.Handle(r, ctx)
		// non-nil resp means the handler decided to skip sending the request
		// and return canned response instead.
		if resp != nil {
			break
		}
	}
	return
}

func (proxy *ProxyHTTPServer) filterResponse(respOrig *http.Response, ctx *ProxyCtx) (resp *http.Response) {
	resp = respOrig
	for _, h := range proxy.respHandlers {
		ctx.Resp = resp
		resp = h.Handle(resp, ctx)
	}
	return
}

func removeProxyHeaders(ctx *ProxyCtx, r *http.Request) {
	r.RequestURI = "" // this must be reset when serving a request with the client
	ctx.Logf("Sending request %v %v", r.Method, r.URL.String())
	// If no Accept-Encoding header exists, Transport will add the headers it can accept
	// and would wrap the response body with the relevant reader.
	r.Header.Del("Accept-Encoding")
	// curl can add that, see
	// https://jdebp.eu./FGA/web-proxy-connection-header.html
	r.Header.Del("Proxy-Connection")
	r.Header.Del("Proxy-Authenticate")
	r.Header.Del("Proxy-Authorization")
	// Connection, Authenticate and Authorization are single hop Header:
	// http://www.w3.org/Protocols/rfc2616/rfc2616.txt
	// 14.10 Connection
	//   The Connection general-header field allows the sender to specify
	//   options that are desired for that particular connection and MUST NOT
	//   be communicated by proxies over further connections.
	r.Header.Del("Connection")
}

// HeaderToWireArray converts a http header to struct slice
func HeaderToWireArray(header http.Header) (res []*bridge.HttpHeader) {
	for name, values := range header {
		for _, value := range values {
			item := bridge.HttpHeader{
				Key:   name,
				Value: value,
			}
			res = append(res, &item)
		}
	}
	return
}

// Standard net/http function. Shouldn't be used directly, http.Serve will use it.
func (proxy *ProxyHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//r.Header["X-Forwarded-For"] = w.RemoteAddr()
	go proxy.readWebRTCMessageLoop()

	if r.Method == "CONNECT" {
		proxy.handleTunnel(w, r)
	} else {
		ctx := &ProxyCtx{Req: r, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy}
		reqBody, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		var err error
		ctx.Logf("Got request %v %v %v %v %v", r.URL.Path, r.Host, r.Method, r.URL.String(), string(reqBody))
		byteURL := []byte(r.URL.String())
		wireRequest := bridge.WireMessage{
			SessionId:  util.UniqueId(byteURL),
			Type:       bridge.MessageType_request,
			Method:     r.Method,
			URL:        byteURL,
			Header:     HeaderToWireArray(r.Header),
			Body:       reqBody,
			Size:       uint32(len(byteURL)),
			Oridinal:   0,
			Compressed: false,
		}
		data, err := proto.Marshal(wireRequest.ProtoReflect().Interface())
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		_, err = proxy.DataChannelWriter.Write(data)
		if err != nil {
			log.Fatal("WebRTC DataChannel writer error: ", err)
		}
		ctx.Logf("ServeHTTP sent protocol buffer request message via WebRTC to %v: %v", r.Host, wireRequest.GetSessionId())
		counter++
		timeout := time.Second * 30
		response, headers, err := proxy.waitForWebRTCMessage(wireRequest.GetSessionId(), timeout)
		if err != nil {
			response = []byte(err.Error())
			ctx.Logf("ServeHTTP %v error while waiting for WebRTC response for %v: %v", r.Host, wireRequest.GetSessionId(), err)
		}
		ctx.Logf("ServeHTTP response size %d", len(response))

		resp := http.Response{
			Header: w.Header(),
			Body:   ioutil.NopCloser(bytes.NewBuffer(response)),
		}
		text := resp.Status

		statusCode := strconv.Itoa(200) + " "
		if strings.HasPrefix(text, statusCode) {
			text = text[len(statusCode):]
		}
		resp.Header.Del("Content-Length")
		for _, header := range headers {
			//ctx.Logf("ServeHTTP creating header: key %v, value %v", head.Key, head.Value)
			if header.Key != "Content-Length" {
				resp.Header.Add(header.Key, header.Value)
			}
		}
		// http.ResponseWriter will take care of filling the correct response length
		// Setting it now, might impose wrong value, contradicting the actual new
		// body the user returned.
		// We keep the original body to remove the header only if things changed.
		// This will prevent problems with HEAD requests where there's no body, yet,
		// the Content-Length header should be set.
		resp.StatusCode = 200
		//resp = proxy.filterResponse(&resp, ctx)
		ctx.Logf("ServeHTTP Copying response to client %v [%d]", resp.Status, resp.StatusCode)
		// Force connection close otherwise chrome will keep CONNECT tunnel open forever
		resp.Header.Set("Connection", "close")
		w.WriteHeader(resp.StatusCode)
		nr, err := io.Copy(w, resp.Body)
		if err != nil {
			ctx.Warnf("ServeHTTP Can't copy reponse body to writer %v", err)
		} else {
			ctx.Logf("ServeHTTP Copied %v bytes to response writer", nr)
		}
		if err := resp.Body.Close(); err != nil {
			ctx.Warnf("ServeHTTP Can't close response body %v", err)
		}
	}
}

// NewProxyHTTPServer creates and returns a proxy server, logging to stderr by default
func NewProxyHTTPServer() *ProxyHTTPServer {
	proxy := ProxyHTTPServer{
		Logger:        log.New(os.Stderr, "", log.LstdFlags),
		reqHandlers:   []ReqHandler{},
		respHandlers:  []RespHandler{},
		httpsHandlers: []HttpsHandler{},
		NonProxyHandler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, "This is a proxy server. Does not respond to non-proxy requests.", 500)
		}),
		Tr: &http.Transport{TLSClientConfig: tlsClientSkipVerify, Proxy: http.ProxyFromEnvironment},
	}
	proxy.ConnectDial = dialerFromEnv(&proxy)

	return &proxy
}
