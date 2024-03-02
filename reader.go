package main

import (
	"encoding/csv"
	"io/fs"
	"log"
	"os"
)

func fileToStruct(i int, fc <-chan fs.DirEntry, dChan chan<- []DeviceData) {

	f := <-fc

	log.Println("Processing file: ", i)

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

	log.Println("file:", i, "records read", len(re))

    d = ParseRecord(re)

	log.Println("records ready to write", len(d))

	dChan <- d
	wg.Done()
}
