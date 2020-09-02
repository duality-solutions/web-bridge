package bridge

import (
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

var mapGetAnswers map[string]dynamic.GetMessageReturnJSON = make(map[string]dynamic.GetMessageReturnJSON)

// SendAnswers uses VPG instant messages to send an answer to a WebRTC offer
func SendAnswers(stopchan *chan struct{}) bool {
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
						go WaitForRTC(link)
					}
				}
			}
		case <-*stopchan:
			util.Info.Println("SendAnswers stopped")
			return false
		}
	}
	return true
}

// GetAnswers checks Dynamicd for bridge messages received
func GetAnswers(stopchan *chan struct{}) bool {
	l := len(linkBridges.connected)
	getAnswersChan := make(chan dynamic.GetVGPMessageReturn, l)
	for _, link := range linkBridges.connected {
		select {
		default:
			if link.State == StateWaitForAnswer {
				dynamicd.GetLinkMessagesAsync(link.LinkID(), link.MyAccount, link.LinkAccount, "webrtc-answer", getAnswersChan)
			} else {
				l--
			}
		case <-*stopchan:
			util.Info.Println("GetAnswers stopped")
			return false
		}
	}
getAnswersLoop:
	for i := 0; i < l; i++ {
		select {
		default:
			answers := <-getAnswersChan
			link := linkBridges.connected[answers.LinkID]
			if len(answers.Messages) > 0 && link.State == StateWaitForAnswer {
				currentAnswer := mapGetAnswers[link.LinkID()]
				if link.PeerConnection == nil {
					pc, err := ConnectToIceServicesDetached(config)
					if err == nil {
						link.PeerConnection = pc
					}
				}
				var answer dynamic.GetMessageReturnJSON
				for _, res := range answers.Messages {
					if res.TimestampEpoch > answer.TimestampEpoch {
						answer = res
					}
				}
				if currentAnswer != answer && answer.MessageSize > 0 {
					util.Info.Println("GetAnswers answers found. Size", answer.MessageSize)
					if len(answer.Message) < MinimumAnswerValueLength {
						util.Info.Println("GetAnswers for", link.LinkAccount, "not found. Value too short.", len(answer.Message))
						break getAnswersLoop
					}
					var newAnswer webrtc.SessionDescription
					err := util.DecodeObject(answer.Message, &newAnswer)
					if err != nil {
						util.Error.Println("GetAnswers DecodeObject error", link.LinkAccount, err)
						break getAnswersLoop
					}
					if newAnswer != link.Answer {
						mapGetAnswers[link.LinkID()] = answer
						link.Answer = newAnswer
						link.State = StateEstablishRTC
						link.OnStateChangeEpoch = time.Now().Unix()
						util.Info.Println("Answer found for", link.LinkAccount, link.LinkID(), "EstablishRTC...")
						// send anwser and wait for connection or timeout.
						go EstablishRTC(link)
					}
				}
			}
		case <-*stopchan:
			util.Info.Println("GetAnswers stopped")
			return false
		}
	}
	return true
}
