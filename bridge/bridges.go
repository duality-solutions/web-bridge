package bridge

import (
	"errors"
	"time"

	"github.com/duality-solutions/web-bridge/init/settings"
	"github.com/duality-solutions/web-bridge/internal/util"
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

var bridgeControler *Controller
var dynamicd dynamic.Dynamicd
var config settings.Configuration
var accounts []dynamic.Account
var links dynamic.ActiveLinks

// Bridges hold all link WebRTC bridges
type Bridges struct {
	connected   map[string]*Bridge
	unconnected map[string]*Bridge
}

// Controler returns the internal bridge controller if it exists
func Controler() (*Controller, error) {
	if bridgeControler == nil {
		return nil, errors.New("The bridge controller is not intialized yet")
	}
	return bridgeControler, nil
}

func setupBridges(stopchan *chan struct{}, links dynamic.ActiveLinks, accounts []dynamic.Account) bool {
	util.Info.Println("setupBridges Started")
	// init bridge controller
	bridgeControler = NewController()
	sessionID := uint16(0)
	for _, link := range links.Links {
		select {
		default:
			var linkBridge = NewBridge(sessionID, link, accounts)
			bridgeControler.PutUnconnected(&linkBridge)
		case <-*stopchan:
			util.Info.Println("setupBridges stopped")
			return false
		}
		sessionID = sessionID + 2
	}
	return true
}

func initializeBridges(stopchan *chan struct{}) bool {
	if setupBridges(stopchan, links, accounts) {
		// Get notifications received while loading and unlocking wallet
		if GetLinkNotifications(stopchan) {
			// Notify links that you are online
			if NotifyLinksOnline(stopchan) {
				util.Info.Println("Sent all online notification messages.", bridgeControler.Count())
				time.Sleep(time.Second * 60) // wait 1 minute for links to respond.
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
// Wait 1 minute for offer reponses
// send bridge result to upstream channel
// put offers in the DHT for unconnected links
// check for answers loop
// if new answer found, create a WebRTC bridge and send bridge result to upstream channel
// on shutdown, clear all DHT offers
func StartBridges(stopchan *chan struct{}, c settings.Configuration, d dynamic.Dynamicd, a []dynamic.Account, l dynamic.ActiveLinks) {
	dynamicd = d
	config = c
	accounts = a
	links = l
	if dynamicd.WaitForSync(stopchan, 10, 10) {
		if initializeBridges(stopchan) {
			for {
				select {
				case <-time.After(3 * time.Second):
					if !GetLinkNotifications(stopchan) {
						return
					}
					if !GetAnswers(stopchan) {
						return
					}
					if !GetOffers(stopchan) {
						return
					}
					if !DisconnectedLinks(stopchan) {
						return
					}
					/*
						if !StopDisconnected(stopchan) {
							return
						}
					*/
					// Update online start time every hour
				case <-*stopchan:
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
func ShutdownBridges(stopchan *chan struct{}) {
	//TODO: disconnect all active/connected WebRTC bridges
	close(*stopchan)
	if !NotifyLinksOffline() {
		util.Error.Println("ShutdownBridges NotifyLinksOffline error")
	}
}
