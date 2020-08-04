package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// SendAnswers uses VPG instant messages to send an answer to a WebRTC offer
func SendAnswers(stopchan chan struct{}) bool {
	util.Info.Println("SendAnswers Started")
	for _, link := range linkBridges.unconnected {
		select {
		default:
			if link.State == StateSendAnswer && link.PeerConnection != nil && len(link.Offer.SDP) > MinimumOfferValueLength {
				err := link.PeerConnection.SetRemoteDescription(link.Offer)
				if err != nil {
					// move to unconnected
					util.Error.Printf("SendAnswers failed to connect to link %s. Error %s\n", link.LinkAccount, err)
				} else {
					answer, err := link.PeerConnection.CreateAnswer(nil)
					if err != nil {
						util.Error.Println(link.LinkAccount, "SendAnswers error", err)
						// clear offer since it didn't work
						// remove from connected and add to unconnected
					} else {
						link.Answer = answer
						encoded, err := util.EncodeObject(answer)
						if err != nil {
							util.Error.Println("SendAnswers EncodeObject error", link.LinkAccount, err)
						}
						_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded, "webrtc-answer")
						if err != nil {
							util.Error.Println("SendAnswers dynamicd.SendLinkMessage error", link.LinkAccount, err)
						}
						link.State = StateWaitForRTC
						link.OnStateChangeEpoch = time.Now().Unix()
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
						go WaitForRTC(link, answer)
					}
				}
			}
		case <-stopchan:
			util.Info.Println("SendAnswers stopped")
			return false
		}
	}
	return true
}

// GetAnswers checks Dynamicd for bridge messages received
func GetAnswers(stopchan chan struct{}) bool {
	//util.Info.Println("GetAnswers Started")
getAnswers:
	for _, link := range linkBridges.connected {
		select {
		default:
			if link.State != StateWaitForAnswer {
				continue getAnswers
			}
			if link.PeerConnection == nil {
				pc, err := ConnectToIceServicesDetached(config)
				if err == nil {
					link.PeerConnection = pc
				}
			}
			if link.PeerConnection != nil && link.State == StateWaitForAnswer {
				answers, err := dynamicd.GetLinkMessages(link.MyAccount, link.LinkAccount, "webrtc-answer")
				if err != nil {
					util.Error.Println("GetAnswers error", link.LinkAccount, err)
				} else if len(*answers) > 0 {
					var answer dynamic.GetMessageReturnJSON
					for _, res := range *answers {
						if res.TimestampEpoch > answer.TimestampEpoch {
							answer = res
						}
					}
					util.Info.Println("GetAnswers offer found. Size", answer.MessageSize)
					if len(answer.Message) < MinimumAnswerValueLength {
						util.Info.Println("GetAnswers for", link.LinkAccount, "not found. Value too short.", len(answer.Message))
						break getAnswers
					}
					var newAnswer webrtc.SessionDescription

					err = util.DecodeObject(answer.Message, &newAnswer)
					if err != nil {
						util.Error.Println("GetAnswers DecodeObject error", link.LinkAccount, err)
						break getAnswers
					}
					if newAnswer != link.Answer {
						link.Answer = newAnswer
						link.State = StateEstablishRTC
						link.OnStateChangeEpoch = time.Now().Unix()
						util.Info.Println("Answer found for", link.LinkAccount, link.LinkID(), "EstablishRTC...")
						// send anwser and wait for connection or timeout.
						go EstablishRTC(link)
					}
				}
			}
		case <-stopchan:
			util.Info.Println("GetAnswers stopped")
			return false
		}
	}
	//util.Info.Println("GetAnswers complete")
	return true
}
