package dynamic

import (
	"fmt"
	"strings"

	ps "github.com/shirou/gopsutil/process"
)

// FindDynamicdProcess returns the dynamicd processes running locally
func FindDynamicdProcess() (*ps.Process, error) {
	processList, err := ps.Processes()
	if err != nil {
		return nil, fmt.Errorf("ps.Processes() Failed")
	}
	for _, process := range processList {
		name, _ := process.Name()
		if strings.HasPrefix(name, "dynamicd") {
			fmt.Println("Running dynamicd process found", name)
			cmd, _ := process.Cmdline()
			// TODO check datadir as well
			if strings.Index(cmd, "-port=33334") > 0 && strings.Index(cmd, "-rpcport=33335") > 0 {
				return process, nil
			}
			fmt.Println("Incorrect dynamicd process cmd", cmd)
		}
	}
	return nil, fmt.Errorf("Running dynamicd process not found")
}
