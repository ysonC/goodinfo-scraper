package scraper

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/playwright-community/playwright-go"

	"github.com/ysonC/multi-stocks-download/internal/storage"
)

func ScrapeAllStocks(
	pw *playwright.Playwright,
	stocks, scraperTypes []string,
	startDate, endDate string,
	maxWorkers int,
	downloadDir string,
) ([]string, []string) {
	var (
		wg           sync.WaitGroup
		mutex        sync.Mutex
		successCount = make(map[string]int)
		totalTypes   = len(scraperTypes)
	)

	for _, stock := range stocks {
		successCount[stock] = 0
	}

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

				instance, err := NewScraper(scraperType, pw)
				if err != nil {
					log.Printf("Scraper creation error (%s) %s: %v", scraperType, stockNumber, err)
					return
				}

				data, err := instance.Scrape(stockNumber, startDate, endDate)
				if err != nil {
					log.Printf("Scraping error (%s) %s: %v", scraperType, stockNumber, err)
					return
				}

				if err := storage.WriteCSV(outputFile, data); err != nil {
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

	// Return successful and error stocks
	return checkDownloadStocks(successCount, totalTypes)
}

func checkDownloadStocks(successCount map[string]int, totalTypes int) ([]string, []string) {
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
	if len(errorStocks) > 0 {
		log.Println("Some tasks failed. Please check the logs for more information.")
	}
	return successfulStocks, errorStocks
}
