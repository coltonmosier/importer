package cleaner

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"importer/internal"
	"importer/models"
)

/* GLOBALS */
var (
	db            *sql.DB
	mu            sync.Mutex
	Concurrency   int
	SerialNumbers = []string{}
	Logger        = internal.Logs{}

	DATA_DIR = os.Getenv("DATA_DIR")
)

func main() {
	begin := time.Now()

	Concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Concurrency: ", Concurrency)

	dirtyFiles, err := os.ReadDir(DATA_DIR + "/dirty")
	if err != nil {
		log.Fatal(err)
	}

	// fs.DirEntry channel buffered to 5 which means only 5 funcs at a time
	fChan := make(chan fs.DirEntry, Concurrency)

	var d []models.DeviceData
	var invalidRecordCount int

	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range dirtyFiles {
		fChan <- file
		go func() {
			r := internal.FileToStruct(i, file)
			res, invalid := ParseDirtyRecord(r)
			mu.Lock()
			invalidRecordCount += invalid
			d = append(d, res...)
			mu.Unlock()
			<-fChan
		}()
	}
	close(fChan)

	// NOTE: we now have clean data that needs to be written to a clean file

	Logger.AddInfo(models.Message{Message: "size of clean data from files: " + strconv.Itoa(len(d)),
		Time: time.Now()})

	WriteClean(d)

	elapsed := time.Since(begin)
	Logger.AddInfo(models.Message{Message: "Time to process all records:" + elapsed.String() + "\n",
		Time: time.Now()})
	Logger.AddInfo(models.Message{Message: "Invalid records: " + strconv.Itoa(invalidRecordCount) + "\n",
		Time: time.Now()})

	Logger.WriteLogs()
}
