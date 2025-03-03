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

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

const (
	downloadDir = "output_stock"
	inputDir    = "input_stock"
	inputFile   = "input_stock/stock_numbers.txt"
	maxWorkers  = 10 // Controls concurrency level
)

// Reads stock numbers from every CSV file in the specified folder.
func readStockNumbersFromFolder(folderPath string) ([]string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", folderPath, err)
	}

	var stockNumbers []string
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip subdirectories.
			continue
		}
		filePath := filepath.Join(folderPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			// Log error and continue with the next file.
			fmt.Printf("Error opening file %s: %v\n", filePath, err)
			continue
		}

		reader := csv.NewReader(file)
		for {
			record, err := reader.Read()
			if err != nil {
				break // EOF or error; proceed to next file.
			}
			if len(record) > 0 {
				stockNumbers = append(stockNumbers, strings.TrimSpace(record[0]))
			}
		}
		file.Close()
	}

	return stockNumbers, nil
}

// isStockDataUpToDate checks if the CSV file for the given stock exists and if its modification time is today.
func isStockDataUpToDate(stockNumber string) bool {
	filePath := filepath.Join(downloadDir, stockNumber+".csv")
	fi, err := os.Stat(filePath)
	if err != nil {
		// File doesn't exist.
		return false
	}
	now := time.Now()
	modTime := fi.ModTime()
	// Check if file was modified on the same day.
	return now.Year() == modTime.Year() && now.YearDay() == modTime.YearDay()
}

// downloadStockData downloads and saves data for one stock.
func downloadStockData(stockNumber string, pw *playwright.Playwright, wg *sync.WaitGroup) {
	defer wg.Done()
	start_date := "2025-02-28"
	end_date := time.Now().Format("2006-01-02")

	// Launch a headless browser.
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-setuid-sandbox"},
	})
	if err != nil {
		log.Printf("[ERROR] Failed to launch browser for %s: %v", stockNumber, err)
		return
	}
	defer browser.Close()

	// Create a new page.
	page, err := browser.NewPage()
	if err != nil {
		log.Printf("[ERROR] Failed to open page for %s: %v", stockNumber, err)
		return
	}

	// Construct stock URL.
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowK_ChartFlow.asp?RPT_CAT=PER&STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=F&START_DT=%s&END_DT=%s",
		stockNumber,
		start_date,
		end_date,
	)

	// Visit the URL.
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		log.Printf("[ERROR] Error visiting URL for stock %s: %v", stockNumber, err)
		return
	}

	// Wait for the table container to load.
	tableLocator := page.Locator("#tblDetail")
	err = tableLocator.WaitFor()
	if err != nil {
		log.Printf("[ERROR] Table not found for stock %s", stockNumber)
		// Even if not found, continue to attempt extraction.
	}

	// Extract the inner HTML of the table container.
	tableHTML, err := tableLocator.InnerHTML()
	if err != nil {
		log.Printf("[ERROR] Failed to extract table for stock %s: %v", stockNumber, err)
		return
	}

	// Parse table data from the HTML.
	data, err := extractTableData(tableHTML)
	if err != nil {
		log.Printf("[ERROR] Failed to extract table data for stock %s: %v", stockNumber, err)
		return
	}

	// Save to CSV with header.
	outputFilePath := filepath.Join(downloadDir, stockNumber+".csv")
	saveToCSV(data, outputFilePath)
	log.Printf("[SUCCESS] Stock %s data saved.", stockNumber)
}

// extractTableData parses the provided HTML (assumed to be a table) and returns a 2D slice.
// It wraps the fragment in a <table> so that <tr> elements are properly parsed.
func extractTableData(html string) ([][]string, error) {
	var data [][]string

	// Wrap the HTML fragment.
	wrappedHTML := "<table>" + html + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}

	// Find all rows.
	rows := doc.Find("tr")
	if rows.Length() == 0 {
		return nil, fmt.Errorf("no <tr> elements found in the table HTML")
	}

	// Iterate over rows, skipping the header row.
	rows.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return // Skip header.
		}
		var rowData []string
		s.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			rowData = append(rowData, text)
		})
		if len(rowData) > 0 {
			// Skip rows whose first cell (Date) ends with "W53"
			if strings.HasSuffix(rowData[0], "W53") {
				return
			}
			data = append(data, rowData)
		}
	})

	return data, nil
}

// saveToCSV writes the CSV file with a fixed header.
func saveToCSV(data [][]string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Date", "Price", "Change", "% Change", "EPS", "PER",
		"8X", "9.8X", "11.6X", "13.4X", "15.2X", "17X",
	}
	if err := writer.Write(header); err != nil {
		log.Fatalf("[ERROR] Failed to write header: %v", err)
	}

	for _, row := range data {
		if len(row) > 0 {
			if err := writer.Write(row); err != nil {
				log.Printf("[ERROR] Failed to write row: %v", err)
			}
		}
	}
}

// downloadStockDataConcurrently processes stocks concurrently, skipping stocks whose data is up to date.
func downloadStockDataConcurrently(stockNumbers []string) {
	var wg sync.WaitGroup

	// Start Playwright once.
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	semaphore := make(chan struct{}, maxWorkers)
	for _, stockNumber := range stockNumbers {
		if isStockDataUpToDate(stockNumber) {
			log.Printf("[INFO] Stock %s is up-to-date; skipping download.", stockNumber)
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}
		go func(stock string) {
			defer func() { <-semaphore }()
			downloadStockData(stock, pw, &wg)
		}(stockNumber)
	}

	wg.Wait()
	log.Println("[INFO] All stock downloads completed.")
}

func main() {
	log.Println("[INFO] Script execution started.")

	// Read stock numbers from file.
	stockNumbers, err := readStockNumbersFromFolder(inputDir)
	// stockNumbers = []string{"2330"}
	if err != nil {
		log.Fatalf("[ERROR] Failed to read stock numbers: %v", err)
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Fatalf("[ERROR] Failed to create download directory: %v", err)
	}

	downloadStockDataConcurrently(stockNumbers)

	log.Println("[INFO] Script execution finished.")
}
