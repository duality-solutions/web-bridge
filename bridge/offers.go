package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

var mapOffers map[string]int = make(map[string]int)

// GetOffers checks the DHT for WebRTC offers from all links
func GetOffers(stopchan *chan struct{}) bool {
	bridges := bridgeControler.Unconnected()
	l := len(bridges)
	getOffersChan := make(chan dynamic.GetVGPMessageReturn, l)
	for _, link := range bridges {
		select {
		default:
			dynamicd.GetLinkMessagesAsync(link.LinkID(), link.MyAccount, link.LinkAccount, "webrtc-offer", getOffersChan)
		case <-*stopchan:
			util.Info.Println("GetOffers stopped")
			return false
		}
	}
getOffersLoop:
	for i := 0; i < l; i++ {
		select {
		default:
			offers := <-getOffersChan
			link := bridgeControler.GetUnconnected(offers.LinkID)
			var offer dynamic.GetMessageReturnJSON
			for _, res := range offers.Messages {
				if res.TimestampEpoch > offer.TimestampEpoch {
					offer = res
				}
			}
			if offer.MessageSize > 0 && !(mapOffers[offer.MessageID] > 0) {
				util.Info.Println("GetOffers new offer found. Size", offer.MessageSize)
				mapOffers[offer.MessageID] = offer.TimestampEpoch
				if len(offer.Message) < MinimumAnswerValueLength {
					util.Info.Println("GetOffers for", link.LinkAccount, "not found. Value too short.", len(offer.Message))
					break getOffersLoop
				}
				var newOffer webrtc.SessionDescription
				err := util.DecodeObject(offer.Message, &newOffer)
				if err != nil {
					util.Error.Println("GetOffers DecodeObject error", link.LinkAccount, err)
					break getOffersLoop
				}
				util.Info.Println("GetOffers offer found. Size", offer.MessageSize)
				if newOffer != link.Offer() {
					link.SetOffer(newOffer)
					pc, err := ConnectToIceServicesDetached(config)
					if err != nil {
						util.Error.Println("GetOffers ConnectToIceServices error", link.LinkAccount, err)
						break getOffersLoop
					}
					link.SetPeerConnection(pc)
					err = link.PeerConnection().SetRemoteDescription(link.Offer())
					if err != nil {
						util.Error.Println("GetOffers SetRemoteDescription offer error", link.LinkAccount, err)
						break getOffersLoop
					}
					answer, err := link.PeerConnection().CreateAnswer(nil)
					if err != nil {
						util.Error.Println("GetOffers CreateAnswer error", link.LinkAccount, err)
						link.SetOffer(webrtc.SessionDescription{})
						break getOffersLoop
					}
					encoded, err := util.EncodeObject(answer)
					if err != nil {
						util.Error.Println("GetOffers EncodeObject answer error", link.LinkAccount, err)
						break getOffersLoop
					}
					_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-answer")
					if err != nil {
						util.Error.Println("GetOffers dynamicd.SendLinkMessage answer error", link.LinkAccount, err)
						break getOffersLoop
					}
					link.SetAnswerStateEpoch(answer, StateWaitForRTC, time.Now().Unix())
					bridgeControler.MoveUnconnectedToConnected(link)
					util.Info.Println("Offer found for", link.LinkAccount, link.LinkID(), "WaitForRTC...")
					// send anwser and wait for connection or timeout.
					go WaitForRTC(link)
				}
			}
		case <-*stopchan:
			util.Info.Println("GetOffers stopped")
			return false
		}
	}
	return true
}

// SendOffer sends a message with the WebRTC offer embedded to the link
func SendOffer(link *Bridge) bool {
	pc, err := ConnectToIceServicesDetached(config)
	if err != nil {
		util.Error.Println("SendOffer error connecting tot ICE services", err)
		return false
	}
	link.SetPeerConnection(pc)
	dataChannel, err := link.PeerConnection().CreateDataChannel(link.LinkParticipants(), nil)
	if err != nil {
		util.Error.Println("SendOffer error creating dataChannel for", link.LinkAccount, link.LinkID())
		return false
	}
	link.SetDataChannel(dataChannel)
	offer, err := link.PeerConnection().CreateOffer(nil)
	if err != nil {
		util.Error.Println("SendOffer error CreateOffer", err)
		return false
	}
	link.SetOffer(offer)
	encoded, err := util.EncodeObject(link.Offer())
	if err != nil {
		util.Error.Println("SendOffer error EncodeObject", err)
		return false
	}
	_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-offer")
	if err != nil {
		util.Error.Println("GetOffers dynamicd.SendLinkMessage answer error", link.LinkAccount, err)
		return false
	}
	link.SetState(StateWaitForAnswer)
	link.SetOnStateChangeEpoch(time.Now().Unix())
	return true
}

// DisconnectedLinks reinitializes the WebRTC link bridge struct
func DisconnectedLinks(stopchan *chan struct{}) bool {
	bridges := bridgeControler.Unconnected()
	endEpoch = 0
	startEpoch = time.Now().Unix()
	for _, link := range bridges {
		if link.State() == StateInit && link.RTCState() == "closed" {
			notification := OnlineNotification{
				StartTime: startEpoch,
				EndTime:   endEpoch,
			}
			encoded, err := util.EncodeObject(notification)
			if err != nil {
				util.Error.Println("DisconnectedLinks EncodeObject error", link.LinkAccount, err)
				break
			}
			sendOnlineChan := make(chan dynamic.MessageReturnJSON, 1)
			dynamicd.SendNotificationMessageAsync(link.MyAccount, link.LinkAccount, encoded, sendOnlineChan)
			util.Info.Println("DisconnectedLinks sent online notification to", link.LinkAccount, encoded)
			link.SetRTCState("")
		}
	}
	return true
}
