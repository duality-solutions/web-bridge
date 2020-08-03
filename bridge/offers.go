package bridge

import (
	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// GetOffers checks the DHT for WebRTC offers from all links
func GetOffers(stopchan chan struct{}) bool {
	util.Info.Println("GetOffers started")
	for _, linkBridge := range linkBridges.unconnected {
		select {
		default:
			var link = linkBridges.unconnected[linkBridge.LinkID()]
			if link.State == StateNew || link.State == StateWaitForOffer {
				offers, err := dynamicd.GetLinkMessages(link.MyAccount, link.LinkAccount, "webrtc-offer")
				if err != nil {
					util.Error.Println("GetOffers error", link.LinkAccount, err)
				} else {
					var offer dynamic.GetMessageReturnJSON
					for _, res := range *offers {
						if res.TimestampEpoch > offer.TimestampEpoch {
							offer = res
						}
					}
					if len(offer.Message) < MinimumAnswerValueLength {
						util.Info.Println("GetOffers for", link.LinkAccount, "not found")
						continue
					}
					var newOffer webrtc.SessionDescription
					err = util.DecodeObject(offer.Message, &newOffer)
					if err != nil {
						util.Error.Println("GetOffers DecodeObject error", link.LinkAccount, err)
						continue
					}
					if newOffer != link.Offer {
						link.Offer = newOffer
						err := link.PeerConnection.SetRemoteDescription(link.Offer)
						if err != nil {
							util.Error.Println("GetOffers SetRemoteDescription offer error", link.LinkAccount, err)
							continue
						}
						answer, err := link.PeerConnection.CreateAnswer(nil)
						if err != nil {
							util.Error.Println(link.LinkAccount, "GetOffers error", err)
							link.Offer = webrtc.SessionDescription{}
							continue
						}
						link.Answer = answer
						encoded, err := util.EncodeObject(answer)
						if err != nil {
							util.Error.Println("GetOffers EncodeObject error", link.LinkAccount, err)
						}
						_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-answer")
						if err != nil {
							util.Error.Println("GetOffers dynamicd.SendLinkMessage error", link.LinkAccount, err)
						}
						link.State = StateWaitForRTC
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
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
	} else {
		link.PeerConnection = pc
		dataChannel, err := link.PeerConnection.CreateDataChannel(link.LinkParticipants(), nil)
		if err != nil {
			util.Error.Println("SendOffer error creating dataChannel for", link.LinkAccount, link.LinkID())
			return false
		}
		link.DataChannel = dataChannel
	}
	link.Offer, err = link.PeerConnection.CreateOffer(nil)
	if err != nil {
		util.Error.Println("SendOffer error CreateOffer", err)
	}
	encoded, err := util.EncodeObject(link.Offer)
	if err != nil {
		util.Error.Println("SendOffer error EncodeObject", err)
	}
	dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-offer")
	link.State = StateWaitForAnswer
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
