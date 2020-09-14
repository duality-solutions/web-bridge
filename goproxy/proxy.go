package goproxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
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

	"github.com/duality-solutions/web-bridge/internal/util"
	"google.golang.org/protobuf/proto"
)

// MaxTransmissionBytes is the maxium bytes allow per WebRTC send
var MaxTransmissionBytes = 65535

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
	BridgeID          string
	BridgeLinkNames   string
	DataChannelWriter io.Writer
	DataChannelReader io.Reader
	mapWebRTCMessages map[string]chan *WireMessage
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
func (proxy *ProxyHTTPServer) readWebRTCMessageLoop(ctx *ProxyCtx) {
	proxy.mapWebRTCMessages = make(map[string]chan *WireMessage)
	// TODO: add a channel to stop this loop
	for {
		buffer := make([]byte, MaxTransmissionBytes)
		_, err := proxy.DataChannelReader.Read(buffer)
		if err != nil {
			ctx.Logf("readWebRTCMessageLoop Read error: %v", err)
			return
		}
		buffer = bytes.Trim(buffer, "\x00")
		if len(buffer) > 32 {
			wr := WireMessage{}
			err = proto.Unmarshal(buffer, &wr)
			if err != nil {
				ctx.Logf("readWebRTCMessageLoop (%v) data length (%v) unmarshaling data size %d error: (%v)", proxy.BridgeLinkNames, proxy.BridgeID, len(buffer), err)
				continue
			}
			if wr.GetType() == MessageType_response {
				sessionID := wr.GetSessionId()
				proxy.mapWebRTCMessages[sessionID] <- &wr
				defer close(proxy.mapWebRTCMessages[sessionID])
			} else if wr.GetType() == MessageType_request {
				go proxy.sendResponse(&wr, ctx)
			} else {
				ctx.Logf("readWebRTCMessageLoop unknown message type received %v %v", proxy.BridgeLinkNames, proxy.BridgeID)
			}
		} else {
			ctx.Logf("readWebRTCMessageLoop short message from %v %v :%v", proxy.BridgeLinkNames, proxy.BridgeID, string(buffer))
		}
	}
}

// waitForWebRTCMessage tries to get a response for the given sessionID before the timeout duration
func (proxy *ProxyHTTPServer) waitForWebRTCMessage(sessionID string, timeout time.Duration) ([]byte, []*HttpHeader, error) {
	messages := make(map[uint32]*WireMessage)
	var response []byte
	var headers []*HttpHeader
	var extraSize int = 0
	var max, chunks uint32 = 0, 0
	// Initialize map session id channel
	proxy.mapWebRTCMessages[sessionID] = make(chan *WireMessage, 1)
	for {
		select {
		case wireResponse := <-proxy.mapWebRTCMessages[sessionID]:
			if headers == nil {
				headers = wireResponse.Header
				extraSize = (100 * len(headers)) + 200
				max = uint32(MaxTransmissionBytes - extraSize)
				chunks = (wireResponse.GetSize() / max) + 1
			}
			messages[wireResponse.GetOridinal()] = wireResponse
			if uint32(len(messages)) >= chunks {
				// make sure response is created with the correct ordinal order
				for i := uint32(0); i < chunks; i++ {
					response = append(response, messages[i].GetBody()...)
				}
				// check if last char is null before calling trim
				if len(response) > 0 {
					lastByByte := response[len(response)-1]
					if lastByByte == 0 {
						response = bytes.Trim(response, "\x00") // remove trailing null characters from response
					}
				}
				return response, headers, nil
			}
			proxy.mapWebRTCMessages[sessionID] = make(chan *WireMessage, 1)
		case <-time.After(timeout):
			return response, nil, fmt.Errorf("waitForWebRTCMessage response for %v %v, timeout for session %v", proxy.BridgeLinkNames, proxy.BridgeID, sessionID[:9])
		}
	}
}

func (proxy *ProxyHTTPServer) sendResponse(wrReq *WireMessage, ctx *ProxyCtx) {
	targetURL := string(wrReq.URL)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	reqBodyCloser := ioutil.NopCloser(bytes.NewBuffer(wrReq.GetBody()))
	req, err := http.NewRequest(wrReq.Method, targetURL, reqBodyCloser)
	req.Proto = "HTTP/1.1"
	req.Header.Add("Cache-Control", "no-cache")
	for _, head := range wrReq.GetHeader() {
		req.Header.Add(head.Key, head.Value)
	}
	resp, err := client.Do(req)
	if err != nil {
		respError := http.Response{
			Body: ioutil.NopCloser(bytes.NewBuffer([]byte(err.Error()))),
		}
		resp = &respError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyLen := uint32(len(body))
	headers := HeaderToWireArray(resp.Header)
	extraSize := (100 * len(headers)) + 200
	max := uint32(MaxTransmissionBytes - extraSize)
	if bodyLen > max {
		chunks := bodyLen/max + 1
		pos := uint32(0)
		for i := uint32(0); i < chunks; i++ {
			if i != chunks {
				wrResp := WireMessage{
					SessionId:  wrReq.SessionId,
					Type:       MessageType_response,
					Method:     wrReq.Method,
					URL:        wrReq.URL,
					Header:     HeaderToWireArray(resp.Header),
					Body:       body[pos : pos+max],
					Size:       bodyLen,
					Oridinal:   i,
					Compressed: false,
				}
				protoData, err := proto.Marshal(wrResp.ProtoReflect().Interface())
				if err != nil {
					ctx.Logf("sendResponse %v marshaling error: %v", proxy.BridgeLinkNames, err)
				}
				_, err = proxy.DataChannelWriter.Write(protoData)
				if err != nil {
					ctx.Logf("sendResponse %v DataChannelWriter.Write error: %v", proxy.BridgeLinkNames, err)
				} else {
					ctx.Logf("sendResponse DataChannelWriter.Write protoData len %d", len(protoData))
				}
			} else {
				wrResp := WireMessage{
					SessionId:  wrReq.SessionId,
					Type:       MessageType_response,
					Method:     wrReq.Method,
					URL:        wrReq.URL,
					Header:     HeaderToWireArray(resp.Header),
					Body:       body[pos : bodyLen-pos],
					Size:       bodyLen,
					Oridinal:   0,
					Compressed: false,
				}
				protoData, err := proto.Marshal(wrResp.ProtoReflect().Interface())
				if err != nil {
					ctx.Logf("sendResponse marshaling error: %v", err)
				}
				proxy.DataChannelWriter.Write(protoData)
				_, err = proxy.DataChannelWriter.Write(protoData)
				if err != nil {
					ctx.Logf("sendResponse %v DataChannelWriter.Write error: %v", proxy.BridgeLinkNames, err)
				} else {
					ctx.Logf("sendResponse %v DataChannelWriter.Write protoData len %d", proxy.BridgeLinkNames, len(protoData))
				}
			}
			pos = pos + max
		}
	} else {
		wrResp := WireMessage{
			SessionId:  wrReq.SessionId,
			Type:       MessageType_response,
			Method:     wrReq.Method,
			URL:        wrReq.URL,
			Header:     HeaderToWireArray(resp.Header),
			Body:       body,
			Size:       bodyLen,
			Oridinal:   0,
			Compressed: false,
		}
		//fmt.Println("sendResponse body ", string(wrResp.GetBody()))
		protoData, err := proto.Marshal(wrResp.ProtoReflect().Interface())

		if err != nil {
			ctx.Logf("sendResponse %v marshaling error: %v", proxy.BridgeLinkNames, err)
		}
		_, err = proxy.DataChannelWriter.Write(protoData)
		if err != nil {
			ctx.Logf("sendResponse %v DataChannelWriter.Write error: ", proxy.BridgeLinkNames, err)
		} else {
			ctx.Logf("sendResponse %v DataChannelWriter.Write protoData len %d", proxy.BridgeLinkNames, len(protoData))
		}
	}
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
func HeaderToWireArray(header http.Header) (res []*HttpHeader) {
	for name, values := range header {
		for _, value := range values {
			item := HttpHeader{
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
	ctx := &ProxyCtx{Req: r, Session: atomic.AddInt64(&proxy.sess, 1), Proxy: proxy}
	go proxy.readWebRTCMessageLoop(ctx)
	if r.Method == "CONNECT" {
		proxy.handleTunnel(w, r)
	} else {
		reqBody, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		var err error
		ctx.Logf("Got request %v %v %v %v %v", r.URL.Path, r.Host, r.Method, r.URL.String(), string(reqBody))
		byteURL := []byte(r.URL.String())
		wireRequest := WireMessage{
			SessionId:  util.UniqueId(byteURL),
			Type:       MessageType_request,
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
			ctx.Logf("ServeHTTP %v marshaling error: %v", proxy.BridgeLinkNames, err)
			return
		}
		_, err = proxy.DataChannelWriter.Write(data)
		if err != nil {
			ctx.Logf("WebRTC DataChannel %v writer error: %v", proxy.BridgeLinkNames, err)
			return
		}
		ctx.Logf("ServeHTTP sent protocol buffer request message via WebRTC to %v: %v", r.Host, wireRequest.GetSessionId()[:9])
		counter++
		timeout := time.Second * 10
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
