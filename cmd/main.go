package main

import (
	"log"
)

const (
	inputDir       = "input_stock"
	downloadDir    = "downloaded_stock"
	finalOutputDir = "final_output"
)

func main() {
	log.Println("Starting scraper application...")

	setupDirectories(downloadDir, finalOutputDir)
	stocks := getStockNumbers(inputDir)
	maxWorkers := promptMaxWorkers()
	startDate, endDate := promptDateRange()

	pw := setupPlaywright()
	defer pw.Stop()

	scraperTypes := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	successStocks, errorStocks := scrapeAllStocks(
		pw,
		stocks,
		scraperTypes,
		startDate,
		endDate,
		maxWorkers,
	)
	if len(errorStocks) > 0 {
		log.Println("Some tasks failed. Please check the logs for more information.")
	}

	combineSuccessfulStocks(successStocks)

	log.Println("All tasks completed successfully.")
}
