package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// GetOffers checks the DHT for WebRTC offers from all links
func GetOffers(stopchan chan struct{}) bool {
	//util.Info.Println("GetOffers started")
getOffers:
	for _, link := range linkBridges.connected {
		select {
		default:
			if link.State == StateNew || link.State == StateWaitForOffer {
				offers, err := dynamicd.GetLinkMessages(link.MyAccount, link.LinkAccount, "webrtc-offer")
				if err != nil {
					util.Error.Println("GetOffers error", link.LinkAccount, err)
				} else if len(*offers) > 0 {
					var offer dynamic.GetMessageReturnJSON
					for _, res := range *offers {
						if res.TimestampEpoch > offer.TimestampEpoch {
							offer = res
						}
					}
					util.Info.Println("GetOffers offer found. Size", offer.MessageSize)
					if len(offer.Message) < MinimumAnswerValueLength {
						util.Info.Println("GetOffers for", link.LinkAccount, "not found. Value too short.", len(offer.Message))
						break getOffers
					}
					var newOffer webrtc.SessionDescription
					err = util.DecodeObject(offer.Message, &newOffer)
					if err != nil {
						util.Error.Println("GetOffers DecodeObject error", link.LinkAccount, err)
						break getOffers
					}
					util.Info.Println("GetOffers offer found. Size", offer.MessageSize)
					if newOffer != link.Offer {
						link.Offer = newOffer
						pc, err := ConnectToIceServices(config)
						if err != nil {
							util.Error.Println("GetOffers ConnectToIceServices error", link.LinkAccount, err)
							break getOffers
						}
						link.PeerConnection = pc
						err = link.PeerConnection.SetRemoteDescription(link.Offer)
						if err != nil {
							util.Error.Println("GetOffers SetRemoteDescription offer error", link.LinkAccount, err)
							break getOffers
						}
						answer, err := link.PeerConnection.CreateAnswer(nil)
						if err != nil {
							util.Error.Println("GetOffers CreateAnswer error", link.LinkAccount, err)
							link.Offer = webrtc.SessionDescription{}
							break getOffers
						}
						link.Answer = answer
						encoded, err := util.EncodeObject(answer)
						if err != nil {
							util.Error.Println("GetOffers EncodeObject answer error", link.LinkAccount, err)
							break getOffers
						}
						_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-answer")
						if err != nil {
							util.Error.Println("GetOffers dynamicd.SendLinkMessage answer error", link.LinkAccount, err)
							break getOffers
						}
						link.State = StateWaitForRTC
						link.OnStateChangeEpoch = time.Now().Unix()
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
						util.Info.Println("Offer found for", link.LinkAccount, link.LinkID(), "WaitForRTC...")
						// send anwser and wait for connection or timeout.
						go WaitForRTC(link, answer)
					}
				}
			} else {
				util.Info.Println("GetOffers skipped", link.LinkAccount)
			}
		case <-stopchan:
			util.Info.Println("GetOffers stopped")
			return false
		}
	}
	return true
}

// SendOffer sends a message with the WebRTC offer embedded to the link
func SendOffer(link *Bridge) bool {
	pc, err := ConnectToIceServices(config)
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
func DisconnectedLinks(stopchan chan struct{}) bool {
	l := len(linkBridges.unconnected)
	putOffers := make(chan dynamic.DHTPutReturn, l)
	for _, link := range linkBridges.unconnected {
		if link.State == StateInit {
			util.Info.Println("DisconnectedLinks for", link.LinkParticipants(), link.LinkID())
			var linkBridge = NewLinkBridge(link.LinkAccount, link.MyAccount, accounts)
			linkBridge.SessionID = link.SessionID
			linkBridge.Get = link.Get
			pc, err := ConnectToIceServices(config)
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
		case <-stopchan:
			util.Info.Println("DisconnectedLinks stopped")
			return false
		}
	}
	return true
}
