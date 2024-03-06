package internal

import (
	"log"
	"os"

	"importer/models"
)

var (
	// contains speeds, errors, and warnings
	LOG_FILE = os.Getenv("ERROR_LOG_FILE")
	// contains the bad data
	BAD_DATA_FILE = os.Getenv("BAD_DATA_FILE")
)

// Struct to hold logging information for the application
type Logs struct {
	Info    []string
	Warn    []string
	Error   []string
	BadData []string
}

func (l *Logs) AddBadData(msg models.Message) {
	l.BadData = append(l.BadData, msg.Message)
}

func (l *Logs) AddInfo(msg models.Message) {
	dt := msg.Time.Format("2006/01/02 15:04:05")
	l.Info = append(l.Info, dt+" [INFO] "+msg.Message)
}
func (l *Logs) AddWarn(msg models.Message) {
	dt := msg.Time.Format("2006/01/02 15:04:05")
	l.Warn = append(l.Warn, dt+" [WARN] "+msg.Message)
}
func (l *Logs) AddErr(msg models.Message) {
	dt := msg.Time.Format("2006/01/02 15:04:05")
	l.Error = append(l.Error, dt+" [ERROR] "+msg.Message)
}

func (l *Logs) ClearLogs() {
	l.Info = []string{}
	l.Warn = []string{}
	l.Error = []string{}
	l.BadData = []string{}
}

func (l *Logs) WriteLogs() {
	logFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	for _, msg := range l.Info {
		logFile.WriteString(msg)
	}

	for _, msg := range l.Warn {
		logFile.WriteString(msg)
	}

	for _, msg := range l.Error {
		logFile.WriteString(msg)
	}
	l.ClearLogs()
}

func (l *Logs) WriteBadData() {
	badDataFile, err := os.OpenFile(BAD_DATA_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer badDataFile.Close()
	for _, msg := range l.BadData {
		badDataFile.WriteString(msg + "\n")
	}
	l.ClearLogs()
}
