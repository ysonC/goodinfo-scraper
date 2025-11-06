package main

import (
	"log"
	"time"

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
	start := time.Now()

	flow.SetupDirectories(inputDir, downloadDir, finalOutputDir)
	stocks := flow.GetStockNumbers(inputDir)
	maxWorkers := flow.PromptMaxWorkers()
	startDate, endDate := flow.PromptDateRange()

	pw := flow.SetupPlaywright()
	defer pw.Stop()

	scraperTypes := []string{"per", "stockdata", "monthlyrevenue", "cashflow"}

	downloadStart := time.Now()
	successStocks, errorStocks := scraper.ScrapeAllStocks(
		pw,
		stocks,
		scraperTypes,
		startDate,
		endDate,
		maxWorkers,
		downloadDir,
	)
	log.Printf("Download process completed in %s", time.Since(downloadStart))

	if len(errorStocks) > 0 {
		log.Println("Some tasks failed. Please check the logs for more information.")
	}

	err := storage.CombineSuccessfulStocks(successStocks, downloadDir, finalOutputDir)
	if err != nil {
		log.Fatalf("Error combining successful stocks: %v", err)
	}

	log.Printf("All tasks completed successfully in %s.", time.Since(start))
}
