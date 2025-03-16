package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/playwright-community/playwright-go"

	"github.com/ysonC/multi-stocks-download/scraper"
	"github.com/ysonC/multi-stocks-download/storage"
)

func scrapeAllStocks(
	pw *playwright.Playwright,
	stocks, scraperTypes []string,
	startDate, endDate string,
	maxWorkers int,
) ([]string, []string) {
	var (
		wg           sync.WaitGroup
		mutex        sync.Mutex
		successCount = make(map[string]int)
		totalTypes   = len(scraperTypes)
	)

	sem := make(chan struct{}, maxWorkers)

	for _, stock := range stocks {
		for _, sType := range scraperTypes {
			wg.Add(1)
			sem <- struct{}{}
			go func(stockNumber, scraperType string) {
				defer wg.Done()
				defer func() { <-sem }()

				stockOutputDir := filepath.Join(downloadDir, stockNumber)
				os.MkdirAll(stockOutputDir, 0755)

				outputFile := filepath.Join(stockOutputDir, scraperType+".csv")
				if storage.IsFileUpToDate(outputFile) {
					log.Printf("%s (%s) up-to-date, skipped.", stockNumber, scraperType)
					mutex.Lock()
					successCount[stockNumber]++
					mutex.Unlock()
					return
				}

				instance, err := scraper.NewScraper(scraperType, pw)
				if err != nil {
					log.Printf("Scraper creation error (%s) %s: %v", scraperType, stockNumber, err)
					return
				}

				data, err := instance.Scrape(stockNumber, startDate, endDate)
				if err != nil {
					log.Printf("Scraping error (%s) %s: %v", scraperType, stockNumber, err)
					return
				}

				if err := storage.SaveToCSV(data, outputFile); err != nil {
					log.Printf("CSV save error (%s) %s: %v", scraperType, stockNumber, err)
					return
				}

				mutex.Lock()
				successCount[stockNumber]++
				mutex.Unlock()

				log.Printf("Successfully scraped %s: %s", scraperType, stockNumber)
			}(stock, sType)
		}
	}

	wg.Wait()

	var successfulStocks []string
	var errorStocks []string
	for stock, count := range successCount {
		if count == totalTypes {
			successfulStocks = append(successfulStocks, stock)
		} else {
			errorStocks = append(errorStocks, stock)
			log.Printf("Incomplete data for stock %s; skipping combine.", stock)
		}
	}
	return successfulStocks, errorStocks
}

func combineSuccessfulStocks(stocks []string) {
	for _, stock := range stocks {
		stockDir := filepath.Join(downloadDir, stock)
		finalOutput := filepath.Join(finalOutputDir, stock+".xlsx")
		if err := storage.CombineAllCSVInFolderToXLSX(stockDir, finalOutput); err != nil {
			log.Printf("Error combining stock %s: %v", stock, err)
			continue
		}
		log.Printf("Successfully combined data for stock %s", stock)
	}
}
