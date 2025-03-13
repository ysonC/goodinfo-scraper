package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

// PERScraper implements the scraper for PER data.
type StockDataScraper struct {
	pw *playwright.Playwright
}

// NewPERScraper returns a new PERScraper.
func NewStockDataScraper(pw *playwright.Playwright) *StockDataScraper {
	return &StockDataScraper{pw: pw}
}

// Scrape navigates to the PER URL, extracts the table HTML, and parses it.
func (p *StockDataScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
	browser, err := p.pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-setuid-sandbox"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	url := fmt.Sprintf(
		"https://goodinfo.tw/tw/ShowK_Chart.asp?STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=T&SHEET=%%E5%%80%%8B%%E8%%82%%A1%%E8%%82%%A1%%E5%%83%%B9%%E3%%80%%81%%E6%%B3%%95%%E4%%BA%%BA%%E8%%B2%%B7%%E8%%B3%%A3%%E5%%8F%%8A%%E8%%9E%%8D%%E8%%B3%%87%%E5%%88%%B8&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)

	log.Printf("Scraping URL: %s", url)
	if _, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return nil, fmt.Errorf("failed to goto URL: %w", err)
	}

	tableLocator := page.Locator("#tblDetail")
	if err := tableLocator.WaitFor(); err != nil {
		log.Printf("Warning: table not found for stock %s", stockNumber)
	}

	// Wait for the table's inner text to be non-empty.
	if _, err := page.WaitForFunction(`() => {
    return document.querySelector("#tblDetail").querySelectorAll("tr").length >= 200;
}`, nil); err != nil {
		return nil, fmt.Errorf("timeout waiting for table data: %w", err)
	}

	tableHTML, err := tableLocator.InnerHTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get table HTML: %w", err)
	}

	// Parse the HTML using goquery.
	var data [][]string
	wrappedHTML := "<table>" + tableHTML + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}

	// Extract the table data.
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
