# Stock Scraper

This repository provides a Go-based scraper to gather weekly stock price and valuation data from [goodinfo.tw](https://goodinfo.tw). The scraper reads a list of stock numbers from CSV files in an `input_stock` folder, fetches the data, and saves the results as CSV files into an `output_stock` folder.

## Features

- **Concurrent Downloads**  
  Utilizes a worker pool (`maxWorkers = 10` by default) to download data for multiple stocks in parallel.

- **Up-to-Date Checks**  
  Skips downloading data if the output CSV for a stock is already current (i.e., generated or modified today).

- **Automatic CSV Output**  
  Each stock’s data is written to a separate CSV file (e.g., `2330.csv`) in the `output_stock` folder.

- **Headless Browsing with Playwright**  
  Uses Playwright (through the [playwright-go](https://github.com/playwright-community/playwright-go) library) to load and parse dynamic web content.

## Directory Structure

```
.
├── .gitignore
├── go.mod
├── go.sum
├── input_stock/
│   └── all_stocks_number.txt
├── one_stock.txt
├── output_stock/
└── scraper.go
```

- **`input_stock/`**: Contains one or more files with stock numbers. Each line in a represents a single stock number.
- **`output_stock/`**: Populated with CSV files for each stock that is successfully scraped.
- **`go.mod` / `go.sum`**: Go modules and dependency tracking.
- **`scraper.go`**: Main Go application that reads stock numbers, fetches data from the target site, and writes results.

> **Note**: By default, the scraper reads every CSV file in `input_stock`. An example file is `all_stocks_number.txt`, which lists a large set of stock numbers, one per line.

## Prerequisites

1. **Go 1.24 or later**  
   Make sure you have [Go](https://go.dev/) installed and properly set up (`GOPATH`, etc.).
   
2. **Playwright dependencies**  
   The scraper uses Playwright for Go. On most systems, you need to install the necessary browser dependencies.  
   - You can see the list of dependencies in [Playwright’s official docs](https://playwright.dev/), but typically it includes `chromium` dependencies and OS libraries for Chrome/Chromium.
   - Installing [Node.js](https://nodejs.org/) (≥ v14) can help if you need to use the `npx playwright install` command manually. However, the code can handle browser installation if run in an environment that allows Playwright to download the drivers.

## Installation

1. **Clone the repository** (or download it):
   ```bash
   git clone https://github.com/your-username/stock-scraper.git
   cd stock-scraper
   ```
   
2. **Download Go modules**:
   ```bash
   go mod download
   ```
   This will fetch all required dependencies as specified in `go.mod` / `go.sum`.

## Usage

1. **Prepare your stock numbers**  
   - Place CSV files in the `input_stock` folder.  
   - Each file can contain one stock number per line, e.g.:
     ```
     2330
     2317
     3008
     ```
   - The script will scan all CSV files in `input_stock` and gather all listed stock numbers.

2. **Run the scraper**  
   ```bash
   go run scraper.go
   ```
   or build and then run:
   ```bash
   go build -o scraper .
   ./scraper
   ```
   As it runs, you will see logs indicating which stock numbers are being processed or skipped (if up-to-date).

3. **Check the output**  
   - CSV files for each processed stock will appear in the `output_stock` folder.
   - Filenames match the stock number (`<stock>.csv`).  
   - Each output CSV includes a header row and subsequent weekly data rows.

## Customizing the Concurrency

- In `scraper.go`, find `const maxWorkers = 10`.
- Modify the value (e.g., `maxWorkers = 5`) to adjust how many stocks are scraped concurrently.

## Troubleshooting

1. **Playwright Browser Issues**  
   - Ensure you have the required OS libraries to run headless Chromium.
   - Check that your environment allows Playwright to download and install browser drivers.

2. **File Not Found or Permission Errors**  
   - Make sure `input_stock` and `output_stock` exist (the code will create `output_stock` if it doesn’t).
   - Verify you have read/write permissions for these folders.

3. **Data Missing or Incorrect**  
   - Verify the target site (`goodinfo.tw`) is accessible.
   - Check that your stock numbers are valid and properly formatted (no extra spaces, etc.).

## License

This project is distributed under the MIT License—see the [LICENSE](LICENSE) file for more details (if included in the repository).

