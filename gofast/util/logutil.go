package util

import (
	"log"
	"os"
)

var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func init() {
	InfoLogger, ErrorLogger = newSystemLogger()
}

func NewActorLogger(name string) (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "INFO: ["+name+"] ", log.Ldate|log.Ltime), log.New(os.Stderr, "ERROR: ["+name+"] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func newSystemLogger() (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}
