package main

import (
	"log"

	"github.com/ysonC/multi-stocks-download/internal/flow"
	"github.com/ysonC/multi-stocks-download/internal/scraper"
	"github.com/ysonC/multi-stocks-download/internal/storage"
)

const (
	inputDir       = "input_stock"
	downloadDir    = "downloaded_stock"
	finalOutputDir = "final_output"
)

func main() {
	log.Println("Starting scraper application...")

	flow.SetupDirectories(downloadDir, finalOutputDir)
	stocks := flow.GetStockNumbers(inputDir)
	maxWorkers := flow.PromptMaxWorkers()
	startDate, endDate := flow.PromptDateRange()

	pw := flow.SetupPlaywright()
	defer pw.Stop()

	scraperTypes := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}
	successStocks, errorStocks := scraper.ScrapeAllStocks(
		pw,
		stocks,
		scraperTypes,
		startDate,
		endDate,
		maxWorkers,
		downloadDir,
	)
	if len(errorStocks) > 0 {
		log.Println("Some tasks failed. Please check the logs for more information.")
	}

	err := storage.CombineSuccessfulStocks(successStocks, downloadDir, finalOutputDir)
	if err != nil {
		log.Fatalf("Error combining successful stocks: %v", err)
	}

	log.Println("All tasks completed successfully.")
}
