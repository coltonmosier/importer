package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
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
	Runs               = 1
	Concurrency        int
    SerialNumbers = []string{}
    Logger = Logs{}
)

const DATA_DIR = "/home/ubuntu/data/"

func main() {
	db = InitDatabase()
	defer db.Close()
    Concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Concurrency: ", Concurrency)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
    Logger.AddInfo(Message{Message: "Connected to MySQL Test Database", Time: time.Now()})

	files, err := os.ReadDir(DATA_DIR)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	fChan := make(chan fs.DirEntry, 5)

	begin := time.Now()

	for i := 0; i < Concurrency; i++ {
        wg.Add(1)
		go func() {
			for file := range fChan {
				fileToDb(file)
			}
		}()
	}

	// Loop through the files and send them to the channel
	for _, file := range files {
		fChan <- file
	}
	close(fChan)

	wg.Wait()

	elapsed := time.Since(begin)
    Logger.AddInfo(Message{Message: "Time for all queries: " + elapsed.String(), Time: time.Now()})
    Logger.AddInfo(Message{Message: "Invalid records: " + strconv.Itoa(InvalidRecordCount), Time: time.Now()})
    Logger.WriteLogs()
}

// fileToDb will read the file and insert the data into the database and log the time it took concurrently
func fileToDb(f fs.DirEntry) {
	var d []DeviceData
	count := 0

	// Open the file
	file, err := os.Open(DATA_DIR + f.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Start a timer for the file
	begin := time.Now()
	// Read the csv file
	r := csv.NewReader(file)
	// Parse the file
	for {
		// Read the record aka the line
		record, err := r.Read()
		if err != nil {
			// If the error is EOF, break the loop because the file is done
			if err.Error() == "EOF" {
				break
			}
		}

		data := ParseRecord(record)
		// If the data is empty, skip it
		if data == (DeviceData{}) {
			continue
		}
		d = append(d, data)

		if len(d) == 1250 {
			WriteDeviceData(d)
			count += len(d)
			d = nil
		}

	}

	WriteDeviceData(d)
	count += len(d)
	elapsed := time.Since(begin)
    msg := fmt.Sprintf("Rows per second: %f in %s elapsed time on run %d", float64(count)/elapsed.Seconds(), elapsed.String(), Runs)
    Logger.AddInfo(Message{Message: msg, Time: time.Now()})
    log.Println(Runs, "completed")
	Runs++
	wg.Done()
}
