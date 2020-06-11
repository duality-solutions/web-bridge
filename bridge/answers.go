package bridge

import (
	"fmt"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// SendAnswers uses VPG instant messages to send an answer to a WebRTC offer
func SendAnswers(stopchan chan struct{}) bool {
	fmt.Println("SendAnswers Started")
	for _, link := range linkBridges.unconnected {
		select {
		default:
			if link.State == 3 && link.PeerConnection != nil && len(link.Offer.SDP) > 10 {
				err := link.PeerConnection.SetRemoteDescription(link.Offer)
				if err != nil {
					// move to unconnected
					fmt.Printf("SendAnswers failed to connect to link %s. Error %s\n", link.LinkAccount, err)
				} else {
					answer, err := link.PeerConnection.CreateAnswer(nil)
					if err != nil {
						fmt.Println(link.LinkAccount, "SendAnswers error", err)
						// clear offer since it didn't work
						// remove from connected and add to unconnected
					} else {
						link.Answer = answer
						encoded, err := util.EncodeObject(answer)
						if err != nil {
							fmt.Println("SendAnswers EncodeObject error", link.LinkAccount, err)
						}
						_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded)
						if err != nil {
							fmt.Println("SendAnswers dynamicd.SendLinkMessage error", link.LinkAccount, err)
						}
						go WaitForRTC(link, answer)
						link.State = 4
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
					}
				}
			}
		case <-stopchan:
			fmt.Println("SendAnswers stopped")
			return false
		}
	}
	return true
}

// GetAnswers checks Dynamicd for bridge messages received
func GetAnswers(stopchan chan struct{}) bool {
	fmt.Println("GetAnswers Started")
	for _, link := range linkBridges.unconnected {
		select {
		default:
			if link.PeerConnection == nil {
				pc, err := ConnectToIceServices(config)
				if err == nil {
					link.PeerConnection = pc
				}
			}
			if link.PeerConnection != nil && link.State == 2 {
				answers, err := dynamicd.GetLinkMessages(link.MyAccount, link.LinkAccount)
				if err != nil {
					fmt.Println("GetAnswers error", link.LinkAccount, err)
				} else {
					var answer dynamic.GetMessageReturnJSON
					for _, res := range *answers {
						if res.TimestampEpoch > answer.TimestampEpoch {
							answer = res
						}
					}
					if len(answer.Message) < 10 {
						fmt.Println("GetAnswers for", link.LinkAccount, "not found")
						continue
					}
					var newAnswer webrtc.SessionDescription
					err = util.DecodeObject(answer.Message, &newAnswer)
					if err != nil {
						fmt.Println("GetAnswers DecodeObject error", link.LinkAccount, err)
						continue
					}
					if newAnswer != link.Answer {
						link.Answer = newAnswer
						go EstablishRTC(link)
						link.State = 5
						delete(linkBridges.unconnected, link.LinkID())
						linkBridges.connected[link.LinkID()] = link
					}
				}
			}
		case <-stopchan:
			fmt.Println("GetAnswers stopped")
			return false
		}
	}
	fmt.Println("GetAnswers complete")
	return true
}
