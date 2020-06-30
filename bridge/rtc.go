package bridge

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
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

// ResponseLoop shows how to read from the datachannel directly
func ResponseLoop(r io.Reader, w io.Writer) {
	for {
		var buffer []byte
		n, err := r.Read(buffer)
		if n > 10 {
			util.Info.Println("Message from DataChannel:", string(buffer[:n]))
			if err != nil {
				util.Info.Println("Datachannel Read error", err)
				continue
			}
			buf := bytes.NewReader(buffer)
			buf2 := bufio.NewReader(buf)
			var req *http.Request
			if req, err = http.ReadRequest(buf2); err != nil { // deserialize request
				util.Error.Println("Datachannel deserialize error", err)
				continue
			}
			destConn, err2 := net.DialTimeout("tcp", req.Host, 10*time.Second)
			if err2 != nil {
				util.Error.Println("Datachannel DialTimeout error", err)
				continue
			}
			var buffer2 []byte
			destConn.Read(buffer2)
			w.Write(buffer2)
		}
	}
}

// EstablishRTC tries to establish a real time connection (RTC) bridge with the link
func EstablishRTC(link *Bridge) {
	keepAlive := true
	stopchan := make(chan struct{})
	if link.PeerConnection == nil {
		util.Info.Println("EstablishRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	util.Info.Println("EstablishRTC found answer!", link.LinkAccount, link.LinkID())
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		link.OnStateChangeEpoch = time.Now().Unix()
		link.RTCState = connectionState.String()
		util.Info.Printf("EstablishRTC ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if connectionState.String() == "disconnected" {
			keepAlive = false
			close(stopchan)
		}
	})

	// Register channel opening handling
	link.DataChannel.OnOpen(func() {
		link.OnOpenEpoch = time.Now().Unix()
		link.State = StateOpenConnection
		util.Info.Printf("EstablishRTC Data channel '%s'-'%d' open.\n", link.DataChannel.Label(), link.DataChannel.ID())
		raw, dErr := link.DataChannel.Detach()
		if dErr != nil {
			util.Error.Println("EstablishRTC Data channel error", link.LinkParticipants(), dErr)
		}
		link.ReadWriteCloser = raw
		link.StartBridgeNetwork()
		go ResponseLoop(raw, raw)
	})

	// Register text message handling
	link.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		if link.DataChannel != nil {
			link.LastDataEpoch = time.Now().Unix()
			util.Info.Printf("EstablishRTC Message from DataChannel '%s': '%s'\n", link.DataChannel.Label(), string(msg.Data))
		}
	})

	link.DataChannel.OnError(func(err error) {
		link.OnErrorEpoch = time.Now().Unix()
		util.Error.Printf("EstablishRTC DataChannel OnError '%s': '%s'\n", link.DataChannel.Label(), err.Error())
		keepAlive = false
		close(stopchan)
	})

	// Set the local SessionDescription
	err := link.PeerConnection.SetLocalDescription(link.Offer)
	if err != nil {
		util.Error.Println("EstablishRTC error SetLocalDescription", link.LinkParticipants(), err)
	}

	// Set the remote SessionDescription
	err = link.PeerConnection.SetRemoteDescription(link.Answer)
	if err != nil {
		util.Error.Println("EstablishRTC SetRemoteDescription error ", link.LinkParticipants(), err)
	}
	util.Info.Println("EstablishRTC SetRemoteDescription", link.LinkAccount)

	for keepAlive {
		select {
		default:
			if !keepAlive {
				break
			}
		case <-stopchan:
			break
		}
	}
	if link.DataChannel != nil {
		link.DataChannel = nil
	}
	failedICEConnection := (link.PeerConnection.ICEConnectionState().String() == "failed")
	if failedICEConnection {
		util.Info.Println("EstablishRTC close peer connection", link.LinkParticipants(), link.LinkID())
		link.PeerConnection.Close()
	}
	delete(linkBridges.connected, link.LinkID())
	linkBridges.unconnected[link.LinkID()] = link
	util.Info.Println("EstablishRTC stopped!", link.LinkParticipants())
	link.State = StateInit
}

// WaitForRTC waits for a real time connection (RTC) bridge with the link
func WaitForRTC(link *Bridge, answer webrtc.SessionDescription) {
	keepAlive := true
	stopchan := make(chan struct{})
	if link.PeerConnection == nil {
		util.Info.Println("WaitForRTC PeerConnection nil for", link.LinkAccount)
		return
	}
	util.Info.Println("WaitForRTC created answer!", link.LinkAccount, link.LinkID())

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		link.OnStateChangeEpoch = time.Now().Unix()
		link.RTCState = connectionState.String()
		util.Info.Printf("WaitForRTC ICE Connection State has changed for %s: %s\n", link.LinkParticipants(), connectionState.String())
		if connectionState.String() == "disconnected" {
			keepAlive = false
			close(stopchan)
		}
	})

	// Register data channel creation handling
	link.PeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		link.DataChannel = d
		util.Info.Printf("WaitForRTC New DataChannel %s %d\n", link.DataChannel.Label(), link.DataChannel.ID())
		// Register channel opening handling
		link.DataChannel.OnOpen(func() {
			link.OnOpenEpoch = time.Now().Unix()
			link.State = StateOpenConnection
			util.Info.Printf("WaitForRTC Data channel '%s'-'%d' open.\n", link.DataChannel.Label(), link.DataChannel.ID())
			raw, dErr := link.DataChannel.Detach()
			if dErr != nil {
				util.Error.Println("WaitForRTC Data channel error", link.LinkParticipants(), dErr)
			}
			link.ReadWriteCloser = raw
			link.StartBridgeNetwork()
			go ResponseLoop(raw, raw)
		})

		// Register text message handling
		link.DataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			if link.DataChannel != nil {
				link.LastDataEpoch = time.Now().Unix()
				util.Info.Printf("WaitForRTC Message from DataChannel '%s': '%s'\n", link.DataChannel.Label(), string(msg.Data))
			}
		})

		link.DataChannel.OnError(func(err error) {
			link.OnErrorEpoch = time.Now().Unix()
			util.Error.Printf("WaitForRTC DataChannel OnError '%s': '%s'\n", link.DataChannel.Label(), err.Error())
			keepAlive = false
			close(stopchan)
		})
	})

	// Set the local SessionDescription
	err := link.PeerConnection.SetLocalDescription(answer)
	if err != nil {
		util.Error.Println("WaitForRTC SetLocalDescription error ", link.LinkParticipants(), err)
	}
	util.Info.Println("WaitForRTC SetLocalDescription", link.LinkAccount)

	for keepAlive {
		select {
		default:
			if !keepAlive {
				break
			}
		case <-stopchan:
			break
		}
	}
	if link.DataChannel != nil {
		link.DataChannel = nil
	}
	failedICEConnection := (link.PeerConnection.ICEConnectionState().String() == "failed")
	if failedICEConnection {
		util.Info.Println("WaitForRTC close peer connection", link.LinkParticipants(), link.LinkID())
		link.PeerConnection.Close()
	} else {
		util.Info.Println("WaitForRTC ICEConnectionState", link.LinkParticipants(), link.PeerConnection.ICEConnectionState().String())
	}
	delete(linkBridges.connected, link.LinkID())
	linkBridges.unconnected[link.LinkID()] = link
	util.Info.Println("WaitForRTC stopped!", link.LinkParticipants())
	link.State = StateInit
}
