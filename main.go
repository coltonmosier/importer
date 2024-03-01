package main

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

/* GLOBALS */
var (
	db                 *sql.DB
	wg                 sync.WaitGroup
	InvalidRecordCount = 0
	Concurrency        int
	SerialNumbers      = []string{}
	Logger             = Logs{}
)

const DATA_DIR = "/home/ubuntu/data/"

func main() {
	begin := time.Now()

	db = InitDatabase()
	defer db.Close()
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	Concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Concurrency: ", Concurrency)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected to MySQL Test Database")

	files, err := os.ReadDir(DATA_DIR)
	if err != nil {
		log.Fatal(err)
	}

	fChan := make(chan fs.DirEntry, 5)
	dChan := make(chan []DeviceData)

    for i := range files {
        wg.Add(1)
        go fileToStruct(i+1,fChan, dChan)
    }
	// Loop through the files and send them to the channel
	for _, file := range files {
		fChan <- file
	}
	close(fChan)

	d := <-dChan
	close(dChan)
	wg.Wait() // we know all files have been read and processed

	log.Println("size of data from files: ", len(d))

	elapsed := time.Since(begin)
	//Logger.AddInfo(Message{Message: "Time for all queries: " + elapsed.String(), Time: time.Now()})
	//Logger.AddInfo(Message{Message: "Invalid records: " + strconv.Itoa(InvalidRecordCount), Time: time.Now()})
	Logger.AddInfo(Message{Message: "Time to process all files: " + elapsed.String(), Time: time.Now()})
	Logger.WriteLogs()
}
