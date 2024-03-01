package main

import (
	"database/sql"
	"encoding/csv"
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
	InfoLog            = InitInfoLogger()
	WarnLog            = InitWarnLogger()
	ErrorLog           = InitErrorLogger()
	wg                 sync.WaitGroup
	InvalidRecordCount = 0
	Runs               = 1
	Concurrency        int
)

const DATA_DIR = "/home/ubuntu/data/"

func main() {
	db = InitDatabase()
	defer db.Close()
    Concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
    if err != nil {
        ErrorLog.Fatal(err)
    }

    log.Println("Concurrency: ", Concurrency)

	pingErr := db.Ping()
	if pingErr != nil {
		ErrorLog.Fatal(pingErr)
	}
	InfoLog.Println("Connected to MySQL Test Database")

	files, err := os.ReadDir(DATA_DIR)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	fChan := make(chan fs.DirEntry, 5)

	begin := time.Now()
	for i := 0; i < Concurrency; i++ {
		go func() {
			for file := range fChan {
				fileToDb(file)
				defer wg.Done()
			}
		}()
	}

	// Loop through the files and send them to the channel
	for _, file := range files {
		fChan <- file
	}
	close(fChan)

	wg.Add(2)
	wg.Wait()

	elapsed := time.Since(begin)
	InfoLog.Println("Time for all queries: ", elapsed)
	InfoLog.Println("Invalid records: ", InvalidRecordCount)
}

// fileToDb will read the file and insert the data into the database and log the time it took concurrently
func fileToDb(f fs.DirEntry) {
	var d []DeviceData
	count := 0

	// Open the file
	file, err := os.Open(DATA_DIR + f.Name())
	if err != nil {
		ErrorLog.Fatal(err)
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

		if len(d) == 1000 {
			WriteDeviceData(d)
			count += len(d)
			d = nil
		}

	}

	WriteDeviceData(d)
	count += len(d)
	elapsed := time.Since(begin)
	InfoLog.Printf("Rows per second: %.2f in %v elapsed time on run %v\n", float64(count)/elapsed.Seconds(), elapsed, Runs)
	Runs++
}
