package importer

import (
	"importer/models"
	"strconv"
)

// ParseCleanData parses the clean data and returns a slice of DeviceData
func ParseCleanRecord(record [][]string) []models.DeviceData {
    var d []models.DeviceData
    for _, v := range record {
        ln, _ := strconv.Atoi(v[0])
        d = append(d, models.DeviceData{
            Line_number:   ln,
            Device_type:   v[1],
            Manufacturer:  v[2],
            Serial_number: v[3],
        })
    }
    return d
}
