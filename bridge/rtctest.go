package bridge

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/pion/webrtc/v2"
)

func TestCreateOffer() {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println("TestCreateOffer NewPeerConnection error ", err)
		return
	}

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		fmt.Println("TestCreateOffer CreateDataChannel error ", err)
		return
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		//report := peerConnection.GetStats()
		//connStats, _ := report.GetConnectionStats(peerConnection)
		//fmt.Println("GetConnectionStats", connStats)
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(5 * time.Second).C {
			message, _ := util.RandomString(16)
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendErr := dataChannel.SendText(message)
			if sendErr != nil {
				panic(sendErr)
			}
		}
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		fmt.Println("TestCreateOffer CreateOffer error ", err)
		return
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		fmt.Println("TestCreateOffer SetLocalDescription error ", err)
	}
	ret, _ := util.EncodeObject(offer)
	fmt.Println(ret)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("paste answer> ")
	answerSDP, _ := reader.ReadString('\n')
	var answer webrtc.SessionDescription
	err = util.DecodeObject(answerSDP, &answer)
	if err != nil {
		fmt.Println("TestCreateOffer DecodeObject error ", err)
		return
	}
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		fmt.Println("TestCreateOffer SetRemoteDescription error ", err)
		return
	}

	// Block forever
	select {}
}

func TestWaitForOffer(offer string) {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		fmt.Println("TestWaitForOffer Error with NewPeerConnection", err)
		return
	}
	fmt.Println("ICE Connection config", config)
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("TestWaitForOffer ICE Connection State has changed: %s\n", connectionState.String())
		//report := peerConnection.GetStats()
		//connStats, _ := report.GetConnectionStats(peerConnection)
		//fmt.Println("TestWaitForOffer GetConnectionStats", connStats)
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("TestWaitForOffer New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("TestWaitForOffer Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(5 * time.Second).C {
				message, _ := util.RandomString(16)
				fmt.Printf("TestWaitForOffer Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					panic(sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("TestWaitForOffer Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})
	var offerDescription webrtc.SessionDescription
	err = util.DecodeObject(offer, &offerDescription)
	if err != nil {
		fmt.Println("TestWaitForOffer Error with DecodeObject", err)
		return
	}
	err = peerConnection.SetRemoteDescription(offerDescription)
	if err != nil {
		fmt.Println("TestWaitForOffer Error with SetRemoteDescription", err)
		return
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		fmt.Println("TestWaitForOffer Error with CreateAnswer", err)
		return
	}
	answerEncoded, _ := util.EncodeObject(answer)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		fmt.Println("TestWaitForOffer Error with SetLocalDescription", err)
		return
	}

	fmt.Println(answerEncoded)

	// Block forever
	select {}
}
