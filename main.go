package main

import (
	"database/sql"
	"encoding/csv"
	"io/fs"
	"os"
	"sync"
	"time"
)

var (
	db                 *sql.DB
	InfoLog            = InitInfoLogger()
	WarnLog            = InitWarnLogger()
	ErrorLog           = InitErrorLogger()
	wg                 sync.WaitGroup
	InvalidRecordCount = 0
)

const DATA_DIR = "/home/ubuntu/data/"

func main() {
	fChan := make(chan fs.DirEntry)

	db = InitDatabase()
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		ErrorLog.Fatal(pingErr)
	}
	InfoLog.Println("Connected to MySQL Test Database")

	files, err := os.ReadDir(DATA_DIR)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	db.SetMaxOpenConns(len(files))
	// Start a timer that completes after all go routines are done
	begin := time.Now()
	// Start a go routine for each file
	for i := range files {
		go fileToDb(i+1, fChan)
	}

	for _, file := range files {
		fChan <- file
	}
	wg.Wait()
	close(fChan)
	elapsed := time.Since(begin)
	InfoLog.Println("Time for all queries: ", elapsed)
	InfoLog.Println("Invalid records: ", InvalidRecordCount)
}

// fileToDb will read the file and insert the data into the database and log the time it took concurrently
func fileToDb(i int, f chan fs.DirEntry) {
	var d []DeviceData
	wg.Add(1)
	defer wg.Done()
	count := 0
	//InfoLog.Println("Starting thread: ", i)

	// Get the file from the channel
	dirEntry := <-f

	// Open the file
	file, err := os.Open(DATA_DIR + dirEntry.Name())
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
			ErrorLog.Fatal(err)
		}

		data := ParseRecord(record)
		// If the data is empty, skip it
		if data == (DeviceData{}) {
			continue
		}
		d = append(d, data)

	}

	WriteDeviceData(d)
	count += len(d)
	elapsed := time.Since(begin)
	InfoLog.Println("Time for thread ", i, ": ", elapsed)
	InfoLog.Printf("Rows per second: %.2f\n", float64(count)/elapsed.Seconds())
}
