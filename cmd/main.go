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
	fmt.Println("1. Max  years")
	fmt.Println("2. Custom range (enter start and end dates, format YYYY-MM-DD)")
	fmt.Print("Enter option: ")

	var input string
	fmt.Scanln(&input)
	switch input {
	case "1":
		start := "1965-01-01"
		end := time.Now().Format("2006-01-02")
		return start, end, nil
	case "2":
		var start, end string
		fmt.Print("Enter start date (YYYY-MM-DD): ")
		fmt.Scanln(&start)
		fmt.Print("Enter end date (YYYY-MM-DD): ")
		fmt.Scanln(&end)
		return start, end, nil
	default:
		return "", "", fmt.Errorf("invalid option")
	}
}

// selectScraperType lets the user choose which type of data to scrape.
func selectScraperType() (string, error) {
	fmt.Println("Select scraper type:")
	fmt.Println("1. PER")
	fmt.Println("2. Stock Data")
	fmt.Println("3. Monthly Revenue")
	fmt.Println("4. Cashflow")
	fmt.Print("Enter option: ")

	var input string
	fmt.Scanln(&input)
	switch input {
	case "1":
		return "per", nil
	case "2":
		return "stockdata", nil
	case "3":
		return "monthlyrevenue", nil
	case "4":
		return "cashflow", nil
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

	startDate, endDate, err := selectDateRange()
	if err != nil {
		log.Fatalf("Error selecting date range: %v", err)
	}

	scraperTypes := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}

	// Start Playwright once.
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Failed to start Playwright: %v", err)
	}
	defer pw.Stop()

	// Use a semaphore to limit concurrency.
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	for _, stock := range stocks {
		for _, sType := range scraperTypes {

			wg.Add(1)
			sem <- struct{}{}
			go func(stockNumber, scraperType string) {
				defer wg.Done()
				defer func() { <-sem }()

				// Create an output folder for each stock.
				stockOutputDir := filepath.Join(outputDir, stockNumber)
				if err := os.MkdirAll(stockOutputDir, 0755); err != nil {
					log.Printf(
						"Failed to create output directory for stock %s: %v",
						stockNumber,
						err,
					)
					return
				}

				// Build output filename including the scraper type.
				outputFile := filepath.Join(stockOutputDir, scraperType+".csv")
				if storage.IsFileUpToDate(outputFile) {
					log.Printf("Stock %s data is up-to-date. Skipping.", stockNumber)
					return
				}

				scraperInstance, err := scraper.NewScraper(scraperType, pw)
				if err != nil {
					log.Printf("Error creating scraper for stock %s: %v", stockNumber, err)
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
				log.Printf(
					"Successfully scraped and saved data for %s : %s",
					scraperType,
					stockNumber,
				)
			}(stock, sType)
		}
	}

	wg.Wait()
	log.Println("Scraping completed.")
}
