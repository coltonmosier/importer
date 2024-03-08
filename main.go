package main

import (
	"aswe-importer/cleaner"
	"aswe-importer/importer"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}
    fmt.Println("All .env variables loaded")
    cf := os.Getenv("CLEAN_DATA_FILE")
    rlf := os.Getenv("ERROR_LOG_FILE")
    bdf := os.Getenv("BAD_DATA_FILE")
    dd := os.Getenv("DATA_DIR")
    db := os.Getenv("MYSQL_DB")
    cd := os.Getenv("CLEAN_DIR")
    fmt.Println("Clean data file: ", cf)
    fmt.Println("Error log file: ", rlf)
    fmt.Println("Bad data file: ", bdf)
    fmt.Println("Data directory: ", dd)
    fmt.Println("Clean directory: ", cd)
    fmt.Println("MySQL DB: ", db)
    fmt.Println()



	for {
		fmt.Println("1. Run Cleaner")
		fmt.Println("2. Run Importer")
        fmt.Println("3/q. Exit")
        fmt.Print("Enter option: ")
        var input string
        fmt.Scanln(&input)
		switch input {
		case "1":
			fmt.Println("Running Cleaner")
            cleaner.RunCleaner()
		case "2":
			fmt.Println("Running Importer")
            importer.RunImporter()
        case "3", "q":
            fmt.Println("Exiting")
            return
		default:
			fmt.Println("Invalid input")
		}
        fmt.Println()

	}
}
