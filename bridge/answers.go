package bridge

import (
	"fmt"
	"strings"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// SendAnswers uses VPG instant messages to send an answer to a WebRTC offer
func SendAnswers(bridges *[]Bridge) {
	for _, brd := range *bridges {
		if brd.PeerConnection != nil && len(brd.Offer) > 10 {
			brd.Offer = strings.ReplaceAll(brd.Offer, `""`, "") // remove double quotes in offer
			sd := webrtc.SessionDescription{Type: 1, SDP: brd.Offer}
			//fmt.Println("sendAnswers", brd.LinkAccount, "SessionDescription", sd)
			err := brd.PeerConnection.SetRemoteDescription(sd)
			if err != nil {
				// move to unconnected
				fmt.Printf("SendAnswers failed to connect to link %s. Error %s\n", brd.LinkAccount, err)
			} else {
				answer, err := brd.PeerConnection.CreateAnswer(nil)
				if err != nil {
					fmt.Println(brd.LinkAccount, "SendAnswers error", err)
					// clear offer since it didn't work
					// remove from connected and add to unconnected
				} else {
					//fmt.Println("SendLinkMessage", brd.LinkAccount, answer.SDP)
					_, err := dynamicd.SendLinkMessage(brd.MyAccount, brd.LinkAccount, answer.SDP)
					if err != nil {
						fmt.Println("SendLinkMessage error", brd.LinkAccount, err)
					}
				}
			}
		} else {
			fmt.Println("Error nil PeerConnection", brd.LinkAccount)
		}
	}
}

// GetAnswers checks Dynamicd for bridge messages received
func GetAnswers(bridges *[]Bridge) {
	for _, brd := range *bridges {
		if brd.PeerConnection == nil {
			pc, err := ConnectToIceServices(config)
			if err == nil {
				brd.PeerConnection = pc
			}
		}
		config := brd.PeerConnection.GetConfiguration()
		if brd.PeerConnection != nil {
			answers, err := dynamicd.GetLinkMessages(brd.MyAccount, brd.LinkAccount)
			if err != nil {
				fmt.Println("GetAnswers error", brd.LinkAccount, err)
			} else {
				//fmt.Println("GetAnswers", answers)
				var answer dynamic.GetMessageReturnJSON
				for _, res := range *answers {
					if res.TimestampEpoch > answer.TimestampEpoch {
						answer = res
					}
				}
				brd.Answer = strings.ReplaceAll(answer.Message, `""`, "") // remove double quotes in answer
				sd := webrtc.SessionDescription{Type: 2, SDP: brd.Answer}
				err := brd.PeerConnection.SetRemoteDescription(sd)
				if err != nil {
					fmt.Println("GetAnswers SetRemoteDescription error ", err)
				} else {
					dc, err := brd.PeerConnection.CreateDataChannel(brd.LinkAccount, nil)
					if err != nil {
						fmt.Println("GetAnswers CreateDataChannel error", err)
					}
					fmt.Println("GetAnswers Data Channel", dc)
				}
			}
		} else {
			fmt.Println("Error nil PeerConnection", brd.LinkAccount)
		}
	}
}
