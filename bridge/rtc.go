package bridge

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

import (
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/pion/webrtc/v2"
)

// EstablishRTC tries to establish a real time connection (RTC) bridge with the link
func EstablishRTC(link *Bridge) {
	stopchan := make(chan struct{})
	if link.PeerConnection() == nil {
		util.Info.Println("EstablishRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	util.Info.Println("EstablishRTC found answer!", link.LinkAccount, link.LinkID())
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection().OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		link.SetOnStateChangeEpoch(time.Now().Unix())
		util.Info.Printf("EstablishRTC ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if link.RTCState() == "checking" && connectionState.String() == "failed" {
			close(stopchan)
		}
		link.SetRTCState(connectionState.String())
		if connectionState.String() == "disconnected" {
			link.ShutdownHTTPProxyServers()
			close(stopchan)
		}
	})

	// Register channel opening handling
	link.DataChannel().OnOpen(func() {
		if link.DataChannel() != nil {
			link.SetOnOpenEpoch(time.Now().Unix())
			link.SetState(StateOpenConnection)
			util.Info.Printf("EstablishRTC Data channel '%s'-'%d' open.\n", link.DataChannel().Label(), link.DataChannel().ID())
			// Detach the data channel
			raw, err := link.DataChannel().Detach()
			if err != nil {
				util.Error.Println("EstablishRTC link DataChannel OnOpen error", err)
				close(stopchan)
			} else {
				go link.StartBridgeNetwork(raw, raw)
			}
		} else {
			util.Error.Println("EstablishRTC link.DataChannel().OnOpen nil DataChannel. Stopping")
			close(stopchan)
		}
	})

	link.DataChannel().OnError(func(err error) {
		if link.DataChannel() != nil {
			link.SetOnErrorEpoch(time.Now().Unix())
			util.Error.Printf("EstablishRTC DataChannel OnError '%s': '%s'\n", link.DataChannel().Label(), err.Error())
			close(stopchan)
		} else {
			link.SetOnErrorEpoch(time.Now().Unix())
			util.Error.Printf("EstablishRTC DataChannel.OnError DataChannel is nil\n")
			close(stopchan)
		}
	})

	// Set the local SessionDescription
	err := link.PeerConnection().SetLocalDescription(link.Offer())
	if err != nil {
		util.Error.Println("EstablishRTC error SetLocalDescription", link.LinkParticipants(), err)
	}

	// Set the remote SessionDescription
	err = link.PeerConnection().SetRemoteDescription(link.Answer())
	if err != nil {
		util.Error.Println("EstablishRTC SetRemoteDescription error ", link.LinkParticipants(), err)
	}
	util.Info.Println("EstablishRTC SetRemoteDescription", link.LinkAccount)

	for true {
		select {
		default:
		case <-stopchan:
			break
		}
	}
	if link.DataChannel() != nil {
		link.SetDataChannel(nil)
	}
	for link.PeerConnection().ICEConnectionState().String() != "failed" {
		time.Sleep(10 * time.Second)
	}
	failedICEConnection := (link.PeerConnection().ICEConnectionState().String() == "failed")
	if failedICEConnection {
		util.Info.Println("EstablishRTC close peer connection", link.LinkParticipants(), link.LinkID())
		link.PeerConnection().Close()
		link.SetOnErrorEpoch(time.Now().Unix())
	}
	bridgeControler.MoveConnectedToUnconnected(link)
	util.Info.Println("EstablishRTC stopped!", link.LinkParticipants())
	link.SetState(StateInit)

	return
}

// WaitForRTC waits for a real time connection (RTC) bridge with the link
// TODO: add timeout
func WaitForRTC(link *Bridge) {
	stopchan := make(chan struct{})
	if link.PeerConnection() == nil {
		util.Error.Println("WaitForRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	util.Info.Println("WaitForRTC created answer!", link.LinkAccount, link.LinkID())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection().OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		link.SetOnStateChangeEpoch(time.Now().Unix())
		util.Info.Printf("WaitForRTC ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if link.RTCState() == "checking" && connectionState.String() == "failed" {
			close(stopchan)
		}
		link.SetRTCState(connectionState.String())
		if connectionState.String() == "disconnected" {
			link.ShutdownHTTPProxyServers()
			close(stopchan)
		}
	})

	// Register data channel creation handling
	link.PeerConnection().OnDataChannel(func(d *webrtc.DataChannel) {
		link.SetDataChannel(d)
		util.Info.Printf("WaitForRTC New DataChannel %s %d\n", link.DataChannel().Label(), link.DataChannel().ID())
		// Register channel opening handling
		link.DataChannel().OnOpen(func() {
			if link.DataChannel() != nil {
				link.SetOnOpenEpoch(time.Now().Unix())
				link.SetState(StateOpenConnection)
				util.Info.Printf("WaitForRTC Data channel '%s'-'%d' open.\n", link.DataChannel().Label(), link.DataChannel().ID())
				// Detach the data channel
				raw, err := link.DataChannel().Detach()
				if err != nil {
					util.Error.Println("WaitForRTC link DataChannel OnOpen error", err)
					close(stopchan)
				} else {
					go link.StartBridgeNetwork(raw, raw)
				}
			} else {
				util.Error.Println("WaitForRTC link.DataChannel().OnOpen nil DataChannel. Stopping")
				close(stopchan)
			}
		})

		link.DataChannel().OnError(func(err error) {
			if link.DataChannel() != nil {
				link.SetOnErrorEpoch(time.Now().Unix())
				util.Error.Printf("WaitForRTC DataChannel OnError '%s': '%s'\n", link.DataChannel().Label(), err.Error())
				close(stopchan)
			} else {
				link.SetOnErrorEpoch(time.Now().Unix())
				util.Error.Printf("WaitForRTC DataChannel.OnError DataChannel is nil\n")
				close(stopchan)
			}
		})
	})

	// Set the local SessionDescription
	err := link.PeerConnection().SetLocalDescription(link.Answer())
	if err != nil {
		util.Error.Println("WaitForRTC SetLocalDescription error ", link.LinkParticipants(), err)
		close(stopchan)
	}
	util.Info.Println("WaitForRTC SetLocalDescription", link.LinkAccount)

	for true {
		select {
		default:
		case <-stopchan:
			break
		}
	}
	if link.DataChannel() != nil {
		link.SetDataChannel(nil)
	}
	for link.PeerConnection().ICEConnectionState().String() != "failed" {
		time.Sleep(10 * time.Second)
	}
	failedICEConnection := (link.PeerConnection().ICEConnectionState().String() == "failed")
	if failedICEConnection {
		util.Info.Println("WaitForRTC close peer connection", link.LinkParticipants(), link.LinkID())
		link.PeerConnection().Close()
		link.SetOnErrorEpoch(time.Now().Unix())
	} else {
		util.Info.Println("WaitForRTC ICEConnectionState", link.LinkParticipants(), link.PeerConnection().ICEConnectionState().String())
	}
	util.Info.Println("WaitForRTC stopped!", link.LinkParticipants())
	link.SetState(StateInit)
	bridgeControler.MoveConnectedToUnconnected(link)
	return
}
