package helper

import (
	"fmt"
	"log"
	"time"
)

func getMaxWorkers() int {
	for {
		fmt.Println("Select max number of workers to download for you:")
		fmt.Println("1. 10 (recommended)\n2. 20\n3. 30\n4. 100")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "1":
			return 10
		case "2":
			return 20
		case "3":
			return 30
		case "4":
			return 100
		default:
			fmt.Println("Invalid option, try again.")
		}
	}
}

func getDateRange() (string, string) {
	fmt.Println("Select date range:\n1. Max years\n2. Custom range")
	var input string
	fmt.Scanln(&input)
	switch input {
	case "1":
		return "1965-01-01", time.Now().Format("2006-01-02")
	case "2":
		var start, end string
		fmt.Print("Start (YYYY-MM-DD): ")
		fmt.Scanln(&start)
		fmt.Print("End (YYYY-MM-DD): ")
		fmt.Scanln(&end)
		return start, end
	default:
		log.Fatal("Invalid date range option")
		return "", ""
	}
}
