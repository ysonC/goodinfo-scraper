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
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) > 0 {
			stockNumbers = append(stockNumbers, strings.TrimSpace(record[0]))
		}
	}
	return stockNumbers, nil
}

// Download stock data (Worker function)
func downloadStockData(stockNumber string, pw *playwright.Playwright, wg *sync.WaitGroup) {
	defer wg.Done()

	// Launch a headless browser
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

	// Wait for the table container to load.
	tableLocator := page.Locator("#tblDetail")
	err = tableLocator.WaitFor()
	if err != nil {
		log.Printf("[ERROR] Table not found for stock %s", stockNumber)
		// Even if the locator isn't found, we try to continue.
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
// It skips the first (header) row and any row with a Date (first cell) ending with "W53".
func extractTableData(html string) ([][]string, error) {
	var data [][]string

	// Create a goquery document from the HTML string.
	wrappedHTML := "<table>" + html + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}

	fmt.Println(html)
	// Find all rows in the table.
	rows := doc.Find("tr")
	if rows.Length() == 0 {
		return nil, fmt.Errorf("no <tr> elements found in the table HTML")
	}

	// Iterate over rows, skipping the header row (assumed to be the first row).
	rows.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// Skip header row.
			return
		}
		var rowData []string
		// Extract text from each <td> cell.
		s.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			rowData = append(rowData, text)
		})
		// Only add non-empty rows.
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

// saveToCSV writes a header and data rows into a CSV file.
func saveToCSV(data [][]string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row matching the Python output.
	header := []string{
		"Date",
		"Price",
		"Change",
		"% Change",
		"EPS",
		"PER",
		"8X",
		"9.8X",
		"11.6X",
		"13.4X",
		"15.2X",
		"17X",
	}
	if err := writer.Write(header); err != nil {
		log.Fatalf("[ERROR] Failed to write header: %v", err)
	}

	// Write each data row.
	for _, row := range data {
		if len(row) > 0 {
			if err := writer.Write(row); err != nil {
				log.Printf("[ERROR] Failed to write row: %v", err)
			}
		}
	}
}

// downloadStockDataConcurrently runs the stock downloads concurrently.
func downloadStockDataConcurrently(stockNumbers []string) {
	var wg sync.WaitGroup

	// Start Playwright once for all workers.
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	// Create a worker pool.
	semaphore := make(chan struct{}, maxWorkers)

	for _, stockNumber := range stockNumbers {
		wg.Add(1)
		semaphore <- struct{}{} // Block if maxWorkers is reached.

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
	stockNumbers, err := readStockNumbersFromFile("one_stock.txt")
	if err != nil {
		log.Fatalf("[ERROR] Failed to read stock numbers: %v", err)
	}

	// Run concurrent stock downloads.
	downloadStockDataConcurrently(stockNumbers)

	log.Println("[INFO] Script execution finished.")
}
