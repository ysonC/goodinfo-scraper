# Stock Scraper

A Go-based web scraper that collects various financial data, including weekly stock prices, PER valuations, monthly revenue, and cash flow statements from [goodinfo.tw](https://goodinfo.tw).

## Key Features

- **Concurrent Scraping**: Configurable concurrency with interactive prompts (default: 10 workers).
- **Date Range Selection**: Supports custom or predefined date ranges.
- **Data Types**: Includes PER, stock data, monthly revenue, and cash flow.
- **Dynamic Web Scraping**: Utilizes headless browsers via [Playwright-Go](https://github.com/playwright-community/playwright-go).
- **Structured CSV Output**: Results are organized into individual CSV files per stock in structured directories.

## Project Structure

```
.
├── .gitignore
├── Dockerfile
├── README.md
├── cmd
│   └── main.go
├── criteria.md
├── go.mod
├── go.sum
├── sample
│   └── final-output.xlsx
├── scraper
│   ├── base.go
│   ├── cashflow.go
│   ├── factory.go
│   ├── per.go
│   ├── sale.go
│   ├── scraper.go
│   └── stockdata.go
└── storage
    └── csv_writer.go
```

- **`input_stock/`**: Place files containing stock numbers (one per line).
- **`output_stock/`**: Automatically generated CSV output.

## Local Installation

### Requirements

- **Go 1.24+** ([Download](https://go.dev/))
- **Playwright dependencies**:
  ```bash
  go run github.com/playwright-community/playwright-go/cmd/playwright install
  ```

### Setup

```bash
git clone https://github.com/your-username/stock-scraper.git
cd stock-scraper
go mod download
```

## Running Locally

```bash
go run cmd/main.go
```

Follow the interactive prompts to:
- Select the maximum number of workers.
- Choose data types.
- Specify date ranges.

## Docker Usage

### Building Docker Image

```bash
docker build -t my-scraper .
```

### Running Docker Container

```bash
docker run -it --rm \
  -v $(pwd)/output_stock:/app/output_stock \
  -v $(pwd)/input_stock:/app/input_stock \
  my-scraper
```

## Customizing Concurrency

Modify the number of concurrent workers through the interactive prompt upon running the scraper, or adjust defaults directly in `cmd/main.go`.

---

