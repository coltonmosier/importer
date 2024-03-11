package cleaner

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
	mu            sync.RWMutex
	Wg            sync.WaitGroup
	Concurrency   int
	SerialNumbers = []string{}
	Logger        = internal.NewLogger()
)

func cleaner() {
	begin := time.Now()

	Concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY"))
	if err != nil {
		log.Fatal(err)
	}
	data_directory := os.Getenv("DATA_DIR")

	log.Println("Cleaner started with concurrency: ", Concurrency)

	dirtyFiles, err := os.ReadDir(data_directory)
	if err != nil {
		log.Fatal(err)
	}

	// fs.DirEntry channel buffered to 5 which means only 5 funcs at a time
	 fChan := make(chan fs.DirEntry, Concurrency)

	var d []models.DeviceData
	var ic []models.InvalidError

	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range dirtyFiles {
		Wg.Add(1)
		fChan <- file
		go func() {
			r := internal.FileToStruct(i, file)
			res, invalid := ParseDirtyRecord(r)
            log.Println("Finished parsing file: ", i)
			mu.Lock()
			ic = append(ic, invalid)
			d = append(d, res...)
			mu.Unlock()
			<-fChan
			Wg.Done()
		}()
	}
	close(fChan)
	Wg.Wait()
    log.Println("All files processed writting clean data")

	// NOTE: we now have clean data that needs to be written to a clean file

	WriteClean(d)

	elapsed := time.Since(begin)
    Logger.WriteLogs()
    Logger.AddInfo(models.Message{Message: "\n", Time: time.Now()})
	Logger.AddInfo(models.Message{Message: "Time to process all records:" + elapsed.String() + "\n",
		Time: time.Now()})

	var invalid models.InvalidError
	var totalInvalid int
	for _, v := range ic {
		invalid.MissingFields += v.MissingFields
		invalid.DeviceTypeMissing += v.DeviceTypeMissing
		invalid.ManufacturerMissing += v.ManufacturerMissing
		invalid.SerialNumberMissing += v.SerialNumberMissing
		invalid.SerialNumberLength += v.SerialNumberLength
		invalid.SerialNumberExists += v.SerialNumberExists
		totalInvalid += v.MissingFields + v.DeviceTypeMissing + v.ManufacturerMissing + v.SerialNumberMissing + v.SerialNumberLength + v.SerialNumberExists
	}

	Logger.AddInfo(models.Message{Message: "Total invalid records: " + strconv.Itoa(totalInvalid) + "\n", Time: time.Now()})
	Logger.AddInfo(models.Message{Message: "size of clean data from files: " + strconv.Itoa(len(d)) + "\n",
		Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with missing fields: " + strconv.Itoa(invalid.MissingFields) + "\n", Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with missing device type: " + strconv.Itoa(invalid.DeviceTypeMissing) + "\n", Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with missing manufacturer: " + strconv.Itoa(invalid.ManufacturerMissing) + "\n", Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with missing serial number: " + strconv.Itoa(invalid.SerialNumberMissing) + "\n", Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with invalid serial number length: " + strconv.Itoa(invalid.SerialNumberLength) + "\n", Time: time.Now()})
    Logger.AddInfo(models.Message{Message: "records with duplicate serial number: " + strconv.Itoa(invalid.SerialNumberExists) + "\n", Time: time.Now()})

	Logger.WriteLogs()
}

func RunCleaner() {
	cleaner()
}
