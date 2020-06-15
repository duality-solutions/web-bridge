package bridge

import (
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
	// StateWaitForAnswer is waiting for offer = 2
	StateWaitForAnswer
	// StateSendAnswer offer received, send answer = 3
	StateSendAnswer
	// StateWaitForRTC offer received and answer sent  = 4
	StateWaitForRTC
	// StateEstablishRTC offer sent and answer received  = 5
	StateEstablishRTC
)

// Bridge hold information about a link WebRTC bridge connection
type Bridge struct {
	MyAccount   string
	LinkAccount string
	Offer       webrtc.SessionDescription
	Answer      webrtc.SessionDescription
	State
	LastPingEpoch  int
	PeerConnection *webrtc.PeerConnection
	DataChannel    *webrtc.DataChannel
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
