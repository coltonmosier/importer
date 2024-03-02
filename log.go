package main

import (
	"log"
	"os"
)

const (
	INFO_LOG_FILE  = "/home/ubuntu/logs/importer_info.log"
	WARN_LOG_FILE  = "/home/ubuntu/logs/importer_warn.log"
	ERROR_LOG_FILE = "/home/ubuntu/logs/importer_error.log"
)

// InitInfoLogger initializes the info logger should be used for general information
func InitInfoLogger() *log.Logger {
	file, err := os.OpenFile(INFO_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[INFO]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// InitErrorLogger initializes the error logger should be used for fatal errors
func InitErrorLogger() *log.Logger {
	file, err := os.OpenFile(ERROR_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// InitWarnLogger initializes the warn logger should be used for non-fatal errors
func InitWarnLogger() *log.Logger {
	file, err := os.OpenFile(WARN_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[WARN]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
