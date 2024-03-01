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
func WriteDeviceData(data []DeviceData) {
	// Insert the data into the database using the prepared statement
	stmt, err := db.Prepare("INSERT INTO devices(device_type, manufacturer, serial_number) VALUES(?, ?, ?)")
	if err != nil {
		ErrorLog.Fatalf("%v: creating prepared statement\n", err)
	}

	for _, d := range data {
		_, err := stmt.Exec(d.device_type, d.manufacturer, d.serial_number)
		if err != nil {
			ErrorLog.Fatalf("%v: executing prepared statement %v\n", err, d)
		}
	}
    stmt.Close()
}
