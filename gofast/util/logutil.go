package util

import (
	"log"
	"os"
)

var infoLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewActorLogger(name string) (*log.Logger, *log.Logger) {
	return log.New(os.Stdout, "INFO: ["+name+"] ", log.Ldate|log.Ltime), log.New(os.Stderr, "ERROR: ["+name+"] ", log.Ldate|log.Ltime)
}

func LogInfo(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func LogFatal(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}
