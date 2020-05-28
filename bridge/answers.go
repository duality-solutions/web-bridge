package bridge

import (
	"fmt"
	"strings"

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
				fmt.Printf("sendAnswers failed to connect to link %s. Error %s\n", brd.LinkAccount, err)
			} else {
				answer, err := brd.PeerConnection.CreateAnswer(nil)
				if err != nil {
					fmt.Println(brd.LinkAccount, "CreateAnswer error", err)
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
