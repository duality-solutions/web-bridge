package bridge

import "github.com/duality-solutions/web-bridge/internal/dynamic"

// NewBridge creates a new bridge struct
func NewBridge(l dynamic.Link, acc []dynamic.Account) Bridge {
	var brd Bridge
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
