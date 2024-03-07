package importer

import (
	"aswe-importer/models"
)

// ParseCleanData parses the clean data and returns a slice of DeviceData
func ParseCleanRecord(record [][]string) []models.DeviceData {
    var d []models.DeviceData
    for _, v := range record {
        d = append(d, models.DeviceData{
            Line_number:   v[0],
            Device_type:   v[1],
            Manufacturer:  v[2],
            Serial_number: v[3],
        })
    }
    return d
}
