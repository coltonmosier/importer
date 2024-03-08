package internal

import (
	"encoding/csv"
	"io/fs"
	"log"
	"os"
)

var DATA_DIR = "data/dirty/"
var CLEAN_DIR = "data/clean/"

func FileToStruct(i int, fc fs.DirEntry) [][]string{



	var re [][]string

	// Open the file
	file, err := os.Open(CLEAN_DIR + fc.Name())
	if err != nil {
        log.Println("Error opening file in filetostruct: ", fc.Name())
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
