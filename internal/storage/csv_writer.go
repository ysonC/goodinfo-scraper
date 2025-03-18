package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

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

func ReadDirFiles(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteCSV(filepath string, data [][]string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func CheckFileExist(fileNames []string) error {
	checkList := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	for _, check := range checkList {
		found := false

		for _, file := range fileNames {
			if strings.Contains(file, check) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing file: %s", check)
		}
	}
	return nil
}

func MergeCSVData(csv1, csv2 [][]string) ([][]string, error) {
	if len(csv1) == 0 || len(csv2) == 0 {
		return nil, fmt.Errorf("empty csv data")
	}

	var merged [][]string
	maxRows := max(len(csv1), len(csv2))
	for i := range maxRows {
		row1 := []string{"-"}
		row2 := []string{"-"}
		if i < len(csv1) {
			row1 = csv1[i]
		}
		if i < len(csv2) {
			row2 = csv2[i]
		}
		combinedRows := append(row1, "")
		combinedRows = append(combinedRows, row2...)
		merged = append(merged, combinedRows)
	}
	return merged, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
