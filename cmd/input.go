package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getStockNumbers(folderPath string) []string {
	var stocks []string
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read stock folder: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(folderPath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file %s: %v", path, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				parts := strings.Split(line, ",")
				stocks = append(stocks, strings.TrimSpace(parts[0]))
			}
		}
		f.Close()
	}
	if len(stocks) == 0 {
		log.Fatalf("No stock numbers found in %s", folderPath)
	}
	return stocks
}

func promptMaxWorkers() int {
	options := map[string]int{"1": 10, "2": 20, "3": 30, "4": 100}
	for {
		fmt.Println("Select max workers:")
		fmt.Println(
			"1. 10  (recommended)\n2. 20  (maybe recommended)\n3. 30  (maybe not recommended)\n4. 100 (are you sure?)",
		)
		var choice string
		fmt.Scanln(&choice)
		if workers, ok := options[choice]; ok {
			return workers
		}
		fmt.Println("Invalid option, try again.")
	}
}

func promptDateRange() (string, string) {
	fmt.Println("Select date range:\n1. Max\n2. Custom")
	var choice string
	fmt.Scanln(&choice)
	switch choice {
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
		log.Fatal("Invalid date range option selected")
		return "", ""
	}
}
