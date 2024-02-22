package main

import (
	"log"
	"os"
)

const LOG_FILE = "/home/ubuntu/logs/sql_write.log"

// InitInfoLogger initializes the info logger should be used for general information
func InitInfoLogger() *log.Logger {
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[INFO]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// InitErrorLogger initializes the error logger should be used for fatal errors
func InitErrorLogger() *log.Logger {
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

// InitWarnLogger initializes the warn logger should be used for non-fatal errors
func InitWarnLogger() *log.Logger {
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(file, "[WARN]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}
