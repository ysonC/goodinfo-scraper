# Stock Scraper

A Go-based web scraper that collects various financial data—including weekly stock prices, PER valuations, monthly revenue, and cash flow statements—from [goodinfo.tw](https://goodinfo.tw).

## Key Features

- **Concurrent Scraping**: Configurable concurrency with interactive prompts (default: 10 workers).
- **Date Range Selection**: Supports custom or predefined date ranges.
- **Data Types**: Scrapes PER, stock data, monthly revenue, and cash flow.
- **Dynamic Web Scraping**: Uses headless browsers via [Playwright-Go](https://github.com/playwright-community/playwright-go).
- **Structured Output**:  
  - **CSV Files**: Saved per stock in the `downloaded_stock/` directory (one subfolder per stock).  
  - **XLSX Files**: Final combined spreadsheets stored in the `final_output/` directory (one per stock).

## Project Structure

```
.
├── Dockerfile
├── README.md
├── cmd
│   └── scraper
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── helper
│   │   ├── helper.go
│   │   ├── input.go
│   │   ├── setup.go
│   │   └── user_input.go
│   ├── scraper
│   │   ├── base.go
│   │   ├── cashflow.go
│   │   ├── factory.go
│   │   ├── per.go
│   │   ├── sale.go
│   │   ├── scrape.go
│   │   ├── scraper.go
│   │   └── stockdata.go
│   └── storage
│       └── csv_writer.go
```

- **`input_stock/`**: Place files here that contain stock numbers (one per line).
- **`downloaded_stock/`**: CSV output is saved here for each stock in its own subfolder.
- **`final_output/`**: Final XLSX output for each stock is generated here by combining CSV files.

## Local Installation

### Requirements

- **Go 1.24+** ([Download](https://go.dev/))
- **Playwright Dependencies**:  
  Install the Playwright driver by running:
  
  ```bash
  go run github.com/playwright-community/playwright-go/cmd/playwright install
  ```

### Setup

Clone the repository and download dependencies:

```bash
git clone https://github.com/your-username/stock-scraper.git
cd stock-scraper
go mod download
```

## Running Locally

Run the scraper using:

```bash
go run cmd/scraper/main.go
```

You will be prompted to:
- Select the maximum number of workers.
- Specify the date range (either default or custom).

After scraping completes:
- **CSV files** for each stock are saved under `downloaded_stock/` (each stock has its own subfolder).
- A **final XLSX file** is generated per stock in `final_output/`, combining the CSV data into multiple sheets.

## Docker Usage

### Building the Docker Image

```bash
docker build -t my-scraper .
```

### Running the Docker Container

```bash
docker run -it --rm \
  -v "$(pwd)/downloaded_stock:/app/downloaded_stock" \
  -v "$(pwd)/final_output:/app/final_output" \
  -v "$(pwd)/input_stock:/app/input_stock" \
  my-scraper
```

Adjust the volume mounts if your local directories differ.

## Customizing Concurrency

- The number of concurrent workers is configurable via the interactive prompt when running the scraper. The default is set to 10 workers.
- The website may block your IP if you set the number of workers too high. If you encounter issues, reduce the number of workers.
- Tested with up to 100 words without issues.

## Additional Information

- **CSV Combination**: The application verifies that required files (e.g., files containing "per", "stockdata", "monthlyrevenue", "cashflow") exist in each stock's output folder before combining them into the final XLSX.
- **XLSX Output**: The combined XLSX file features multiple sheets (e.g., one sheet for PER/stock data and another for monthly revenue/cash flow). This format is designed to provide a clear overview of all scraped data for each stock.
