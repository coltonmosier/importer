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
	mu                 sync.Mutex
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

	// fs.DirEntry channel buffered to 5 which means only 5 funcs at a time
	fChan := make(chan fs.DirEntry, 5)

    var d []DeviceData


	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range files {
		fChan <- file
        go func() {
            res := fileToStruct(i, file)
            mu.Lock()
            d = append(d, res...)
            mu.Unlock()
            <-fChan
        }()
	}
	close(fChan)

	log.Println("size of data from files: ", len(d))

	elapsed := time.Since(begin)
    Logger.AddInfo(Message{"Time for all queries: " + elapsed.String(), time.Now()})
    Logger.AddInfo(Message{"Invalid records: " + strconv.Itoa(InvalidRecordCount), time.Now()})
    Logger.WriteLogs()
}
