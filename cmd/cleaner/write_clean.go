package cleaner

import (
	"encoding/csv"
	"importer/models"
	"os"
	"strconv"
)

func WriteClean(d []models.DeviceData) {
	cleanFile := os.Getenv("CLEAN_DATA_FILE")
	f, err := os.OpenFile(cleanFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r [][]string

	// This loop is to convert the struct to a slice of slices of strings
	for _, v := range d {
		r = append(r, []string{strconv.Itoa(v.Line_number), v.Device_type,
			v.Manufacturer, v.Serial_number})
	}

	// Write all the clean data to the clean file
	w := csv.NewWriter(f)
	w.WriteAll(r)

	if err := w.Error(); err != nil {
		panic(err)
	}
}
