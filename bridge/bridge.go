package bridge

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"sort"
	"sync"

	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/pion/webrtc/v2"
)

// Bridge hold information about a link WebRTC bridge connection
type Bridge struct {
	SessionID          uint16
	MyAccount          string
	LinkAccount        string
	bridgeMut          *sync.RWMutex
	offer              webrtc.SessionDescription
	answer             webrtc.SessionDescription
	onOpenEpoch        int64
	onErrorEpoch       int64
	onStateChangeEpoch int64
	rtcState           string
	onLastDataEpoch    int64
	peerConnection     *webrtc.PeerConnection
	dataChannel        *webrtc.DataChannel
	proxyHTTP          *http.Server
	proxyHTTPS         *net.Listener
	state              State
}

// NewBridge creates a new bridge struct
func NewBridge(s uint16, l dynamic.Link, acc []dynamic.Account) Bridge {
	var brd Bridge
	brd.SessionID = s
	brd.bridgeMut = new(sync.RWMutex)
	brd.state = StateNew
	for _, a := range acc {
		if a.ObjectID == l.GetRequestorObjectID() {
			brd.MyAccount = l.GetRequestorObjectID()
			brd.LinkAccount = l.GetRecipientObjectID()
			return brd
		} else if a.ObjectID == l.GetRecipientObjectID() {
			brd.MyAccount = l.GetRecipientObjectID()
			brd.LinkAccount = l.GetRequestorObjectID()
			return brd
		}
	}
	return brd
}

// ResetBridge clones an existing bridge into a new struct
func ResetBridge(b *Bridge) Bridge {
	var newBridge Bridge
	newBridge.bridgeMut = new(sync.RWMutex)
	newBridge.SessionID = b.SessionID
	newBridge.MyAccount = b.MyAccount
	newBridge.LinkAccount = b.LinkAccount
	newBridge.SetState(StateNew)
	newBridge.SetRTCState("")
	return newBridge
}

// NewLinkBridge creates a new bridge struct
func NewLinkBridge(s uint16, requester string, recipient string, acc []dynamic.Account) Bridge {
	var brd Bridge
	brd.SessionID = s
	brd.bridgeMut = new(sync.RWMutex)
	brd.state = StateNew
	for _, a := range acc {
		if a.ObjectID == requester {
			brd.MyAccount = requester
			brd.LinkAccount = recipient
			return brd
		} else if a.ObjectID == recipient {
			brd.MyAccount = recipient
			brd.LinkAccount = requester
			return brd
		}
	}
	return brd
}

// SetOfferAnswerStateEpoch sets the bridge WebRTC offer, answer, state and onStateChangeEpoch
func (b *Bridge) SetOfferAnswerStateEpoch(o webrtc.SessionDescription, a webrtc.SessionDescription, s State, e int64) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.offer = o
	b.answer = a
	b.state = s
	b.onStateChangeEpoch = e
}

// SetAnswerStateEpoch sets the bridge WebRTC answer, state and onStateChangeEpoch
func (b *Bridge) SetAnswerStateEpoch(a webrtc.SessionDescription, s State, e int64) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.answer = a
	b.state = s
	b.onStateChangeEpoch = e
}

// SetOnLastDataEpoch sets the bridge WebRTC last data epoch time
func (b *Bridge) SetOnLastDataEpoch(e int64) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.onLastDataEpoch = e
}

// OnLastDataEpoch returns the bridge WebRTC last data epoch time
func (b *Bridge) OnLastDataEpoch() int64 {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.onLastDataEpoch
}

// SetState sets the bridge WebRTC state
func (b *Bridge) SetState(s State) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.state = s
}

// State returns the bridge WebRTC state
func (b *Bridge) State() State {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.state
}

// SetOnOpenEpoch sets the bridge WebRTC open epoch time
func (b *Bridge) SetOnOpenEpoch(e int64) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.onOpenEpoch = e
}

// OnOpenEpoch returns the bridge WebRTC open epoch time
func (b *Bridge) OnOpenEpoch() int64 {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.onOpenEpoch
}

// SetOnErrorEpoch sets the bridge WebRTC error epoch time
func (b *Bridge) SetOnErrorEpoch(e int64) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.onErrorEpoch = e
}

// OnErrorEpoch sets the bridge WebRTC error epoch time
func (b *Bridge) OnErrorEpoch() int64 {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	return b.onErrorEpoch
}

// SetOnStateChangeEpoch returns the bridge WebRTC set on change epoch time
func (b *Bridge) SetOnStateChangeEpoch(e int64) {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	b.onStateChangeEpoch = e
}

// OnStateChangeEpoch returns the bridge WebRTC set on change epoch time
func (b *Bridge) OnStateChangeEpoch() int64 {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.onStateChangeEpoch
}

// SetRTCState sets the bridge WebRTC state
func (b *Bridge) SetRTCState(s string) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.rtcState = s
}

// RTCState returns the bridge WebRTC state
func (b *Bridge) RTCState() string {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.rtcState
}

// SetDataChannel sets the bridge WebRTC data channel struct pointer
func (b *Bridge) SetDataChannel(dc *webrtc.DataChannel) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.dataChannel = dc
}

// DataChannel returns the bridge WebRTC peer data channel struct pointer
func (b *Bridge) DataChannel() *webrtc.DataChannel {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.dataChannel
}

// SetPeerConnection sets the bridge WebRTC peer connection struct pointer
func (b *Bridge) SetPeerConnection(pc *webrtc.PeerConnection) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.peerConnection = pc
}

// PeerConnection returns the bridge WebRTC peer connection struct pointer
func (b *Bridge) PeerConnection() *webrtc.PeerConnection {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.peerConnection
}

// SetOffer sets the bridge offer variable
func (b *Bridge) SetOffer(o webrtc.SessionDescription) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.offer = o
}

// Offer returns the bridge offer variable
func (b *Bridge) Offer() webrtc.SessionDescription {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.offer
}

// SetAnswer sets the bridge answer variable
func (b *Bridge) SetAnswer(o webrtc.SessionDescription) {
	b.bridgeMut.Lock()
	defer b.bridgeMut.Unlock()
	b.answer = o
}

// Answer returns the bridge answer variable
func (b *Bridge) Answer() webrtc.SessionDescription {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	return b.answer
}

// LinkID returns an hashed id for the link
func (b *Bridge) LinkID() string {
	var ret string = ""
	strs := []string{b.MyAccount, b.LinkAccount}
	sort.Strings(strs)
	for _, str := range strs {
		ret += str
	}
	hash := sha256.New()
	hash.Write([]byte(ret))
	bs := hash.Sum(nil)
	hs := fmt.Sprintf("%x", bs)
	return hs
}

// ListenPort returns the HTTP server listening port
func (b *Bridge) ListenPort() uint16 {
	return b.SessionID + StartHTTPPortNumber
}

// LinkParticipants returns link participants
func (b *Bridge) LinkParticipants() string {
	return (b.MyAccount + "-" + b.LinkAccount)
}

// Shutdown safely closes WebRTC and proxy networks
func (b *Bridge) Shutdown() {
	if b.peerConnection.ConnectionState().String() == "connected" {
		b.dataChannel.Close()
		b.peerConnection.Close()
	}
	b.ShutdownHTTPProxyServers()
}

// ShutdownHTTPProxyServers closes the proxy server and sets bridge to disconnected
func (b *Bridge) ShutdownHTTPProxyServers() {
	b.SetState(StateDisconnected)
	if b.proxyHTTP != nil {
		util.Info.Println("ShutdownHTTPProxyServers http proxyHTTP")
		b.proxyHTTP.Shutdown(context.Background())
	}
	if b.proxyHTTPS != nil {
		listen := *b.proxyHTTPS
		err := listen.Close()
		util.Info.Println("ShutdownHTTPProxyServers listener proxyHTTPS")
		if err != nil {
			util.Error.Println("proxyHTTPS close error", err)
		}
	}
}

// String returns a string representation of the bridge struct
func (b *Bridge) String() string {
	b.bridgeMut.RLock()
	defer b.bridgeMut.RUnlock()
	result := fmt.Sprint("Bridge {",
		"\nSessionID: ", b.SessionID,
		"\nListenPort: ", b.ListenPort(),
		"\nMyAccount: ", b.MyAccount,
		"\nLinkAccount: ", b.LinkAccount,
		"\nLinkID: ", b.LinkID(),
		"\nOffer: ", b.offer.SDP,
		"\nAnswer: ", b.answer.SDP,
		"\nOnOpenEpoch: ", b.onOpenEpoch,
		"\nOnErrorEpoch: ", b.onErrorEpoch,
		"\nOnStateChangeEpoch: ", b.onStateChangeEpoch,
		"\nRTCState: ", b.rtcState,
		"\nOnLastDataEpoch: ", b.onLastDataEpoch,
		"\nState: ", b.state.String(),
	)
	if b.peerConnection != nil {
		result += fmt.Sprint("\nPeerConnection.ICEConnectionState: ", b.peerConnection.ICEConnectionState().String(),
			"\nPeerConnection.ConnectionState: ", b.peerConnection.ConnectionState().String(),
		)
	} else {
		result += fmt.Sprint("\nPeerConnection.ICEConnectionState: nil\nPeerConnection.ConnectionState: nil")
	}
	if b.dataChannel != nil {
		result += fmt.Sprint("\nDataChannel.ReadyState: ", b.dataChannel.ReadyState().String())
	} else {
		result += fmt.Sprint("\nDataChannel.ReadyState: nil")
	}
	result += "\n}"
	return result
}
