package webbridge

import (
	"fmt"
	"log"
	"os"
)

// DebugLog handles console and file debug print lines.
type DebugLog struct {
	Console bool
	*log.Logger
}

var (
	// Trace used to log trace events
	Trace DebugLog
	// Info used to log general information
	Info DebugLog
	// Warning used to log warning information
	Warning DebugLog
	// Error used to log error information
	Error DebugLog
)

// InitDebugLogFile is used to initialize the debug log file path.
func InitDebugLogFile(console bool) {

	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file", err)
	}

	Trace.Console = console
	Trace.Logger = log.New(file,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info.Console = console
	Info.Logger = log.New(file,
		"INFO: ",
		log.Ldate|log.LUTC|log.Ltime)

	Warning.Console = console
	Warning.Logger = log.New(file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error.Console = console
	Error.Logger = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

// Println logs to debug.log file and console
func (d *DebugLog) Println(a ...interface{}) {
	d.Logger.Println(a...)
	if d.Console {
		fmt.Println(a...)
	}
}

// Print logs to debug.log file and console
func (d *DebugLog) Print(a ...interface{}) {
	d.Logger.Print(a...)
	if d.Console {
		fmt.Print(a...)
	}
}
