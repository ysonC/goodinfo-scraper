package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

const (
	downloadDir = "output_stock"
	inputFile   = "input_stock/stock_numbers.txt"
	maxWorkers  = 5 // Controls concurrency level (e.g., 5 workers)
)

// Read stock numbers from a file
func readStockNumbersFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var stockNumbers []string
	scanner := csv.NewReader(file)
	for {
		record, err := scanner.Read()
		if err != nil {
			break
		}
		stockNumbers = append(stockNumbers, strings.TrimSpace(record[0]))
	}
	return stockNumbers, nil
}

// Download stock data (Worker function)
func downloadStockData(stockNumber string, pw *playwright.Playwright, wg *sync.WaitGroup) {
	defer wg.Done()

	// Launch a browser (headless)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-setuid-sandbox"},
	})
	if err != nil {
		log.Printf("[ERROR] Failed to launch browser for %s: %v", stockNumber, err)
		return
	}
	defer browser.Close()

	// Create new page
	page, err := browser.NewPage()
	if err != nil {
		log.Printf("[ERROR] Failed to open page for %s: %v", stockNumber, err)
		return
	}

	// Construct stock URL
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowK_ChartFlow.asp?RPT_CAT=PER&STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=F&START_DT=2001-03-28&END_DT=%s",
		stockNumber,
		time.Now().Format("2006-01-02"),
	)

	// Visit the URL
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		log.Printf("[ERROR] Error visiting URL for stock %s: %v", stockNumber, err)
		return
	}

	// Wait for the table to load
	tableLocator := page.Locator("#tblDetail")
	err = tableLocator.WaitFor()
	if err != nil {
		log.Printf("[ERROR] Table not found for stock %s", stockNumber)
		return
	}

	// Extract table HTML
	tableHTML, err := tableLocator.InnerHTML()
	if err != nil {
		log.Printf("[ERROR] Failed to extract table for stock %s", stockNumber)
		return
	}

	// Parse table data
	data := extractTableData(tableHTML)

	// Save to CSV
	outputFilePath := filepath.Join(downloadDir, stockNumber+".csv")
	saveToCSV(data, outputFilePath)

	log.Printf("[SUCCESS] Stock %s data saved.", stockNumber)
}

// Extract table data from HTML
func extractTableData(html string) [][]string {
	rows := strings.Split(html, "<tr>")
	var data [][]string

	for _, row := range rows {
		if strings.Contains(row, "<td>") {
			cells := strings.Split(row, "<td>")
			var rowData []string
			for _, cell := range cells {
				text := strings.TrimSpace(stripHTMLTags(cell))
				rowData = append(rowData, text)
			}
			data = append(data, rowData)
		}
	}

	return data
}

// Remove HTML tags from a string
func stripHTMLTags(input string) string {
	return strings.NewReplacer("<td>", "", "</td>", "", "<tr>", "", "</tr>", "").Replace(input)
}

// Save data to CSV
func saveToCSV(data [][]string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		if len(row) > 0 {
			_ = writer.Write(row)
		}
	}
}

// Worker pool to download stocks concurrently
func downloadStockDataConcurrently(stockNumbers []string) {
	var wg sync.WaitGroup

	// Start Playwright once for all workers
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	// Create worker pool
	semaphore := make(chan struct{}, maxWorkers)

	// Start downloading stocks concurrently
	for _, stockNumber := range stockNumbers {
		wg.Add(1)
		semaphore <- struct{}{} // Block if maxWorkers is reached

		go func(stock string) {
			defer func() { <-semaphore }() // Release worker slot
			downloadStockData(stock, pw, &wg)
		}(stockNumber)
	}

	wg.Wait()
	log.Println("[INFO] All stock downloads completed.")
}

func main() {
	log.Println("[INFO] Script execution started.")

	// Read stock numbers from file
	stockNumbers, err := readStockNumbersFromFile("all_stocks_number.txt")
	if err != nil {
		log.Fatalf("[ERROR] Failed to read stock numbers: %v", err)
	}

	// Run concurrent stock downloads
	downloadStockDataConcurrently(stockNumbers)

	log.Println("[INFO] Script execution finished.")
}
