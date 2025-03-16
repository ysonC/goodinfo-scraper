package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// PERScraper implements the scraper for PER data.
type PERScraper struct {
	base *BaseScraper
}

// NewPERScraper returns a new PERScraper.
func NewPERScraper(pw *playwright.Playwright) *PERScraper {
	base := NewBaseScraper(pw)
	return &PERScraper{base: base}
}

// Scrape fetches the PER data by building the URL and parsing the table HTML.
func (p *PERScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowK_ChartFlow.asp?RPT_CAT=PER&STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=F&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)
	html, err := p.base.fetchHTML(url)
	if err != nil {
		return nil, err
	}
	// For PER data, extract only the first 6 columns and skip the header row.
	return p.base.extractTableData(html, 6, true)
}
