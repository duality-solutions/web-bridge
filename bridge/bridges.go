package bridge

import (
	"fmt"
	"time"

	"github.com/duality-solutions/web-bridge/init/settings"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
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
	connected   []*Bridge
	unconnected []*Bridge
}

var linkBridges Bridges
var dynamicd dynamic.Dynamicd
var config settings.Configuration
var accounts []dynamic.Account
var links dynamic.ActiveLinks

func initializeBridges(stopchan chan struct{}) bool {
	// check all links for WebRTC offers in the DHT
	if GetAllOffers(stopchan, links, accounts) {
		fmt.Println("Get all offers complete. Found", len(linkBridges.connected), "Not found", len(linkBridges.unconnected))
		// respond to all offers with a WebRTC answer and send it to the link using instant VGP messages
		if SendAnswers(stopchan) {
			fmt.Println("Send answers completed", len(linkBridges.connected))
			// put WebRTC offers for unconnected links
			if PutOffers(stopchan) {
				fmt.Println("Put offers completed", len(linkBridges.unconnected))
			} else {
				fmt.Println("StartBridges stopped after PutOffers")
				return false
			}
		} else {
			fmt.Println("StartBridges stopped after SendAnswers")
			return false
		}
	} else {
		fmt.Println("StartBridges stopped after GetAllOffers")
		return false
	}
	return true
}

// StartBridges runs a goroutine to manage network bridges
// get link offers from DHT
// send answers to offers using VGP instant messaging
// send bridge result to upstream channel
// put offers in the DHT for unconnected links
// check for answers loop
// if new answer found, create a WebRTC bridge and send bridge result to upstream channel
// on shutdown, clear all DHT offers
func StartBridges(stopchan chan struct{}, c settings.Configuration, d dynamic.Dynamicd, a []dynamic.Account, l dynamic.ActiveLinks) {
	dynamicd = d
	config = c
	accounts = a
	links = l
	if dynamicd.WaitForSync(stopchan, 10, 10) {
		if initializeBridges(stopchan) {
			GetAnswers(stopchan)
			fmt.Println("StartBridges stopped after GetAnswers")
		}
	} else {
		fmt.Println("StartBridges stopped after WaitForSync")
	}
}

// ShutdownBridges stops the ManageBridges goroutine
func ShutdownBridges() {
	//TODO: disconnect WebRTC bridges
	//clear all link offers in the DHT
	ClearOffers()
	// sleep for 20 seconds to make sure all clear take effect.
	time.Sleep(time.Second * 20)
}
