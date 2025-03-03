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
├── Dockerfile
├── go.mod
├── go.sum
├── input_stock/
│   └── all_stocks_number.txt
├── output_stock/
├── scraper.go
└── ...
```

- **`input_stock/`**: Contains one or more files with stock numbers. Each line in a file represents a single stock number.
- **`output_stock/`**: Populated with CSV files for each stock that is successfully scraped.
- **`go.mod` / `go.sum`**: Go modules and dependency tracking.
- **`scraper.go`**: Main Go application that reads stock numbers, fetches data from the target site, and writes results.

> **Note**: By default, the scraper reads every CSV file in `input_stock`. An example file is `all_stocks_number.txt`, which lists a large set of stock numbers, one per line.

## Prerequisites (Local Development)

1. **Go 1.24 or later**  
   Make sure you have [Go](https://go.dev/) installed and properly set up (`GOPATH`, etc.).

2. **Playwright dependencies**  
   The scraper uses Playwright for Go. On most systems, you need to install the necessary browser dependencies.  
   - Check [Playwright’s official docs](https://playwright.dev/) for the list of OS packages needed for Chromium/Firefox/WebKit.
   - Installing [Node.js](https://nodejs.org/) can help if you need the `npx playwright install` command.

## Installation (Local)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/stock-scraper.git
   cd stock-scraper
   ```
2. **Download Go modules**:
   ```bash
   go mod download
   ```
   This fetches all required dependencies as specified in `go.mod` / `go.sum`.

## Usage (Local)

1. **Prepare your stock numbers**  
   - Place CSV files in the `input_stock` folder (one stock number per line).
2. **Run the scraper**  
   ```bash
   go run scraper.go
   ```
   or build and then run:
   ```bash
   go build -o scraper .
   ./scraper
   ```
   The scraper will prompt you to select a date range (e.g., nearest 5 years or a custom range).

3. **Check the output**  
   - CSV files appear in the `output_stock` folder.
   - Filenames match the stock number (`<stock>.csv`).  
   - Each CSV includes a header row and subsequent data rows.

## Running via Docker

If you prefer not to install Go or browser dependencies locally, you can build and run the scraper in a Docker container that includes all necessary dependencies and Playwright binaries.

1. **Build the Docker image**:
   ```bash
   docker build -t my-scraper .
   ```
   This uses the `Dockerfile` in the repository (a multi-stage build) to:
   - Compile the Go code.
   - Install the Playwright-Go driver.
   - Copy everything into a Playwright-enabled base image.

2. **Run the Docker container** (interactive mode is required for date-range input):
   ```bash
   docker run -it --rm \
     -v $(pwd)/output_stock:/app/output_stock \
     -v $(pwd)/input_stock:/app/input_stock \
     my-scraper
   ```
   - Run in interactive mode with mounted volume to access CSV files on your host.
   - When prompted, select `1` for “Nearest 5 years” or `2` for “Custom range,” and then provide any required dates.


## Customizing the Concurrency

In `scraper.go`, find:

```go
const maxWorkers = 10
```

Change `10` to whatever concurrency you prefer.

## Troubleshooting

1. **Playwright Browser Issues**  
   - If running locally, ensure you have the OS libraries to run headless Chromium.  
   - In Docker, the provided Dockerfile includes these dependencies.

2. **File Permission Errors**  
   - Ensure `input_stock` and `output_stock` exist and are writable. If mounting volumes in Docker, check your host’s directory permissions.

3. **Data Missing or Incorrect**  
   - Verify `goodinfo.tw` is accessible (sometimes it blocks requests or changes the layout).
   - Check that your stock numbers are valid and properly formatted.

