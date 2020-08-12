package bridge

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"sort"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// State enum stores the bridge state
type State uint16

const (
	// StateInit is the initial bridge state = 0
	StateInit State = 0 + iota
	// StateNew is the state after calling new bridge = 1
	StateNew
	// StateWaitForOffer is when waiting for an offer
	StateWaitForOffer
	// StateWaitForAnswer is waiting for an answer = 3
	StateWaitForAnswer
	// StateSendAnswer offer received, send answer = 4
	StateSendAnswer
	// StateWaitForRTC offer received and answer sent  = 5
	StateWaitForRTC
	// StateEstablishRTC offer sent and answer received  = 6
	StateEstablishRTC
	// StateOpenConnection when WebRTC is connected and open = 7
	StateOpenConnection
	// StateDisconnected when WebRTC goes from connected to diconnected and open = 8
	StateDisconnected
)

func (s State) String() string {
	switch s {
	case StateInit:
		return "StateInit"
	case StateNew:
		return "StateNew"
	case StateWaitForAnswer:
		return "StateWaitForAnswer"
	case StateSendAnswer:
		return "StateSendAnswer"
	case StateWaitForRTC:
		return "StateWaitForRTC"
	case StateEstablishRTC:
		return "StateEstablishRTC"
	case StateOpenConnection:
		return "StateOpenConnection"
	default:
		return "Undefined"
	}
}

// Bridge hold information about a link WebRTC bridge connection
type Bridge struct {
	SessionID          uint16
	MyAccount          string
	LinkAccount        string
	Offer              webrtc.SessionDescription
	Answer             webrtc.SessionDescription
	OnOpenEpoch        int64
	OnErrorEpoch       int64
	OnStateChangeEpoch int64
	RTCState           string
	LastDataEpoch      int64
	PeerConnection     *webrtc.PeerConnection
	DataChannel        *webrtc.DataChannel
	proxyHTTP          *http.Server
	proxyHTTPS         *net.Listener
	Get                dynamic.DHTGetJSON
	Put                dynamic.DHTPutJSON
	State
}

// NewBridge creates a new bridge struct
func NewBridge(l dynamic.Link, acc []dynamic.Account) Bridge {
	var brd Bridge
	brd.State = StateNew
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

// NewLinkBridge creates a new bridge struct
func NewLinkBridge(requester string, recipient string, acc []dynamic.Account) Bridge {
	var brd Bridge
	brd.State = StateNew
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

// LinkID returns an hashed id for the link
func (b Bridge) LinkID() string {
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
func (b Bridge) ListenPort() uint16 {
	return b.SessionID + StartHTTPPortNumber
}

// LinkParticipants returns link participants
func (b Bridge) LinkParticipants() string {
	return (b.MyAccount + "-" + b.LinkAccount)
}

// ShutdownHTTPProxyServers returns the HTTP server listening port
func (b *Bridge) ShutdownHTTPProxyServers() {
	b.State = StateDisconnected
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

func (b Bridge) String() string {
	result := fmt.Sprint("Bridge {",
		"\nSessionID: ", b.SessionID,
		"\nListenPort: ", b.ListenPort(),
		"\nMyAccount: ", b.MyAccount,
		"\nLinkAccount: ", b.LinkAccount,
		"\nLinkID: ", b.LinkID(),
		"\nOffer: ", b.Offer.SDP,
		"\nAnswer: ", b.Answer.SDP,
		"\nOnOpenEpoch: ", b.OnOpenEpoch,
		"\nOnErrorEpoch: ", b.OnErrorEpoch,
		"\nOnStateChangeEpoch: ", b.OnStateChangeEpoch,
		"\nRTCStatus: ", b.RTCState,
		"\nLastDataEpoch: ", b.LastDataEpoch,
		"\nState: ", b.State.String(),
	)
	if b.PeerConnection != nil {
		result += fmt.Sprint("\nPeerConnection.ICEConnectionState: ", b.PeerConnection.ICEConnectionState().String(),
			"\nPeerConnection.ConnectionState: ", b.PeerConnection.ConnectionState().String(),
		)
	} else {
		result += fmt.Sprint("\nPeerConnection.ICEConnectionState: nil\nPeerConnection.ConnectionState: nil")
	}
	if b.DataChannel != nil {
		result += fmt.Sprint("\nDataChannel.ReadyState: ", b.DataChannel.ReadyState().String())
	} else {
		result += fmt.Sprint("\nDataChannel.ReadyState: nil")
	}
	result += "\n}"
	return result
}
