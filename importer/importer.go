package importer

import (
	"database/sql"
	"log"
	"os"
	"sync"
	"time"

	"aswe-importer/internal"
	"aswe-importer/models"
)

/* GLOBALS */
var (
	db            *sql.DB
	Wg            sync.WaitGroup
	SerialNumbers = []string{}
	Logger        = internal.Logs{}
)

func importer() {
	begin := time.Now()

	CLEAN_DIR := os.Getenv("CLEAN_DIR")
	db = internal.InitDatabase()
	defer db.Close()
	db.SetMaxOpenConns(75)
	db.SetMaxIdleConns(10)

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr.Error())
	}
	log.Println("Connected to MySQL Device Database")

	files, err := os.ReadDir(CLEAN_DIR)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Files to process: ", len(files))

	for i, file := range files {
		Wg.Add(1)
		go func() {
			res := internal.FileToStruct(i, file)
			ParseCleanRecord(res)
			Wg.Done()
		}()
	}
    log.Println("All files read -- waiting on queries...")
	Wg.Wait()

	elapsed := time.Since(begin)
	Logger.AddInfo(models.Message{Message: "Total time to process files: " + elapsed.String() + "\n", Time: time.Now()})
	Logger.WriteLogs()
}

func RunImporter() {
	importer()
}
