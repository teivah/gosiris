package gosiris

import (
	"log"
	"os"
)

var InfoLogger *log.Logger
var ErrorLogger *log.Logger
var FatalLogger *log.Logger

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	FatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewActorLogger(name string) (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "INFO: ["+name+"] ", log.Ldate|log.Ltime), log.New(os.Stderr, "ERROR: ["+name+"] ", log.Ldate|log.Ltime)
}
