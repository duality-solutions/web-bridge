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
	l := len(linkBridges.unconnected)
	getOffersChan := make(chan dynamic.GetVGPMessageReturn, l)
	for _, link := range linkBridges.unconnected {
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
			link := linkBridges.unconnected[offers.LinkID]
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
				if newOffer != link.Offer {
					link.Offer = newOffer
					pc, err := ConnectToIceServicesDetached(config)
					if err != nil {
						util.Error.Println("GetOffers ConnectToIceServices error", link.LinkAccount, err)
						break getOffersLoop
					}
					link.PeerConnection = pc
					err = link.PeerConnection.SetRemoteDescription(link.Offer)
					if err != nil {
						util.Error.Println("GetOffers SetRemoteDescription offer error", link.LinkAccount, err)
						break getOffersLoop
					}
					answer, err := link.PeerConnection.CreateAnswer(nil)
					if err != nil {
						util.Error.Println("GetOffers CreateAnswer error", link.LinkAccount, err)
						link.Offer = webrtc.SessionDescription{}
						break getOffersLoop
					}
					link.Answer = answer
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
					link.State = StateWaitForRTC
					link.OnStateChangeEpoch = time.Now().Unix()
					delete(linkBridges.unconnected, link.LinkID())
					linkBridges.connected[link.LinkID()] = link
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
	link.PeerConnection = pc
	dataChannel, err := link.PeerConnection.CreateDataChannel(link.LinkParticipants(), nil)
	if err != nil {
		util.Error.Println("SendOffer error creating dataChannel for", link.LinkAccount, link.LinkID())
		return false
	}
	link.DataChannel = dataChannel
	link.Offer, err = link.PeerConnection.CreateOffer(nil)
	if err != nil {
		util.Error.Println("SendOffer error CreateOffer", err)
		return false
	}
	encoded, err := util.EncodeObject(link.Offer)
	if err != nil {
		util.Error.Println("SendOffer error EncodeObject", err)
		return false
	}
	_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-offer")
	if err != nil {
		util.Error.Println("GetOffers dynamicd.SendLinkMessage answer error", link.LinkAccount, err)
		return false
	}
	link.State = StateWaitForAnswer
	link.OnStateChangeEpoch = time.Now().Unix()
	return true
}

// DisconnectedLinks reinitializes the WebRTC link bridge struct
func DisconnectedLinks(stopchan *chan struct{}) bool {
	l := len(linkBridges.unconnected)
	putOffers := make(chan dynamic.DHTPutReturn, l)
	for _, link := range linkBridges.unconnected {
		if link.State == StateInit {
			util.Info.Println("DisconnectedLinks for", link.LinkParticipants(), link.LinkID())
			var linkBridge = NewLinkBridge(link.LinkAccount, link.MyAccount, accounts)
			linkBridge.SessionID = link.SessionID
			linkBridge.Get = link.Get
			pc, err := ConnectToIceServicesDetached(config)
			if err != nil {
				util.Error.Println("DisconnectedLinks error connecting tot ICE services", err)
				continue
			} else {
				linkBridge.PeerConnection = pc
				dataChannel, err := linkBridge.PeerConnection.CreateDataChannel(link.LinkParticipants(), nil)
				if err != nil {
					util.Error.Println("DisconnectedLinks error creating dataChannel for", link.LinkAccount, link.LinkID())
					continue
				}
				linkBridge.DataChannel = dataChannel
			}
			linkBridge.Offer, _ = linkBridge.PeerConnection.CreateOffer(nil)
			linkBridge.Answer = link.Answer
			encoded, err := util.EncodeObject(linkBridge.Offer)
			if err != nil {
				util.Info.Println("DisconnectedLinks error EncodeObject", err)
			}
			dynamicd.PutLinkRecord(linkBridge.MyAccount, linkBridge.LinkAccount, encoded, putOffers)
			linkBridge.State = StateWaitForAnswer
			link.OnStateChangeEpoch = time.Now().Unix()
			linkBridges.unconnected[linkBridge.LinkID()] = &linkBridge
		} else {
			l--
		}
	}
	for i := 0; i < l; i++ {
		select {
		default:
			offer := <-putOffers
			linkBridge := NewLinkBridge(offer.Sender, offer.Receiver, accounts)
			link := linkBridges.unconnected[linkBridge.LinkID()]
			link.Put = offer.DHTPutJSON
			util.Info.Println("DisconnectedLinks Offer saved", offer)
		case <-*stopchan:
			util.Info.Println("DisconnectedLinks stopped")
			return false
		}
	}
	return true
}
