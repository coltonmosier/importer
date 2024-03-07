package main

import (
	"aswe-importer/cleaner"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file ", err)
	}

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
            // TODO: fix importer
			fmt.Println("Running Importer")
        case "3", "q":
            fmt.Println("Exiting")
            return
		default:
			fmt.Println("Invalid input")
		}
        fmt.Println()

	}
}
