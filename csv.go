package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
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

func main() {
	fmt.Println("Reading CSV files for 1101")
	readCSVFile("output_stock/1101/per.csv")
	readCSVFile("output_stock/1101/stockdata.csv")
	readCSVFile("output_stock/1101/monthlyrevenue.csv")
	readCSVFile("output_stock/1101/cashflow.csv")

	fmt.Println("Reading CSV files for 2330")
	readCSVFile("output_stock/2330/per.csv")
	readCSVFile("output_stock/2330/stockdata.csv")
	readCSVFile("output_stock/2330/monthlyrevenue.csv")
	readCSVFile("output_stock/2330/cashflow.csv")
}
