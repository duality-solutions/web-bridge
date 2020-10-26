package util

import (
	"fmt"
	"io"
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

// HomeDir stores the operating system specific user home directory
var HomeDir string = ""

// InitDebugLogFile is used to initialize the debug log file path.
func InitDebugLogFile(console bool, homeDir string) {
	HomeDir = homeDir
	file, err := os.OpenFile(HomeDir+"debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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

// EndDebugLogFile adds x blank lines to debug.log file
func EndDebugLogFile(x int) {
	file, err := os.OpenFile(HomeDir+"debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open log file", err)
	}
	exitLog := log.New(file, "", 0)
	for n := 0; n < x; n++ {
		exitLog.Println()
	}
}

// Println logs to debug.log file and console
func (d *DebugLog) Println(a ...interface{}) {
	go func() {
		d.Logger.Println(a...)
		if d.Console {
			fmt.Println(a...)
		}
	}()
}

// Print logs to debug.log file and console
func (d *DebugLog) Print(a ...interface{}) {
	go func() {
		d.Logger.Print(a...)
		if d.Console {
			fmt.Print(a...)
		}
	}()
}

// Printf logs to debug.log file and console
func (d *DebugLog) Printf(format string, a ...interface{}) {
	go func() {
		d.Logger.Printf(format, a...)
		if d.Console {
			fmt.Printf(format, a...)
		}
	}()
}

// Fprintf logs to debug.log file and console
func (d *DebugLog) Fprintf(w io.Writer, format string, a ...interface{}) {
	go func() {
		d.Logger.Printf(format, a...)
		if d.Console {
			fmt.Fprintf(w, format, a...)
		}
	}()
}
