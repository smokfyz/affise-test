package log

import (
	"io"
	"log"
	"os"
	"strings"
)

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	Init("WARNING")
}

func Init(logLevel string) {
	logLevel = strings.ToUpper(logLevel)

	Trace = log.New(io.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(io.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(io.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(io.Discard, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.Discard, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	switch logLevel {
	case "TRACE":
		Trace.SetOutput(os.Stdout)
		fallthrough
	case "DEBUG":
		Debug.SetOutput(os.Stdout)
		fallthrough
	case "INFO":
		Info.SetOutput(os.Stdout)
		fallthrough
	case "WARNING":
		Warning.SetOutput(os.Stdout)
		fallthrough
	case "ERROR":
		Warning.SetOutput(os.Stderr)
	default:
		panic("unsupported log level")
	}
}
