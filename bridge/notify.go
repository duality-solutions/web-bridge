package bridge

import (
	"time"

	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
)

var startEpoch int64
var endEpoch int64
var mapGetNotifications map[string]dynamic.GetMessageReturnJSON = make(map[string]dynamic.GetMessageReturnJSON)

// OnlineNotification stores start and end WebBridge session time
type OnlineNotification struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// NotifyLinksOnline sends a VGP message to all links with online status
func NotifyLinksOnline(stopchan *chan struct{}) bool {
	if bridgeControler == nil {
		return true
	}
	util.Info.Println("NotifyLinksOnline Started")
	endEpoch = 0
	startEpoch = time.Now().Unix()
	bridges := bridgeControler.Unconnected()
	for _, link := range bridges {
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
		case <-*stopchan:
			util.Info.Println("NotifyLinksOnline stopped")
			return false
		}
	}
	return true
}

// NotifyLinksOffline sends a VGP message to all links with offline status
func NotifyLinksOffline() bool {
	if bridgeControler == nil {
		return true
	}
	util.Info.Println("NotifyLinksOffline Started")
	l := bridgeControler.Count()
	bridges := bridgeControler.Unconnected()
	sendOfflineChan := make(chan dynamic.MessageReturnJSON, l)
	endEpoch = time.Now().Unix()
	for _, link := range bridges {
		notification := OnlineNotification{
			StartTime: startEpoch,
			EndTime:   endEpoch,
		}
		encoded, err := util.EncodeObject(notification)
		if err != nil {
			util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
			break
		}
		dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOfflineChan)
	}
	bridges = bridgeControler.Connected()
	for _, link := range bridges {
		notification := OnlineNotification{
			StartTime: startEpoch,
			EndTime:   endEpoch,
		}
		encoded, err := util.EncodeObject(notification)
		if err != nil {
			util.Error.Println("NotifyLinksOffline EncodeObject error", link.LinkAccount, err)
			break
		}
		dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOfflineChan)
	}
	for i := uint16(0); i < l; i++ {
		select {
		default:
			notification := <-sendOfflineChan
			util.Info.Println("NotifyLinksOffline", notification.SubjectID)
		case <-time.After(time.Second * 30):
			util.Error.Println("NotifyLinksOffline timeout after 30 seconds")
			return false
		}
	}
	return true
}

// GetLinkNotifications gets online notification messages for all unconnected links
func GetLinkNotifications(stopchan *chan struct{}) bool {
	//util.Info.Println("GetLinkNotifications Started")
	bridges := bridgeControler.Unconnected()
	for _, link := range bridges {
		select {
		default:
			notifications, err := dynamicd.GetNotificationMessages(link.MyAccount, link.LinkAccount)
			if err != nil {
				util.Error.Println("GetLinkNotifications dynamicd.GetNotificationMessages error", link.LinkAccount, err)
			}
			for _, notification := range *notifications {
				currentNotification := mapGetNotifications[link.LinkID()]
				if notification.TimestampEpoch > currentNotification.TimestampEpoch {
					mapGetNotifications[link.LinkID()] = notification
					var online OnlineNotification
					err := util.DecodeObject(notification.Message, &online)
					if err != nil {
						util.Error.Println("GetLinkNotifications EncodeObject error", link.LinkAccount, err)
						break
					}
					if online.EndTime == 0 /* && (time.Now().Unix() - online.StartTime) < 36000 */ {
						// send offer
						util.Info.Println("GetLinkNotifications online message from", link.LinkAccount, "secs:", (time.Now().Unix() - online.StartTime))
						if SendOffer(link) {
							link.SetState(StateWaitForAnswer)
							bridgeControler.MoveUnconnectedToConnected(link)
						} else {
							break
						}
					} else {
						// delete notification from map?
						util.Info.Println("GetLinkNotifications offline message from", link.LinkAccount, "secs:", (time.Now().Unix() - online.EndTime))
					}
				}
			}
		case <-*stopchan:
			util.Info.Println("GetLinkNotifications stopped")
			return false
		}
	}
	return true
}
