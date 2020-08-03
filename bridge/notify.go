package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
)

var startEpoch int64
var endEpoch int64

// OnlineNotification stores start and end WebBridge session time
type OnlineNotification struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// NotifyLinksOnline sends a VGP message to all links with online status
func NotifyLinksOnline(stopchan chan struct{}, links dynamic.ActiveLinks, accounts []dynamic.Account) bool {
	util.Info.Println("NotifyLinksOnline Started")
	endEpoch = 0
	startEpoch = time.Now().Unix()
	for _, link := range links.Links {
		select {
		default:
			var linkBridge = NewBridge(link, accounts)
			linkBridges.unconnected[linkBridge.LinkID()] = &linkBridge
			notification := OnlineNotification{
				StartTime: startEpoch,
				EndTime:   endEpoch,
			}
			encoded, err := util.EncodeObject(notification)
			if err != nil {
				util.Error.Println("NotifyLinksOnline EncodeObject error", linkBridge.LinkAccount, err)
				continue
			}
			_, err = dynamicd.SendNotificationMessage(linkBridge.MyAccount, linkBridge.LinkAccount, encoded)
			if err != nil {
				util.Error.Println("NotifyLinksOnline dynamicd.SendNotificationMessage error", linkBridge.LinkAccount, err)
			}
		case <-stopchan:
			util.Info.Println("NotifyLinksOnline stopped")
			return false
		}
	}
	return true
}

// NotifyLinksOffline sends a VGP message to all links with offline status
func NotifyLinksOffline(stopchan chan struct{}) bool {
	util.Info.Println("NotifyLinksOffline Started")
	endEpoch = time.Now().Unix()
	for _, link := range linkBridges.unconnected {
		select {
		default:
			notification := OnlineNotification{
				StartTime: startEpoch,
				EndTime:   endEpoch,
			}
			encoded, err := util.EncodeObject(notification)
			if err != nil {
				util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
				continue
			}
			_, err = dynamicd.SendNotificationMessage(link.MyAccount, link.LinkAccount, encoded)
			if err != nil {
				util.Error.Println("NotifyLinksOffline dynamicd.SendNotificationMessage error", link.LinkAccount, err)
			}
		case <-stopchan:
			util.Info.Println("NotifyLinksOffline stopped")
			return false
		}
	}
	for _, link := range linkBridges.connected {
		select {
		default:
			var link = linkBridges.unconnected[link.LinkID()]
			notification := OnlineNotification{
				StartTime: startEpoch,
				EndTime:   endEpoch,
			}
			encoded, err := util.EncodeObject(notification)
			if err != nil {
				util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
				continue
			}
			_, err = dynamicd.SendNotificationMessage(link.MyAccount, link.LinkAccount, encoded)
			if err != nil {
				util.Error.Println("NotifyLinksOffline dynamicd.SendNotificationMessage error", link.LinkAccount, err)
			}
		case <-stopchan:
			util.Info.Println("NotifyLinksOffline stopped")
			return false
		}
	}
	return true
}

// GetLinkNotifications get
func GetLinkNotifications(stopchan chan struct{}) bool {
	util.Info.Println("GetLinktNotifications Started")
	for _, link := range linkBridges.unconnected {
		select {
		default:
			notifications, err := dynamicd.GetNotificationMessages(link.MyAccount, link.LinkAccount)
			if err != nil {
				util.Error.Println("GetLinktNotifications dynamicd.GetNotificationMessages error", link.LinkAccount, err)
			}
			for _, notification := range *notifications {
				// send offer
				var online OnlineNotification
				err := util.DecodeObject(notification.Message, online)
				if err != nil {
					util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
					continue
				}
				if online.EndTime == 0 /* && (time.Now().Unix() - online.StartTime) < 36000 */ {
					SendOffer(link)
					link.State = StateWaitForAnswer
					delete(linkBridges.unconnected, link.LinkID())
					linkBridges.connected[link.LinkID()] = link
				}
			}
		case <-stopchan:
			util.Info.Println("NotifyLinksOffline stopped")
			return false
		}
	}
	return true
}
