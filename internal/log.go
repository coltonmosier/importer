package internal

import (
	"log"
	"os"

	"importer/models"
)

var (
    LOG_FILE = os.Getenv("ERROR_LOG_FILE")
    BAD_DATA_FILE = os.Getenv("BAD_DATA_FILE")
)

// Struct to hold logging information for the application
type Logs struct {
	Info  []string
	Warn  []string
	Error []string
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
}

func (l *Logs) WriteLogs() {
	infoFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer infoFile.Close()
	for _, msg := range l.Info {
		infoFile.WriteString(msg)
	}

	warnFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer warnFile.Close()
	for _, msg := range l.Warn {
		warnFile.WriteString(msg)
	}

	errorFile, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer errorFile.Close()
	for _, msg := range l.Error {
		errorFile.WriteString(msg)
	}
    l.ClearLogs()
}
