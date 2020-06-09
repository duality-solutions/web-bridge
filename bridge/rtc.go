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
	keepAlive := true
	if link.PeerConnection == nil {
		fmt.Println("EstablishRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	fmt.Println("EstablishRTC found answer!", link.LinkAccount, link.LinkID())
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if connectionState.String() == "disconnected" {
			keepAlive = false
		}
	})

	// Register channel opening handling
	link.DataChannel.OnOpen(func() {
		fmt.Printf("EstablishRTC Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 30 seconds\n", link.DataChannel.Label(), link.DataChannel.ID())

		for range time.NewTicker(30 * time.Second).C {
			if !keepAlive {
				break
			}
			rand, _ := util.RandomString(7)
			message := "From " + link.MyAccount + " to " + link.LinkAccount + " :" + rand
			fmt.Printf("EstablishRTC Sending '%s'\n", message)

			// Send the message as text
			sendErr := link.DataChannel.SendText(message)
			if sendErr != nil {
				fmt.Printf("EstablishRTC SendText error: %s\n", sendErr)
			}
		}
	})

	// Register text message handling
	link.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("EstablishRTC Message from DataChannel '%s': '%s'\n", link.DataChannel.Label(), string(msg.Data))
	})

	link.DataChannel.OnError(func(err error) {
		fmt.Printf("EstablishRTC DataChannel OnError '%s': '%s'\n", link.DataChannel.Label(), err.Error())
		keepAlive = false
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
	for keepAlive {
		select {}
	}
	fmt.Println("EstablishRTC stopped", link.LinkParticipants())
}

// WaitForRTC waits for a real time connection (RTC) bridge with the link
func WaitForRTC(link *Bridge, answer webrtc.SessionDescription) {
	keepAlive := true
	if link.PeerConnection == nil {
		fmt.Println("WaitForRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	fmt.Println("WaitForRTC created answer!", link.LinkAccount, link.LinkID())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("WaitForRTC ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if connectionState.String() == "disconnected" {
			keepAlive = false
		}
	})

	// Register data channel creation handling
	link.PeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("WaitForRTC New DataChannel %s %d\n", d.Label(), d.ID())
		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("WaitForRTC Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 30 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(30 * time.Second).C {
				if !keepAlive {
					break
				}
				rand, _ := util.RandomString(7)
				message := "From " + link.MyAccount + " to " + link.LinkAccount + " :" + rand
				fmt.Printf("WaitForRTC Sending '%s'\n", message)

				// Send the message as text
				sendErr := d.SendText(message)
				if sendErr != nil {
					fmt.Printf("WaitForRTC SendText error: %s\n", sendErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("WaitForRTC Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})

		d.OnError(func(err error) {
			fmt.Printf("WaitForRTC DataChannel OnError '%s': '%s'\n", link.DataChannel.Label(), err.Error())
			keepAlive = false
		})
	})

	// Set the local SessionDescription
	err := link.PeerConnection.SetLocalDescription(answer)
	if err != nil {
		fmt.Println("WaitForRTC SetLocalDescription error ", link.LinkParticipants(), err)
	}
	fmt.Println("WaitForRTC SetLocalDescription", link.LinkAccount)

	for keepAlive {
		select {}
	}
	fmt.Println("WaitForRTC stopped", link.LinkParticipants())
}
