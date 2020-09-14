package webbridge

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/duality-solutions/web-bridge/internal/util"
	"github.com/duality-solutions/web-bridge/rest"
	"github.com/duality-solutions/web-bridge/rpc/dynamic"
)

func appCommandLoop(shutdown *rest.AppShutdown, d *dynamic.Dynamicd, status *dynamic.SyncStatus, sync bool) {
	go func() {
		var err error
		for {
			select {
			default:
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("web-bridge> ")
				cmdText, _ := reader.ReadString('\n')
				if len(cmdText) > 1 {
					cmdText = strings.Trim(cmdText, "\r\n ") //cmdText[:len(cmdText)-2]
				}
				if strings.HasPrefix(cmdText, "exit") || strings.HasPrefix(cmdText, "shutdown") || strings.HasPrefix(cmdText, "stop") {
					util.Info.Println("Exit command. Stopping services.")
					shutdown.ShutdownAppliction()
					break
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
			case <-*shutdown.Close:
				return
			}
		}
	}()
}
