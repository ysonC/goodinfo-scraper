package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"

	"github.com/ysonC/multi-stocks-download/scraper"
	"github.com/ysonC/multi-stocks-download/storage"
)

const (
	inputDir   = "input_stock"
	outputDir  = "output_stock"
	maxWorkers = 10
)

// selectDateRange lets the user choose a date range if needed.
func selectDateRange() (string, string, error) {
	fmt.Println("Select date range:")
	fmt.Println("1. Nearest 5 years (for PER scraping)")
	fmt.Println("2. Custom range (enter start and end dates, format YYYY-MM-DD)")
	fmt.Println("3. No date range (for scrapers that don’t require dates)")
	fmt.Print("Enter option: ")

	var input string
	fmt.Scanln(&input)
	switch input {
	case "1":
		start := "2020-03-14"
		end := time.Now().Format("2006-01-02")
		return start, end, nil
	case "2":
		var start, end string
		fmt.Print("Enter start date (YYYY-MM-DD): ")
		fmt.Scanln(&start)
		fmt.Print("Enter end date (YYYY-MM-DD): ")
		fmt.Scanln(&end)
		return start, end, nil
	case "3":
		return "", "", nil
	default:
		return "", "", fmt.Errorf("invalid option")
	}
}

// selectScraperType lets the user choose which type of data to scrape.
func selectScraperType() (string, error) {
	fmt.Println("Select scraper type:")
	fmt.Println("1. PER")
	fmt.Println("2. Cash Flow")
	fmt.Println("3. Monthly Revenue")
	fmt.Print("Enter option: ")

	var input string
	fmt.Scanln(&input)
	switch input {
	case "1":
		return "per", nil
	case "2":
		return "cashflow", nil
	case "3":
		return "revenue", nil
	default:
		return "", fmt.Errorf("invalid option")
	}
}

// readStockNumbersFromFolder reads stock numbers (one per line) from files in the input folder.
func readStockNumbersFromFolder(folderPath string) ([]string, error) {
	var stocks []string
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(folderPath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file %s: %v", path, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				// In case the file is CSV formatted, split on comma.
				parts := strings.Split(line, ",")
				stocks = append(stocks, strings.TrimSpace(parts[0]))
			}
		}
		f.Close()
	}
	return stocks, nil
}

func main() {
	log.Println("Starting scraper application...")

	// Ensure output directory exists.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	stocks, err := readStockNumbersFromFolder(inputDir)
	if err != nil {
		log.Fatalf("Failed to read stock numbers: %v", err)
	}
	if len(stocks) == 0 {
		log.Fatalf("No stock numbers found in %s", inputDir)
	}

	scraperType, err := selectScraperType()
	if err != nil {
		log.Fatalf("Error selecting scraper type: %v", err)
	}

	startDate, endDate, err := selectDateRange()
	if err != nil {
		log.Fatalf("Error selecting date range: %v", err)
	}

	// Start Playwright once.
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	// Create scraper instance based on the chosen type.
	var scraperInstance interface {
		Scrape(stockNumber, startDate, endDate string) ([][]string, error)
	}
	switch scraperType {
	case "per":
		scraperInstance = scraper.NewPERScraper(pw)
	case "cashflow":
		scraperInstance = scraper.NewCashFlowScraper(pw)
	default:
		log.Fatalf("Unknown scraper type: %s", scraperType)
	}

	// Use a semaphore to limit concurrency.
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for _, stock := range stocks {
		wg.Add(1)
		sem <- struct{}{}
		go func(stockNumber string) {
			defer wg.Done()
			defer func() { <-sem }()

			// Build output filename including the scraper type.
			outputFile := filepath.Join(outputDir, stockNumber+"_"+scraperType+".csv")
			if storage.IsFileUpToDate(outputFile) {
				log.Printf("Stock %s data is up-to-date. Skipping.", stockNumber)
				return
			}

			data, err := scraperInstance.Scrape(stockNumber, startDate, endDate)
			if err != nil {
				log.Printf("Error scraping stock %s: %v", stockNumber, err)
				return
			}

			err = storage.SaveToCSV(data, outputFile)
			if err != nil {
				log.Printf("Error saving CSV for stock %s: %v", stockNumber, err)
				return
			}
			log.Printf("Successfully scraped and saved data for stock %s", stockNumber)
		}(stock)
	}

	wg.Wait()
	log.Println("Scraping completed.")
}
