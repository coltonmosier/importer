package main

import (
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type DeviceData struct {
	device_type   string
	manufacturer  string
	serial_number string
}

func InitDatabase() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		ErrorLog.Fatal("Error loading .env file ", err)
	}

	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               os.Getenv("MYSQL_DB"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		ErrorLog.Fatal(err)
	}

	return db
}

// WriteDeviceData writes the device data to the database from the csv file
func WriteDeviceData(data DeviceData) {
	rows, err := db.Query("INSERT INTO devices(device_type, manufacturer, serial_number) VALUES(?, ?, ?)", data.device_type, data.manufacturer, data.serial_number)
	defer rows.Close()

	if err != nil {
		errString := "INSERT INTO devices(device_type, manufacturer, serial_number) VALUES(" + data.device_type + ", " + data.manufacturer + ", " + data.serial_number + ")"
		ErrorLog.Fatal(err, "\n\t", errString)
	}

}
