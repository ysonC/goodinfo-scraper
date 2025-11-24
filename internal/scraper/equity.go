package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type EquityScraper struct {
	base *BaseScraper
}

func NewEquityScraper(pw *playwright.Playwright) *EquityScraper {
	base := NewBaseScraper(pw)
	return &EquityScraper{base: base}
}

func (p *EquityScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/EquityDistributionClassHis.asp?STOCK_ID=%s&PRICE_ADJ=T&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)
	html, err := p.base.fetchHTML(url)
	if err != nil {
		return nil, err
	}
	return p.base.extractFullTableData(html)
}
