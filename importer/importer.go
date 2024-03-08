package importer

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"aswe-importer/internal"
	"aswe-importer/models"
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

func importer() {
	begin := time.Now()

	CLEAN_DIR := os.Getenv("CLEAN_DIR")
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
	log.Println("Connected to MySQL Device Database")

	files, err := os.ReadDir(CLEAN_DIR)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Files to process: ", len(files))

	// fs.DirEntry channel buffered to 5 which means only 5 funcs at a time
	fChan := make(chan fs.DirEntry, Concurrency)

	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range files {
		Wg.Add(1)
		fChan <- file
		go func() {
            gTime := time.Now()
			res := internal.FileToStruct(i, file)
			ParseCleanRecord(res)
			<-fChan
            log.Println("File processed: ", i+1)
            Wg.Done()
            fTime := time.Since(gTime)
            Logger.AddInfo(models.Message{Message: "File " + file.Name() + " processed in " + fTime.String() + "\n", Time: time.Now()})
            Logger.AddInfo(models.Message{Message: "Queries per second: " + strconv.FormatFloat(float64(len(res))/fTime.Seconds(), 'f', 2, 64) + "\n", Time: time.Now()})
		}()
	}
	close(fChan)
    Wg.Wait()

	elapsed := time.Since(begin)
	Logger.AddInfo(models.Message{Message: "Total time to process files: " + elapsed.String() + "\n", Time: time.Now()})
	Logger.WriteLogs()
}

func RunImporter() {
	importer()
}
