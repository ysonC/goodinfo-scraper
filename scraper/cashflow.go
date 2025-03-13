package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// PERScraper implements the scraper for PER data.
type CashflowScraper struct {
	base *BaseScraper
}

// NewPERScraper returns a new PERScraper.
func NewCashflowScraper(pw *playwright.Playwright) *CashflowScraper {
	base := NewBaseScraper(pw)
	return &CashflowScraper{base: base}
}

// Scrape fetches the PER data by building the URL and parsing the table HTML.
func (p *CashflowScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/StockCashFlow.asp?STEP=DATA&STOCK_ID=%s&RPT_CAT=M_QUAR&PRICE_ADJ=F&START_DT=%s&END_DT=%s",
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
