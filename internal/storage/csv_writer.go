package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CombineSuccessfulStocks(stocks []string, downloadDir, finalOutputDir string) error {
	for _, stock := range stocks {
		stockDir := filepath.Join(downloadDir, stock)
		finalOutput := filepath.Join(finalOutputDir, stock+".csv")
		if err := combineAllCSVInFolder(stockDir, finalOutput); err != nil {
			log.Printf("Error combining stock %s: %v", stock, err)
			continue
		}
		log.Printf("Successfully combined data for stock %s", stock)
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

func combineAllCSVInFolder(folderPath, finalOutput string) error {
	files, err := ReadDirFiles(folderPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	err = CheckFileExist(files)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %v", err)
	}

	perData, err := ReadCSV(folderPath + "/per.csv")
	if err != nil {
		return fmt.Errorf("failed to read per.csv: %v", err)
	}
	stockData, err := ReadCSV(folderPath + "/stockdata.csv")
	if err != nil {
		return fmt.Errorf("failed to read stockdata.csv: %v", err)
	}
	monthlyRevenueData, err := ReadCSV(folderPath + "/monthlyrevenue.csv")
	if err != nil {
		return fmt.Errorf("failed to read monthlyrevenue.csv: %v", err)
	}
	cashflowData, err := ReadCSV(folderPath + "/cashflow.csv")
	if err != nil {
		return fmt.Errorf("failed to read cashflow.csv: %v", err)
	}
	equityData, err := ReadCSV(folderPath + "/equity.csv")
	if err != nil {
		return fmt.Errorf("failed to read equity.csv: %v", err)
	}

	// Add headers to the data
	perData, err = addPERHeaderNew(perData)
	if err != nil {
		return fmt.Errorf("failed to add header to per data: %v", err)
	}
	stockData, err = addStockDataHeader(stockData)
	if err != nil {
		return fmt.Errorf("failed to add header to stock data: %v", err)
	}
	monthlyRevenueData, err = addMonthlyRevenueHeader(monthlyRevenueData)
	if err != nil {
		return fmt.Errorf("failed to add header to monthly revenue data: %v", err)
	}
	cashflowData, err = addCashflowHeader(cashflowData)
	if err != nil {
		return fmt.Errorf("failed to add header to cashflow data: %v", err)
	}
	equityData, err = addEquityHeader(equityData)
	if err != nil {
		return fmt.Errorf("failed to add header to equity data: %v", err)
	}

	mergedData, err := mergeCSVData(perData, stockData)
	if err != nil {
		return fmt.Errorf("failed to merge per and stock data: %v", err)
	}
	mergedData, err = mergeCSVData(mergedData, monthlyRevenueData)
	if err != nil {
		return fmt.Errorf("failed to merge monthly revenue data: %v", err)
	}
	mergedData, err = mergeCSVData(mergedData, cashflowData)
	if err != nil {
		return fmt.Errorf("failed to merge cashflow data: %v", err)
	}
	mergedData, err = mergeCSVData(mergedData, equityData)
	if err != nil {
		return fmt.Errorf("failed to merge equity data: %v", err)
	}

	// Write the final finalOutput
	err = WriteCSV(finalOutput, mergedData)
	if err != nil {
		return fmt.Errorf("failed to write final output: %v", err)
	}

	return nil
}

func mergeCSVData(csv1, csv2 [][]string) ([][]string, error) {
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

func addPERHeaderNew(data [][]string) ([][]string, error) {
	header := [][]string{
		{
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{"交易週別", "收盤價", "漲跌價", "漲跌幅", "河流圖 EPS(元)", "目前 PER (倍)"},
	}

	return append(header, data...), nil
}

func addStockDataHeader(data [][]string) ([][]string, error) {
	header := [][]string{
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"成交張數",
			"",
			"成交金額",
			"",
			"法人買賣超(千張)",
			"",
			"",
			"",
			"",
			"融資(千張)",
			"",
			"融券(千張)",
			"",
			"",
		},
		{
			"交易週別",
			"交易日數",
			"開盤",
			"最高",
			"最低",
			"收盤",
			"漲跌",
			"漲跌(%)",
			"振幅(%)",
			"千張",
			"日均",
			"億元",
			"日均",
			"外資",
			"投信",
			"自營",
			"合計",
			"外資持股(%)",
			"增減",
			"餘額",
			"增減",
			"餘額",
			"券資比(%)",
		},
	}

	return append(header, data...), nil
}

func addMonthlyRevenueHeader(data [][]string) ([][]string, error) {
	header := [][]string{
		{
			"",
			"當月股價",
			"",
			"",
			"",
			"",
			"",
			"營業收入",
			"",
			"",
			"",
			"",
			"合併營業收入",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"單月",
			"",
			"",
			"累計",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"月別",
			"開盤",
			"收盤",
			"最高",
			"最低",
			"漲跌(元)",
			"漲跌(%)",
			"營收(億)",
			"月增(%)",
			"年增(%)",
			"營收(億)",
			"年增(%)",
			"營收(億)",
			"月增(%)",
			"年增(%)",
			"營收(億)",
			"年增(%)",
		},
	}
	return append(header, data...), nil
}

func addCashflowHeader(data [][]string) ([][]string, error) {
	header := [][]string{
		{
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"",
			"",
			"",
			"季度股價",
			"",
			"",
			"",
			"獲利(億)",
			"",
			"現金流量(億)",
			"",
			"",
			"",
			"",
			"",
			"現金餘額(億)",
			"",
			"",
			"",
		},
		{
			"季度",
			"平均股本(億)",
			"財報評分",
			"上期收盤",
			"本期收盤",
			"漲跌(元)",
			"漲跌(%)",
			"稅前淨利",
			"稅後淨利",
			"營業活動",
			"投資活動",
			"融資活動",
			"其他活動",
			"淨現金流",
			"自由金流",
			"期初餘額",
			"期末餘額",
			"現金流量(%)",
			"稅後EPS(元)",
		},
	}
	return append(header, data...), nil
}

func addEquityHeader(data [][]string) ([][]string, error) {
	header := [][]string{
		{
			"",
			"",
			"",
			"當週股價",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
		{
			"週別",
			"統計日期",
			"收盤",
			"漲跌(元)",
			"漲跌(%)",
			"集保庫存(萬張)",
			"≦10張",
			">10張≦50張",
			">50張≦100張",
			">100張≦200張",
			">200張≦400張",
			">400張≦800張",
			">800張≦1000張",
			">1000張",
		},
	}
	return append(header, data...), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
