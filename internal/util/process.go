package util

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/process"
)

// WaitForProcessShutdown waits for a process to shutdown normally or kills it if it runs past the given timeout
func WaitForProcessShutdown(process *os.Process, timeout time.Duration) bool {
	Info.Printf("WaitForStoppedPID waiting for process pid %v to shutdown. Timeout set to %v seconds\n", process.Pid, timeout.Seconds())
	_, err := os.FindProcess(process.Pid)
	if err != nil {
		Info.Println("WaitForStoppedPID process found. Waiting for normal shutdown or", timeout.String(), "seconds.")
		for {
			select {
			case <-time.After(time.Second * 3):
				_, err = os.FindProcess(process.Pid)
				if err != nil {
					Info.Println("WaitForStoppedPID process not found anymore. Normal shutdown complete.")
					return true
				}
			case <-time.After(timeout):
				Info.Printf("WaitForStoppedPID timeout expired after %v. Killing process!\n", timeout.String())
				if errKill := process.Kill(); errKill != nil {
					Error.Println("WaitForStoppedPID failed to kill process after timeout ", errKill)
					return false
				}
				return true
			}
		}
	} else {
		Info.Println("WaitForStoppedPID process not found")
		return true
	}
}

// FindWebBridgeProcess returns true if a web-bridge process is already running
func FindWebBridgeProcess() (bool, int32, error) {
	processList, err := process.Processes()
	if err != nil {
		return false, -1, fmt.Errorf("process.Processes() Failed")
	}
	thisPID := int32(os.Getpid())
	for _, process := range processList {
		if process.Pid != thisPID {
			name, _ := process.Name()
			if strings.Contains(name, "web-bridge") {
				fmt.Println("FindWebBridgeProcess found", name)
				return true, process.Pid, nil
			}
		}
	}
	return false, -1, nil
}
