package main

import (
	"log"
	"os"
	"time"
)

const (
	INFO_LOG_FILE  = "/home/ubuntu/logs/importer_info.log"
	WARN_LOG_FILE  = "/home/ubuntu/logs/importer_warn.log"
	ERROR_LOG_FILE = "/home/ubuntu/logs/importer_error.log"
)

// Struct to hold logging information for the application
type Logs struct {
    Info  []string
    Warn  []string
    Error []string
}

type Message struct {
    Message string
    Time time.Time
}

func (l *Logs) AddInfo(msg Message) {
    dt := msg.Time.Format("2006/01/02 15:04:05")
    l.Info = append(l.Info, dt + "[INFO]" + msg.Message)
}
func (l *Logs) AddWarn(msg Message) {
    dt := msg.Time.Format("2006/01/02 15:04:05")
    l.Warn = append(l.Warn, dt + "[WARN]" + msg.Message)
}
func (l *Logs) AddErr(msg Message) {
    dt := msg.Time.Format("2006/01/02 15:04:05")
    l.Error = append(l.Error, dt + "[ERROR]" + msg.Message)
}

func (l *Logs) WriteLogs() {
    infoFile, err := os.OpenFile(INFO_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer infoFile.Close()
    for _, msg := range l.Info {
        infoFile.WriteString(msg + "\n")
    }

    warnFile, err := os.OpenFile(WARN_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer warnFile.Close()
    for _, msg := range l.Warn {
        warnFile.WriteString(msg + "\n")
    }

    errorFile, err := os.OpenFile(ERROR_LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer errorFile.Close()
    for _, msg := range l.Error {
        errorFile.WriteString(msg + "\n")
    }
}
