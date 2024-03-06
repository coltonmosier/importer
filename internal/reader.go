package internal

import (
	"encoding/csv"
	"io/fs"
	"log"
	"os"
)

var DATA_DIR = os.Getenv("DATA_DIR")

func FileToStruct(i int, fc fs.DirEntry) [][]string{



	var re [][]string

	// Open the file
	file, err := os.Open(DATA_DIR + fc.Name())
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

    return re
}
