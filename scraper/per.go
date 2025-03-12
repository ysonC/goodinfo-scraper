package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

// PERScraper implements the scraper for PER data.
type PERScraper struct {
	pw *playwright.Playwright
}

// NewPERScraper returns a new PERScraper.
func NewPERScraper(pw *playwright.Playwright) *PERScraper {
	return &PERScraper{pw: pw}
}

// Scrape navigates to the PER URL, extracts the table HTML, and parses it.
func (p *PERScraper) Scrape(stockNumber, startDate, endDate string) ([][]string, error) {
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
		"https://goodinfo.tw/tw/ShowK_ChartFlow.asp?RPT_CAT=PER&STEP=DATA&STOCK_ID=%s&CHT_CAT=WEEK&PRICE_ADJ=F&START_DT=%s&END_DT=%s",
		stockNumber,
		startDate,
		endDate,
	)
	if _, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return nil, fmt.Errorf("failed to goto URL: %w", err)
	}

	tableLocator := page.Locator("#tblDetail")
	if err := tableLocator.WaitFor(); err != nil {
		log.Printf("Warning: table not found for stock %s", stockNumber)
	}

	tableHTML, err := tableLocator.InnerHTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get table HTML: %w", err)
	}

	return extractTableData(tableHTML)
}

// extractTableData wraps the fragment in a table and parses it.
func extractTableData(html string) ([][]string, error) {
	var data [][]string
	wrappedHTML := "<table>" + html + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}
	// Iterate over table rows.
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// Skip header row.
			return
		}
		var row []string
		s.Find("td").Each(func(j int, cell *goquery.Selection) {
			// For PER, keep the first six columns.
			if j < 6 {
				row = append(row, strings.TrimSpace(cell.Text()))
			}
		})
		// Optionally filter out rows (e.g. those ending with "W53").
		if len(row) > 0 && !strings.HasSuffix(row[0], "W53") {
			data = append(data, row)
		}
	})
	return data, nil
}
