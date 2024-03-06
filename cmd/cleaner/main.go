package cleaner

import (
	"database/sql"
	"fmt"
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
	var ic []models.InvalidError

	// Loop through the files and send them to the channel
	// acts like a semaphore
	for i, file := range dirtyFiles {
		fChan <- file
		go func() {
			r := internal.FileToStruct(i, file)
			res, invalid := ParseDirtyRecord(r)
			mu.Lock()
            ic = append(ic, invalid)
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
    invalidMsg := fmt.Sprintf("MissingFields: %d\nDeviceTypeMissing: %d\nManufacturerMissing: %d\nSerialNumberMissing: %d\nSerialNumberLength: %d\nSerialNumberExists: %d\n",
        invalid.MissingFields, invalid.DeviceTypeMissing,
        invalid.ManufacturerMissing, invalid.SerialNumberMissing, 
        invalid.SerialNumberLength, invalid.SerialNumberExists)
    Logger.AddInfo(models.Message{Message: invalidMsg, Time: time.Now()})
    

	Logger.WriteLogs()
}
