package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ysonC/multi-stocks-download/storage"
)

// func checkSpace(record []string) {
// 	if slices.Contains(record, "") {
// 		fmt.Println("Empty value found,")
// 		for i, v := range record {
// 			if v == "" {
// 				fmt.Printf("Empty value found at index %d\n", i)
// 			}
// 		}
// 		return
// 	}
// 	fmt.Println("No empty value found")
// }

func readCSVFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	counterInt := 0

	// Read line by line.
	for {
		_, err := reader.Read()
		if err != nil {
			break
		}

		counterInt++
	}
	fmt.Printf("Total lines: %d\n", counterInt)
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

func combineAllCSVInFolder(folderPath string) error {
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
	storage.SaveToCSV(perAndStockData, folderPath+"/combine.csv")
	fmt.Println("Combined per and stock data")
	monthlyRevenueAndCashflow := combineTwoCSV(
		folderPath+"/monthlyrevenue.csv",
		folderPath+"/cashflow.csv",
	)
	storage.SaveToCSV(monthlyRevenueAndCashflow, folderPath+"/combine2.csv")
	fmt.Println("Combined monthly revenue and cashflow")
	return nil
}

func main() {
	// out := combineTwoCSV("output_stock/1101/per.csv", "output_stock/1101/stockdata.csv")
	// // out2 := combineTwoCSV("output_stock/1101/cashflow.csv", "output_stock/1101/monthlyrevenue.csv")
	// storage.SaveToCSV(out, "output_stock/1101/combine.csv")
	combineAllCSVInFolder("output_stock/2330")
}
