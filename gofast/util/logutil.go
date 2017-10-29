package util

import (
	"log"
	"os"
)

func NewActorLogger(name string) (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "INFO: ["+name+"] ", log.Ldate|log.Ltime), log.New(os.Stderr, "ERROR: ["+name+"] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewSystemLogger() (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}
