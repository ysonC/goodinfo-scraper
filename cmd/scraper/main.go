package main

import (
	"log"

	"github.com/ysonC/multi-stocks-download/internal/helper"
	"github.com/ysonC/multi-stocks-download/internal/scraper"
)

const (
	inputDir       = "input_stock"
	downloadDir    = "downloaded_stock"
	finalOutputDir = "final_output"
)

func main() {
	log.Println("Starting scraper application...")

	helper.SetupDirectories(downloadDir, finalOutputDir)
	stocks := helper.GetStockNumbers(inputDir)
	maxWorkers := helper.PromptMaxWorkers()
	startDate, endDate := helper.PromptDateRange()

	pw := helper.SetupPlaywright()
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

	scraper.CombineSuccessfulStocks(successStocks, downloadDir, finalOutputDir)

	log.Println("All tasks completed successfully.")
}
