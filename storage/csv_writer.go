package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
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

func CombineAllCSVInFolderToXLSX(folderPath, xlsxOutputPath string) error {
	// Define required keywords.
	checkList := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	// Read folder files.
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Ensure each required keyword is found in at least one file name.
	for _, check := range checkList {
		found := false
		for _, file := range files {
			if !file.IsDir() &&
				strings.Contains(strings.ToLower(file.Name()), strings.ToLower(check)) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing file containing: %s", check)
		}
	}

	// Combine the CSV data from the two groups.
	perAndStockData := combineTwoCSV(
		filepath.Join(folderPath, "per.csv"),
		filepath.Join(folderPath, "stockdata.csv"),
	)
	monthlyRevenueAndCashflow := combineTwoCSV(
		filepath.Join(folderPath, "monthlyrevenue.csv"),
		filepath.Join(folderPath, "cashflow.csv"),
	)

	// Create a new XLSX file.
	f := excelize.NewFile()
	// Sheet1: write perAndStockData.
	sheet1 := f.GetSheetName(f.GetActiveSheetIndex())
	for i, row := range perAndStockData {
		for j, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return fmt.Errorf("failed to get cell name: %v", err)
			}
			if err := f.SetCellValue(sheet1, cellName, cellValue); err != nil {
				return fmt.Errorf("failed to set cell value: %v", err)
			}
		}
	}

	// Sheet2: create and write monthlyRevenueAndCashflow.
	sheet2Name := "Sheet2"
	_, err = f.NewSheet(sheet2Name)
	if err != nil {
		return fmt.Errorf("failed to create new sheet: %v", err)
	}
	for i, row := range monthlyRevenueAndCashflow {
		for j, cellValue := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				return fmt.Errorf("failed to get cell name: %v", err)
			}
			if err := f.SetCellValue(sheet2Name, cellName, cellValue); err != nil {
				return fmt.Errorf("failed to set cell value: %v", err)
			}
		}
	}

	// Set active sheet to Sheet1.
	idx, err := f.GetSheetIndex(sheet1)
	if err != nil {
		return fmt.Errorf("failed to set active sheet: %v", err)
	}
	f.SetActiveSheet(idx)

	// Save the Excel file.
	if err := f.SaveAs(xlsxOutputPath); err != nil {
		return fmt.Errorf("failed to save XLSX file: %v", err)
	}

	fmt.Println("Combined data written to XLSX file:", xlsxOutputPath)
	return nil
}
