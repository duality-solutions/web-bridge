package bridge

import (
	"time"

	"github.com/duality-solutions/web-bridge/init/settings"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
)

var linkBridges Bridges
var dynamicd dynamic.Dynamicd
var config settings.Configuration
var accounts []dynamic.Account
var links dynamic.ActiveLinks

// Bridges hold all link WebRTC bridges
type Bridges struct {
	connected   map[string]*Bridge
	unconnected map[string]*Bridge
}

func initializeBridges(stopchan chan struct{}) bool {
	// check all links for WebRTC offers in the DHT
	if GetAllOffers(stopchan, links, accounts) {
		util.Info.Println("Get all offers complete. unconnected", len(linkBridges.unconnected))
		// respond to all offers with a WebRTC answer and send it to the link using instant VGP messages
		if SendAnswers(stopchan) {
			if GetOffers(stopchan) {
				util.Info.Println("get offers completed", len(linkBridges.unconnected))
			} else {
				return false
			}
			// put WebRTC offers for unconnected links
			if PutOffers(stopchan) {
				util.Info.Println("Put offers completed", len(linkBridges.unconnected))
			} else {
				util.Info.Println("StartBridges stopped after PutOffers")
				return false
			}
		} else {
			util.Info.Println("StartBridges stopped after SendAnswers")
			return false
		}
	} else {
		util.Info.Println("StartBridges stopped after GetAllOffers")
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
		linkBridges.connected = make(map[string]*Bridge)
		linkBridges.unconnected = make(map[string]*Bridge)
		if initializeBridges(stopchan) {
			for {
				select {
				default:
					if !GetAnswers(stopchan) {
						return
					}
					if !GetOffers(stopchan) {
						return
					}
					if !SendAnswers(stopchan) {
						return
					}
					if !DisconnectedLinks(stopchan) {
						return
					}
					time.Sleep(time.Second * 20)
				case <-stopchan:
					util.Info.Println("StartBridges stopped")
					return
				}
			}
		}
	} else {
		util.Info.Println("StartBridges stopped after WaitForSync")
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
