package bridge

import (
	"time"

	"github.com/duality-solutions/web-bridge/init/settings"
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
)

const (
	// OfferExpireMinutes is the number of minutes an offer found in the DHT is valid.
	OfferExpireMinutes = 360
	// MinimumOfferValueLength is the minimum offer value size to be considered
	MinimumOfferValueLength = 10
	// MinimumAnswerValueLength is the minimum offer value size to be considered
	MinimumAnswerValueLength = 10
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

func setupBridges(stopchan chan struct{}, links dynamic.ActiveLinks, accounts []dynamic.Account) bool {
	util.Info.Println("setupBridges Started")
	for _, link := range links.Links {
		select {
		default:
			var linkBridge = NewBridge(link, accounts)
			linkBridges.unconnected[linkBridge.LinkID()] = &linkBridge
		case <-stopchan:
			util.Info.Println("setupBridges stopped")
			return false
		}
	}
	return true
}

func initializeBridges(stopchan chan struct{}) bool {
	linkBridges.connected = make(map[string]*Bridge)
	linkBridges.unconnected = make(map[string]*Bridge)
	if setupBridges(stopchan, links, accounts) {
		// Get notifications received while loading and unlocking wallet
		if GetLinkNotifications(stopchan) {
			// Notify links that you are online
			if NotifyLinksOnline(stopchan) {
				util.Info.Println("Sent all online notification messages.", len(linkBridges.unconnected))
				time.Sleep(time.Second * 180) // wait 3 minutes for links to respond.
				if !GetOffers(stopchan) {
					util.Error.Println("GetOffers error")
				}
			}
		}
	}
	return true
}

// StartBridges runs a goroutine to manage network bridges
// Send online notifications to links
// Wait 3 minutes for offer reponses
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
			for {
				select {
				default:
					if !GetLinkNotifications(stopchan) {
						return
					}
					if !GetAnswers(stopchan) {
						return
					}

					/*
						if !DisconnectedLinks(stopchan) {
							return
						}
						if !StopDisconnected(stopchan) {
							return
						}
					*/
					// Update online start time every hour
				case <-stopchan:
					util.Info.Println("StartBridges stopped")
					return
				}
			}
		} else {
			util.Info.Println("StartBridges stopped after initializeBridges failed.")
		}
	} else {
		util.Info.Println("StartBridges stopped after WaitForSync")
	}
}

// ShutdownBridges stops the ManageBridges goroutine
func ShutdownBridges(stopchan chan struct{}) {
	//TODO: disconnect all active/connected WebRTC bridges
	close(stopchan)
	if !NotifyLinksOffline(stopchan) {
		util.Error.Println("ShutdownBridges NotifyLinksOffline error")
	}
}
