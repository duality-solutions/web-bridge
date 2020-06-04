package bridge

import (
	"fmt"
	"strings"
	"time"

	util "github.com/duality-solutions/web-bridge/internal/utilities"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/pion/webrtc/v2"
)

// SendAnswers uses VPG instant messages to send an answer to a WebRTC offer
func SendAnswers(stopchan chan struct{}) bool {
	fmt.Println("SendAnswers Started")
	for _, link := range linkBridges.connected {
		select {
		default:
			if link.PeerConnection != nil && len(link.Offer) > 10 {
				link.Offer = strings.ReplaceAll(link.Offer, `""`, "") // remove double quotes in offer
				sd := webrtc.SessionDescription{Type: 1, SDP: link.Offer}
				//fmt.Println("sendAnswers", link.LinkAccount, "SessionDescription", sd)
				err := link.PeerConnection.SetRemoteDescription(sd)
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
						sd = webrtc.SessionDescription{Type: 2, SDP: link.Answer}
						err = link.PeerConnection.SetLocalDescription(sd)
						if err != nil {
							fmt.Println("SendAnswers SetLocalDescription error ", err)
						} else {
							dc, err := link.PeerConnection.CreateDataChannel(link.LinkAccount, nil)
							if err != nil {
								fmt.Println("GetAnswers CreateDataChannel error", err)
							}
							fmt.Println("GetAnswers Data Channel Negotiated", dc.Negotiated())
						}
						//fmt.Println("SendLinkMessage", link.LinkAccount, answer.SDP)
						encoded, err := util.EncodeString(answer.SDP)
						if err != nil {
							fmt.Println("SendAnswers EncodeString error", link.LinkAccount, err)
						}
						_, err = dynamicd.SendLinkMessage(link.MyAccount, link.LinkAccount, encoded)
						if err != nil {
							fmt.Println("SendAnswers dynamicd.SendLinkMessage error", link.LinkAccount, err)
						}
						// Set the handler for ICE connection state
						// This will notify you when the peer has connected/disconnected
						link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
							fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
						})
						// Register data channel creation handling
						link.PeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
							fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
							// Register channel opening handling
							d.OnOpen(func() {
								fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

								for range time.NewTicker(5 * time.Second).C {
									message, _ := util.RandomString(24)
									fmt.Printf("Sending '%s'\n", message)

									// Send the message as text
									sendErr := d.SendText(message)
									if sendErr != nil {
										panic(sendErr)
									}
								}
							})

							// Register text message handling
							d.OnMessage(func(msg webrtc.DataChannelMessage) {
								fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
							})
						})
					}
				}
			} else {
				fmt.Println("Error nil PeerConnection", link.LinkAccount)
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
	for {
		select {
		default:
			for _, link := range linkBridges.unconnected {
				select {
				default:
					if link.PeerConnection == nil {
						pc, err := ConnectToIceServices(config)
						if err == nil {
							link.PeerConnection = pc
						}
					}
					if link.PeerConnection != nil {
						answers, err := dynamicd.GetLinkMessages(link.MyAccount, link.LinkAccount)
						if err != nil {
							fmt.Println("GetAnswers error", link.LinkAccount, err)
						} else {
							//fmt.Println("GetAnswers", answers)
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
							link.Answer, err = util.DecodeString(answer.Message)
							if err != nil {
								fmt.Println("GetAnswers DecodeString error", link.LinkAccount, err)
								continue
							}
							//link.Answer = strings.ReplaceAll(answer.Message, `""`, "") // remove double quotes in answer
							if link.PeerConnection == nil {
								fmt.Println("GetAnswers PeerConnection nil for", link.LinkAccount)
								continue
							}
							sd := webrtc.SessionDescription{Type: 2, SDP: link.Answer}
							err := link.PeerConnection.SetRemoteDescription(sd)
							if err != nil {
								fmt.Println("GetAnswers SetRemoteDescription error ", err)
							} else {
								// Set the handler for ICE connection state
								// This will notify you when the peer has connected/disconnected
								link.PeerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
									fmt.Printf("OnICEConnectionStateChange has changed: %s\n", connectionState.String())
								})
								link.PeerConnection.OnICEGatheringStateChange(func(gathererState webrtc.ICEGathererState) {
									fmt.Printf("OnICEGatheringStateChange has changed: %s\n", gathererState.String())
								})
								link.PeerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
									fmt.Printf("OnICECandidate has changed: %s\n", candidate.ToJSON())
								})
								link.PeerConnection.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
									fmt.Printf("OnConnectionStateChange has changed: %s\n", state.String())
								})
								link.PeerConnection.OnSignalingStateChange(func(sig webrtc.SignalingState) {
									fmt.Printf("OnSignalingStateChange has changed: %s\n", sig.String())
								})
								// Register data channel creation handling
								link.PeerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
									fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
									// Register channel opening handling
									d.OnOpen(func() {
										fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

										for range time.NewTicker(5 * time.Second).C {
											message, _ := util.RandomString(12)
											fmt.Printf("Sending '%s'\n", message)

											// Send the message as text
											sendErr := d.SendText(message)
											if sendErr != nil {
												panic(sendErr)
											}
										}
									})

									// Register text message handling
									d.OnMessage(func(msg webrtc.DataChannelMessage) {
										fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
									})
								})
								dc, err := link.PeerConnection.CreateDataChannel(link.LinkAccount, nil)
								if err != nil {
									fmt.Println("GetAnswers CreateDataChannel error", err)
								}
								fmt.Println("GetAnswers Data Channel Negotiated", dc.Negotiated())
							}
						}
					} else {
						fmt.Println("Error nil PeerConnection", link.LinkAccount)
					}
				case <-stopchan:
					fmt.Println("GetAnswers stopped")
					return false
				}
			}
		case <-stopchan:
			fmt.Println("GetAnswers stopped")
			return false
		}
		// check for message every 30 seconds
		time.Sleep(time.Second * 30)
		fmt.Println("GetAnswers restart")
	}
	return true
}
