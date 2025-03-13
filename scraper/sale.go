package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// PERScraper implements the scraper for PER data.
type MonthlyRevenueScraper struct {
	base *BaseScraper
}

// NewPERScraper returns a new PERScraper.
func NewMonthlyRevenueScraper(pw *playwright.Playwright) *MonthlyRevenueScraper {
	base := NewBaseScraper(pw)
	return &MonthlyRevenueScraper{base: base}
}

// Scrape fetches the PER data by building the URL and parsing the table HTML.
func (p *MonthlyRevenueScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowSaleMonChart.asp?STEP=DATA&STOCK_ID=%s&PRICE_ADJ=T&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)
	html, err := p.base.fetchHTML(url)
	if err != nil {
		return nil, err
	}
	return extractFullTableData(html)
}
