package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/pion/webrtc/v2"
)

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

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

			// Detach the data channel
			raw, dErr := d.Detach()
			if dErr != nil {
				panic(dErr)
			}
			datawriter = raw
			// Handle reading from the data channel
			go ReadLoop(raw)
		})
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	util.Decode(util.MustReadStdin(), &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(util.Encode(answer))

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
		counter++
		fmt.Println("ReadLoop Message from DataChannel:", counter, string(buffer))
		go sendResponse(buffer)
	}
}

func sendResponse(data []byte) {
	var w http.ResponseWriter
	fmt.Println("sendResponse", string(data), w)
	bufReader := bytes.NewReader(data)
	bufIO := bufio.NewReader(bufReader)
	req, err := http.ReadRequest(bufIO) // deserialize request
	if err != nil {                     // this is a response
		fmt.Println("Datachannel ReadRequest error", err)
	} else {
		fmt.Println("sendResponse before DialTimeout", req.Host)
		destConn, err2 := net.DialTimeout("tcp", req.Host, 10*time.Second)
		if err2 != nil {
			fmt.Println("Datachannel DialTimeout error", err)
			return
		}
		buffer := make([]byte, 64000)
		defer destConn.Close()
		destConn.Read(buffer)
		buffer = bytes.Trim(buffer, "\x00")
		fmt.Println("sendResponse destConn", string(buffer))
		_, err = datawriter.Write(buffer)
		if err != nil {
			fmt.Println("sendResponse datawriter Write error", err)
		}
	}
}
