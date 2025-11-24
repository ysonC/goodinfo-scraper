package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

type MonthlyRevenueScraper struct {
	base *BaseScraper
}

func NewMonthlyRevenueScraper(pw *playwright.Playwright) *MonthlyRevenueScraper {
	base := NewBaseScraper(pw)
	return &MonthlyRevenueScraper{base: base}
}

func (p *MonthlyRevenueScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowSaleMonChart.asp?STOCK_ID=%s&PRICE_ADJ=T&START_DT=%s&END_DT=%s",
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
