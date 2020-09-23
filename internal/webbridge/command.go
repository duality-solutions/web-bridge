package webbridge

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/duality-solutions/web-bridge/api/rest"
	"github.com/duality-solutions/web-bridge/blockchain/rpc/dynamic"
	"github.com/duality-solutions/web-bridge/bridge"
	"github.com/duality-solutions/web-bridge/internal/util"
	"golang.org/x/crypto/ssh/terminal"
)

func unlockWallet(d *dynamic.Dynamicd) bool {
	fmt.Print("wallet passphrase> ")
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	walletpassphrase = strings.Trim(string(bytePassword), "\r\n ")
	err := d.UnlockWallet(walletpassphrase)
	if err == nil {
		util.Info.Println("Wallet unlocked.")
		return true
	}
	util.Error.Println(err)
	return false
}

func appCommandLoop(stopBridges *chan struct{}, acc *[]dynamic.Account, al *dynamic.ActiveLinks,
	shutdown *rest.AppShutdown, d *dynamic.Dynamicd,
	status *dynamic.SyncStatus, sync bool) {
	go func() {
		var err error
		var unlocked, bridgesStarted = false, false
		for {
			select {
			default:
				if !unlocked {
					errUnlock := d.UnlockWallet("")
					if errUnlock != nil {
						util.Info.Println("Wallet locked so links are unavailable. Use the unlock command to start your link bridges.")
					} else {
						unlocked = true
						util.Info.Println("Starting bridges.")
						go bridge.StartBridges(stopBridges, config, *d, *acc, *al)
						bridgesStarted = true
					}
				}
				reader := bufio.NewReader(os.Stdin)
				fmt.Print(DefaultName + `> `)
				cmdText, _ := reader.ReadString('\n')
				if len(cmdText) > 1 {
					cmdText = strings.Trim(cmdText, "\r\n ") //cmdText[:len(cmdText)-2]
				}
				if strings.HasPrefix(cmdText, "exit") || strings.HasPrefix(cmdText, "shutdown") || strings.HasPrefix(cmdText, "stop") {
					util.Info.Println("Exit command. Stopping services.")
					shutdown.ShutdownAppliction()
					break
				} else if strings.HasPrefix(cmdText, "unlock") {
					unlocked = unlockWallet(d)
					if unlocked {
						al, err = d.GetActiveLinks(time.Second * 120)
						if err != nil {
							util.Error.Println("GetActiveLinks error", err)
						} else {
							util.Info.Printf("Found %v links\n", len(al.Links))
						}
					}
				} else if strings.HasPrefix(cmdText, "dynamic-cli") {
					req, errNewRequest := dynamic.NewRequest(cmdText)
					if errNewRequest != nil {
						util.Error.Println("Error:", errNewRequest)
					} else {
						strResp, _ := util.BeautifyJSON(<-d.ExecCmdRequest(req))
						util.Info.Println(strResp)
					}
				} else {
					util.Warning.Println("Invalid command", cmdText)
					status, err = d.GetSyncStatus()
				}
				status, err = d.GetSyncStatus()
				if err != nil {
					util.Error.Println("syncstatus unmarshal error", err)
				} else {
					if !sync {
						util.Info.Println("Sync " + fmt.Sprintf("%f", status.SyncProgress*100) + " percent complete!")
						if status.SyncProgress == 1 {
							sync = true
						}
					}
				}
				if unlocked && !bridgesStarted && acc != nil && al != nil {
					util.Info.Println("Starting bridges.")
					go bridge.StartBridges(stopBridges, config, *d, *acc, *al)
					bridgesStarted = true
				}
			case <-*shutdown.Close:
				return
			}
		}
	}()
}
