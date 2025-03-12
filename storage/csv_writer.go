package storage

import (
	"encoding/csv"
	"log"
	"os"
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
