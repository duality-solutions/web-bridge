package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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
		buffer := make([]byte, 128000)
		_, err := d.Read(buffer)
		if err != nil {
			fmt.Println("Datachannel closed; Exit the readloop:", err)
			return
		}
		buffer = bytes.Trim(buffer, "\x00")
		counter++
		if len(buffer) > 100 {
			fmt.Println("ReadLoop Message from DataChannel:", counter, string(buffer[:100]))
			fmt.Println("ReadLoop Message from DataChannel Len:", counter, len(buffer))
		} else {
			fmt.Println("ReadLoop Message from DataChannel:", counter, string(buffer))
		}

		go sendResponse(buffer)
	}
}

func sendResponse(data []byte) {
	targetURL := string(data)
	fmt.Println("sendResponse before http.Client", targetURL)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Get(targetURL)
	if err != nil {
		fmt.Println("sendResponse client.Get error: ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyLen := len(body)
	if bodyLen > 50000 {
		datawriter.Write(body[:50000])
	}

	fmt.Println("sendResponse client.Get successful", len(string(body)))
}
