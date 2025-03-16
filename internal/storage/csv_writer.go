package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// SaveToCSV writes a 2D string slice to a CSV file.
func SaveToCSV(data [][]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Optionally write a header row if needed.
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Printf("Error writing row: %v", err)
		}
	}
	return nil
}

// IsFileUpToDate checks if the file exists and was modified today.
func IsFileUpToDate(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	now := time.Now()
	modTime := info.ModTime()
	return now.Year() == modTime.Year() && now.YearDay() == modTime.YearDay()
}

func combineTwoCSV(filePath1, filePath2 string) [][]string {
	file1, err := os.Open(filePath1)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file1.Close()

	file2, err := os.Open(filePath2)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file2.Close()

	reader1 := csv.NewReader(file1)
	reader2 := csv.NewReader(file2)

	var combineFiles [][]string
	for {
		firstRowFile1, err1 := reader1.Read()
		firstRowFile2, err2 := reader2.Read()

		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		if err1 == io.EOF {
			firstRowFile1 = []string{"-"}
		} else if err1 != nil {
			log.Fatalf("Error reading file1: %v", err1)
		}
		if err2 == io.EOF {
			firstRowFile2 = []string{"-"}
		} else if err2 != nil {
			log.Fatalf("Error reading file2: %v", err2)
		}

		combineRows := append(firstRowFile1, "")
		combineRows = append(combineRows, firstRowFile2...)
		combineFiles = append(combineFiles, combineRows)
	}
	return combineFiles
}

func CombineAllCSVInFolder(folderPath string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	checkList := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	// Ensure all required files are present.
	for _, check := range checkList {
		found := false

		for _, file := range files {
			if strings.Contains(file.Name(), check) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing file: %s", check)
		}
	}

	// Combine all files.
	perAndStockData := combineTwoCSV(folderPath+"/per.csv", folderPath+"/stockdata.csv")
	SaveToCSV(perAndStockData, folderPath+"/combine.csv")
	monthlyRevenueAndCashflow := combineTwoCSV(
		folderPath+"/monthlyrevenue.csv",
		folderPath+"/cashflow.csv",
	)
	SaveToCSV(monthlyRevenueAndCashflow, folderPath+"/combine2.csv")
	return nil
}
