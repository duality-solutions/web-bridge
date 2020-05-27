package bridge

import (
	"fmt"
	"strings"
	"time"

	"github.com/duality-solutions/web-bridge/internal/dynamic"
	"github.com/duality-solutions/web-bridge/internal/settings"
	"github.com/pion/webrtc/v2"
)

// Bridge hold information about a link WebRTC bridge connection
type Bridge struct {
	MyAccount      string
	LinkAccount    string
	Offer          string
	Answer         string
	LastPingEpoch  int
	PeerConnection *webrtc.PeerConnection
	DataChannel    *webrtc.DataChannel
}

// Bridges hold all link WebRTC bridges
type Bridges struct {
	connected   []Bridge
	unconnected []Bridge
}

var linkBridges Bridges
var dynamicd dynamic.Dynamicd
var config settings.Configuration
var accounts []dynamic.Account
var links dynamic.ActiveLinks

// getAllOffers checks the DHT for WebRTC offers from all links
func getAllOffers() {
	getOffers := make(chan dynamic.DHTGetReturn, len(links.Links))
	for _, link := range links.Links {
		var linkBridge = newBridge(link, accounts)
		dynamicd.GetLinkRecord(linkBridge.LinkAccount, linkBridge.MyAccount, getOffers)
	}
	for i := 0; i < len(links.Links); i++ {
		offer := <-getOffers
		if offer.GetValueSize > 0 {
			linkBridge := newLinkBridge(offer.Sender, offer.Receiver, accounts)
			pc, err := ConnectToIceServices(config)
			if err == nil && offer.GetValue != "null" {
				//fmt.Println("Offer found for", offer.Sender)
				linkBridge.Offer = strings.ReplaceAll(offer.GetValue, `""`, "")
				linkBridge.PeerConnection = pc
				linkBridges.connected = append(linkBridges.connected, linkBridge)
			} else {
				//fmt.Println("Offer found for", offer.Sender, "ConnectToIceServices failed", err)
				linkBridges.unconnected = append(linkBridges.unconnected, linkBridge)
			}
		} else {
			linkBridge := newLinkBridge(offer.Sender, offer.Receiver, accounts)
			pc, _ := ConnectToIceServices(config)
			linkBridge.PeerConnection = pc
			linkBridges.unconnected = append(linkBridges.unconnected, linkBridge)
		}
	}
}

func putOffers(bridges []Bridge) {
	putOffers := make(chan dynamic.DHTPutReturn, len(bridges))
	for _, link := range bridges {
		var linkBridge = newLinkBridge(link.LinkAccount, link.MyAccount, accounts)
		offer, _ := link.PeerConnection.CreateOffer(nil)
		dynamicd.PutLinkRecord(linkBridge.MyAccount, linkBridge.LinkAccount, offer.SDP, putOffers)
	}
	for i := 0; i < len(linkBridges.unconnected); i++ {
		offer := <-putOffers
		fmt.Println("Offer saved", offer)
	}
}

func sendAnswers(bridges []Bridge) {
	for _, brd := range bridges {
		if brd.PeerConnection != nil && len(brd.Offer) > 10 {
			brd.Offer = strings.ReplaceAll(brd.Offer, `""`, "") // remove double quotes in offer
			sd := webrtc.SessionDescription{Type: 1, SDP: brd.Offer}
			//fmt.Println("sendAnswers", brd.LinkAccount, "SessionDescription", sd)
			err := brd.PeerConnection.SetRemoteDescription(sd)
			if err != nil {
				// move to unconnected
				fmt.Printf("sendAnswers failed to connect to link %s. Error %s\n", brd.LinkAccount, err)
			} else {
				answer, err := brd.PeerConnection.CreateAnswer(nil)
				if err != nil {
					fmt.Println(brd.LinkAccount, "CreateAnswer error", err)
					// clear offer since it didn't work
					// remove from connected and add to unconnected
				} else {
					//fmt.Println("SendLinkMessage", brd.LinkAccount, answer.SDP)
					_, err := dynamicd.SendLinkMessage(brd.MyAccount, brd.LinkAccount, answer.SDP)
					if err != nil {
						fmt.Println("SendLinkMessage error", brd.LinkAccount, err)
					}
				}
			}
		} else {
			fmt.Println("Error nil PeerConnection", brd.LinkAccount)
		}
	}
}

func clearOffers(bridges []Bridge) {
	clearOffers := make(chan dynamic.DHTPutReturn, len(bridges))
	for _, link := range bridges {
		var linkBridge = newLinkBridge(link.LinkAccount, link.MyAccount, accounts)
		dynamicd.ClearLinkRecord(linkBridge.MyAccount, linkBridge.LinkAccount, clearOffers)
	}
	for i := 0; i < len(bridges); i++ {
		offer := <-clearOffers
		fmt.Println("Offer cleared", offer)
	}
}

func waitForSync() {
	status, _ := dynamicd.GetSyncStatus()
	for status.SyncProgress < 1 {
		time.Sleep(time.Second * 30)
		status, _ = dynamicd.GetSyncStatus()
	}
	time.Sleep(time.Second * 10)
}

// StartBridges runs a goroutine in the background to manage network bridges
// get link offers from DHT
// send answers to offers using VGP instant messaging and create a WebRTC bridge
// send bridge result to upstream channel
// put offers in the DHT for unconnected links
// check for answers loop
// if new answer found, create a WebRTC bridge and send bridge result to upstream channel
// on shutdown, clear all DHT offers
func StartBridges(chanBridge *chan []Bridge, c settings.Configuration, d dynamic.Dynamicd, a []dynamic.Account, l dynamic.ActiveLinks) {
	dynamicd = d
	config = c
	accounts = a
	links = l
	waitForSync()
	fmt.Println("\n\nStarting Link Bridges")
	// check all links for WebRTC offers in the DHT
	getAllOffers()
	fmt.Println("Get offers complete.", len(linkBridges.connected), len(linkBridges.unconnected))
	// respond to all offers with a WebRTC answer and send it to the link using instant VGP messages
	sendAnswers(linkBridges.connected)
	fmt.Println("Send answers complete.", len(linkBridges.connected))
	// put WebRTC offers for unconnected links
	putOffers(linkBridges.unconnected)
	fmt.Println("Put offers complete.", len(linkBridges.unconnected))
}

// ShutdownBridges stops the ManageBridges goroutine
func ShutdownBridges() {
	//TODO: disconnect WebRTC bridges
	//clear all link offers in the DHT
	clearOffers(linkBridges.unconnected)
	clearOffers(linkBridges.connected)
	// sleep for 20 seconds to make sure all clear take effect.
	time.Sleep(time.Second * 20)
}

func newBridge(l dynamic.Link, acc []dynamic.Account) Bridge {
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

func newLinkBridge(requester string, recipient string, acc []dynamic.Account) Bridge {
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
