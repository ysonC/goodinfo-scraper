package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ysonC/multi-stocks-download/internal/flow"
	"github.com/ysonC/multi-stocks-download/internal/scraper"
	"github.com/ysonC/multi-stocks-download/internal/storage"
)

const dataDir = "data"

var (
	inputDir       = filepath.Join(dataDir, "input_stock")
	downloadDir    = filepath.Join(dataDir, "downloaded_stock")
	finalOutputDir = filepath.Join(dataDir, "final_output")
	failedDir      = filepath.Join(dataDir, "failed_stock")
)

func main() {
	log.Println("Starting scraper application...")
	start := time.Now()

	maxWorkersFlag := flag.Int("workers", 10, "maximum number of concurrent workers (default 10)")
	flag.IntVar(maxWorkersFlag, "w", 10, "shorthand for -workers")
	rerunFailedFlag := flag.Bool(
		"rerun-failed",
		false,
		"rerun only failed stocks from the previous run",
	)
	flag.BoolVar(rerunFailedFlag, "rf", false, "shorthand for -rerun-failed")

	startDateFlag := flag.String(
		"start",
		"",
		"start date in YYYY-MM-DD (default 1965-01-01 when omitted together with -end)",
	)
	endDateFlag := flag.String(
		"end",
		"",
		"end date in YYYY-MM-DD (default today when omitted together with -start)",
	)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), `
Examples:
  # Default: 10 workers, full date range
  scraper

  # Custom workers and date range
  scraper -workers=20 -start=2020-01-01 -end=2024-12-31`)
	}

	flag.Parse()

	maxWorkers := *maxWorkersFlag
	rerunFailed := *rerunFailedFlag
	if maxWorkers <= 0 {
		log.Fatalf("Invalid workers value %d, must be > 0", maxWorkers)
	}

	var startDate, endDate string
	if *startDateFlag == "" && *endDateFlag == "" {
		startDate = "1965-01-01"
		endDate = time.Now().Format("2006-01-02")
	} else if *startDateFlag != "" && *endDateFlag != "" {
		startDate = *startDateFlag
		endDate = *endDateFlag
	} else {
		log.Fatal("Both -start and -end must be provided together, or neither for max range")
	}

	flow.SetupDirectories(inputDir, downloadDir, finalOutputDir, failedDir)

	var stocks []string
	if rerunFailed {
		var err error
		stocks, err = storage.LoadFailedStocks(failedDir)
		if err != nil {
			log.Fatalf("Failed to load failed stocks: %v", err)
		}
		if len(stocks) == 0 {
			log.Println("No recorded failed stocks to rerun. Exiting.")
			return
		}
		log.Printf("Rerunning %d previously failed stock(s).", len(stocks))
	} else {
		stocks = flow.GetStockNumbers(inputDir)
	}

	pw := flow.SetupPlaywright()
	defer pw.Stop()

	scraperTypes := []string{"per", "stockdata", "monthlyrevenue", "cashflow", "equity"}

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

	successCount := len(successStocks)
	errorCount := len(errorStocks)

	if err := storage.SaveFailedStocks(failedDir, errorStocks); err != nil {
		log.Fatalf("Failed to write failed stocks list: %v", err)
	}

	err := storage.CombineSuccessfulStocks(successStocks, downloadDir, finalOutputDir)
	if err != nil {
		log.Fatalf("Error combining successful stocks: %v", err)
	}

	failedSummary := ""
	if errorCount > 0 {
		failedSummary = " Failed stocks: " + strings.Join(errorStocks, ", ") + "."
	}

	log.Printf(
		"Run summary: %d stocks succeeded, %d failed.%s Total duration: %s.",
		successCount,
		errorCount,
		failedSummary,
		time.Since(start),
	)
}
