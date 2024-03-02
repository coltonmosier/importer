package main

import (
	"encoding/csv"
	"io/fs"
	"log"
	"os"
)

func fileToStruct(i int, fc <-chan fs.DirEntry, dChan chan<- [][]DeviceData) {

	f := <-fc
    defer wg.Done()


	var d []DeviceData
	var re [][]string

	// Open the file
	file, err := os.Open(DATA_DIR + f.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
		re = append(re, record)
	}

    log.Println("Time to parse:", i)
    d = ParseRecord(re)
    res := [][]DeviceData{d}

	dChan <- res
}
