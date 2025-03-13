package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

// StockDataScraper implements the scraper for stock data.
type StockDataScraper struct {
	base *BaseScraper
}

// NewStockDataScraper returns a new StockDataScraper.
func NewStockDataScraper(pw *playwright.Playwright) *StockDataScraper {
	base := NewBaseScraper(pw)
	return &StockDataScraper{base: base}
}

// Scrape fetches the stock data by building the URL and parsing the table HTML.
func (p *StockDataScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowK_Chart.asp?STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=T&SHEET=%%E5%%80%%8B%%E8%%82%%A1%%E8%%82%%A1%%E5%%83%%B9%%E3%%80%%81%%E6%%B3%%95%%E4%%BA%%BA%%E8%%B2%%B7%%E8%%B3%%A3%%E5%%8F%%8A%%E8%%9E%%8D%%E8%%B3%%87%%E5%%88%%B8&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)

	html, err := p.base.fetchHTML(url)
	if err != nil {
		return nil, err
	}
	// For stock data, extract all columns including header.
	return extractFullTableData(html)
}

// extractFullTableData parses the table HTML without skipping the header.
func extractFullTableData(html string) ([][]string, error) {
	var data [][]string
	wrappedHTML := "<table>" + html + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		var row []string
		s.Find("td").Each(func(j int, cell *goquery.Selection) {
			row = append(row, strings.TrimSpace(cell.Text()))
		})
		if len(row) > 0 {
			data = append(data, row)
		}
	})
	return data, nil
}
