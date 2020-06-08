package bridge

import (
	"fmt"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/pion/webrtc/v2"
)

/*
WebRTC:

Offer Peer (EstablishRTC)					 Answer Peer (WaitForRTC)
1- NewPeerConnection                         1- NewPeerConnection
2- CreateDataChannel                         2- OnICEConnectionStateChange
3- OnICEConnectionStateChange                3- OnDataChannel
4- dataChannel.OnOpen                        4- SetRemoteDescription
5- dataChannel.OnMessage                     5- CreateAnswer
6- CreateOffer                               6- SetLocalDescription
7- SetLocalDescription
8- SetRemoteDescription
*/

// EstablishRTC tries to establish a real time connection (RTC) bridge with the link
func EstablishRTC(link *Bridge) {
	if link.PeerConnection == nil {
		fmt.Println("EstablishRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	fmt.Println("EstablishRTC found answer!", link.LinkAccount, link.LinkID())
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
	})

	// Register channel opening handling
	link.DataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", link.DataChannel.Label(), link.DataChannel.ID())

		for range time.NewTicker(30 * time.Second).C {
			rand, _ := util.RandomString(7)
			message := "From " + link.MyAccount + " to " + link.LinkAccount + " :" + rand
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendErr := link.DataChannel.SendText(message)
			if sendErr != nil {
				fmt.Printf("SendText error: %s\n", sendErr)
			}
		}
	})

	// Register text message handling
	link.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", link.DataChannel.Label(), string(msg.Data))
	})

	// Set the local SessionDescription
	err := link.PeerConnection.SetLocalDescription(link.Offer)
	if err != nil {
		fmt.Println("EstablishRTC error SetLocalDescription", link.LinkParticipants(), err)
	}

	// Set the remote SessionDescription
	err = link.PeerConnection.SetRemoteDescription(link.Answer)
	if err != nil {
		fmt.Println("EstablishRTC SetRemoteDescription error ", link.LinkParticipants(), err)
	}
	fmt.Println("EstablishRTC SetRemoteDescription", link.LinkAccount)
	// Block forever
	select {}
}

// WaitForRTC waits for a real time connection (RTC) bridge with the link
func WaitForRTC(link *Bridge, answer webrtc.SessionDescription) {
	if link.PeerConnection == nil {
		fmt.Println("EstablishRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	fmt.Println("WaitForRTC created answer!", link.LinkAccount, link.LinkID())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
	})

	// Register data channel creation handling
	link.PeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(30 * time.Second).C {
				rand, _ := util.RandomString(7)
				message := "From " + link.MyAccount + " to " + link.LinkAccount + " :" + rand
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					fmt.Printf("SendText error: %s\n", sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Set the local SessionDescription
	err := link.PeerConnection.SetLocalDescription(answer)
	if err != nil {
		fmt.Println("WaitForRTC SetLocalDescription error ", link.LinkParticipants(), err)
	}
	fmt.Println("WaitForRTC SetLocalDescription", link.LinkAccount)
	// Block forever
	select {}
}
