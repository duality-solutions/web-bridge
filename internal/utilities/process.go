package util

import (
	"os"
	"time"
)

// WaitForStoppedPID waits for PID process to stop or it kills it after the time period
func WaitForStoppedPID(process *os.Process, timeout time.Duration) bool {
	_, err := os.FindProcess(process.Pid)
	if err != nil {
		Info.Println("WaitForStoppedPID process found. Waiting for normal shutdown.")
		for {
			select {
			case <-time.After(time.Second * 3):
				_, err = os.FindProcess(process.Pid)
				if err != nil {
					Info.Println("WaitForStoppedPID process not found anymore. Shutdown complete.")
					return true
				}
			case <-time.After(timeout):
				Info.Println("WaitForStoppedPID timeout expired. Killing process.")
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
