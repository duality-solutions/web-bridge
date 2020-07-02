package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/pion/webrtc/v2"
)

const messageSize = 15

var dataChannel *webrtc.DataChannel
var datawriter io.Writer
var counter = 0

func main() {
	// Since this behavior diverges from the WebRTC API it has to be
	// enabled using a settings engine. Mixing both detached and the
	// OnMessage DataChannel API is not supported.

	// Create a SettingEngine and enable Detach
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()

	// Create an API object with the engine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection using the API object
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Create a datachannel with label 'data'
	dataChannel, err = peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", dataChannel.Label(), dataChannel.ID())

		// Detach the data channel
		raw, err := dataChannel.Detach()
		if err != nil {
			panic(err)
		}

		go ReadLoop(raw)
		datawriter = raw
		// Handle reading from the data channel
		StartBridgeNetwork()
	})

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Output the offer in base64 so we can paste it in browser
	fmt.Println(util.Encode(offer))

	// Wait for the answer to be pasted
	answer := webrtc.SessionDescription{}
	util.Decode(util.MustReadStdin(), &answer)

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block forever
	select {}
}

// ReadLoop shows how to read from the datachannel directly
func ReadLoop(d io.Reader) {
	for {
		buffer := make([]byte, 64000)
		_, err := d.Read(buffer)
		if err != nil {
			fmt.Println("Datachannel closed; Exit the readloop:", err)
			return
		}
		buffer = bytes.Trim(buffer, "\x00")
		fmt.Println("ReadLoop Message from DataChannel:", string(buffer))
	}
}

// StartBridgeNetwork listens to a port for http traffic and routes it through a link's WebRTC channel
func StartBridgeNetwork() {
	server := &http.Server{
		Addr: ":7777",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunnel(w, r)
			} else {
				http.Error(w, "HTTP not supported", http.StatusNotImplemented)
				return
			}
		}),
		ConnState: onConnStateEvent,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	go server.ListenAndServe()
}

func transferCloser(dest io.WriteCloser, src io.ReadCloser) {
	defer dest.Close()
	defer src.Close()
	io.Copy(dest, src)
}

// handleTunnel handles link bridge tunnel connection
func handleTunnel(w http.ResponseWriter, r *http.Request) {
	counter++
	fmt.Println("handleTunnel", r.Host, counter)
	byteRequest, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("handleTunnel DumpRequest error", r.Host, err.Error())
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}
	fmt.Println("handleTunnel", "byteRequest len", len(byteRequest))
	fmt.Println("handleTunnel Request", string(byteRequest))
	err = dataChannel.Send(byteRequest)
	if err != nil {
		fmt.Println("handleTunnel Send error", err)
	}
	reqReader := bytes.NewReader(byteRequest)
	reqCloser := ioutil.NopCloser(reqReader)
	//io.Copy(datawriter, reqCloser)
	// todo: wrap send with standard envelop so receive knows if it is a request or response
	// todo: get response from WebRTC messages
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		fmt.Println("Hijacking not supported", r.Host)
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		fmt.Println("handleTunnel Hijack error", r.Host)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transferCloser(clientConn, reqCloser)
}

func onConnStateEvent(conn net.Conn, state http.ConnState) {
	fmt.Println("onConnStateEvent", state.String())
}
