package scraper

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
)

// NewScraper returns a Scraper instance based on the given type.
func NewScraper(scraperType string, pw *playwright.Playwright) (Scraper, error) {
	switch scraperType {
	case "per":
		return NewPERScraper(pw), nil
	case "stockdata":
		return NewStockDataScraper(pw), nil
	default:
		return nil, fmt.Errorf("unknown scraper type: %s", scraperType)
	}
}
