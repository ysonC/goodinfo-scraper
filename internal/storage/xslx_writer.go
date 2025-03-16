package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func combineAllCSVInFolderToXLSX(folderPath, xlsxOutputPath string) error {
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

func CombineSuccessfulStocks(stocks []string, downloadDir, finalOutputDir string) {
	for _, stock := range stocks {
		stockDir := filepath.Join(downloadDir, stock)
		finalOutput := filepath.Join(finalOutputDir, stock+".xlsx")
		if err := combineAllCSVInFolderToXLSX(stockDir, finalOutput); err != nil {
			log.Printf("Error combining stock %s: %v", stock, err)
			continue
		}
		log.Printf("Successfully combined data for stock %s", stock)
	}
}
