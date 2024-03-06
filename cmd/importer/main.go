package importer

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
	Wg            sync.WaitGroup
	mu            sync.Mutex
	Concurrency   int
	SerialNumbers = []string{}
	Logger        = internal.Logs{}
)

const DATA_DIR = "/home/ubuntu/data/"

func main() {
	begin := time.Now()

	db = internal.InitDatabase()
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
	fChan := make(chan fs.DirEntry, Concurrency)

	var d []models.DeviceData
	var InvalidRecordCount int

	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range files {
		fChan <- file
		go func() {
			res, cnt := internal.FileToStruct(i, file)
			mu.Lock()
			InvalidRecordCount += cnt
			d = append(d, res...)
			mu.Unlock()
			<-fChan
		}()
	}
	close(fChan)
	// Stop the CPU profiler

	log.Println("size of data from files: ", len(d))

	elapsed := time.Since(begin)
	Logger.AddInfo(models.Message{Message: "Time to process files:" + elapsed.String() + "\n", Time: time.Now()})
	Logger.AddInfo(models.Message{Message: "Invalid records: " + strconv.Itoa(InvalidRecordCount) + "\n", Time: time.Now()})
	Logger.WriteLogs()
}
