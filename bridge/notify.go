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
func NotifyLinksOnline(stopchan chan struct{}) bool {
	util.Info.Println("NotifyLinksOnline Started")
	endEpoch = 0
	startEpoch = time.Now().Unix()
	for _, link := range linkBridges.unconnected {
		select {
		default:
			notification := OnlineNotification{
				StartTime: startEpoch,
				EndTime:   endEpoch,
			}
			encoded, err := util.EncodeObject(notification)
			if err != nil {
				util.Error.Println("NotifyLinksOnline EncodeObject error", link.LinkAccount, err)
				break
			}
			sendOnlineChan := make(chan dynamic.MessageReturnJSON, 1)
			dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOnlineChan)
			util.Info.Println("NotifyLinksOnline sent", link.LinkAccount, encoded)
		case <-stopchan:
			util.Info.Println("NotifyLinksOnline stopped")
			return false
		}
	}
	return true
}

// NotifyLinksOffline sends a VGP message to all links with offline status
func NotifyLinksOffline() bool {
	util.Info.Println("NotifyLinksOffline Started")
	endEpoch = time.Now().Unix()
	for _, link := range linkBridges.unconnected {
		notification := OnlineNotification{
			StartTime: startEpoch,
			EndTime:   endEpoch,
		}
		encoded, err := util.EncodeObject(notification)
		if err != nil {
			util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
			break
		}
		sendOfflineChan := make(chan dynamic.MessageReturnJSON, 1)
		dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOfflineChan)
	}
	for _, link := range linkBridges.connected {
		notification := OnlineNotification{
			StartTime: startEpoch,
			EndTime:   endEpoch,
		}
		encoded, err := util.EncodeObject(notification)
		if err != nil {
			util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
			break
		}
		sendOfflineChan := make(chan dynamic.MessageReturnJSON, 1)
		dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOfflineChan)
	}
	return true
}

// GetLinkNotifications gets online notification messages for all unconnected links
func GetLinkNotifications(stopchan chan struct{}) bool {
	//util.Info.Println("GetLinkNotifications Started")
	for _, link := range linkBridges.unconnected {
		select {
		default:
			notifications, err := dynamicd.GetNotificationMessages(link.MyAccount, link.LinkAccount)
			if err != nil {
				util.Error.Println("GetLinkNotifications dynamicd.GetNotificationMessages error", link.LinkAccount, err)
			}
			for _, notification := range *notifications {
				var online OnlineNotification
				err := util.DecodeObject(notification.Message, &online)
				if err != nil {
					util.Error.Println("GetLinkNotifications EncodeObject error", link.LinkAccount, err)
					break
				}
				util.Info.Println("GetLinkNotifications message from", link.LinkAccount, "secs:", (time.Now().Unix() - online.StartTime), online.EndTime)
				if online.EndTime == 0 /* && (time.Now().Unix() - online.StartTime) < 36000 */ {
					// send offer
					if SendOffer(link) {
						link.State = StateWaitForAnswer
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
					} else {
						break
					}
				}
			}
		case <-stopchan:
			util.Info.Println("GetLinkNotifications stopped")
			return false
		}
	}
	return true
}
