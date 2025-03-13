package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
	return extractTableData(html, 6, true)
}

// extractTableData wraps the HTML in a table tag and uses goquery to extract rows.
func extractTableData(html string, maxColumns int, skipHeader bool) ([][]string, error) {
	var data [][]string
	wrappedHTML := "<table>" + html + "</table>"
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrappedHTML))
	if err != nil {
		return nil, err
	}
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if skipHeader && i == 0 {
			return
		}
		var row []string
		s.Find("td").Each(func(j int, cell *goquery.Selection) {
			if j < maxColumns {
				row = append(row, strings.TrimSpace(cell.Text()))
			}
		})
		// Optionally filter out rows (e.g., those ending with "W53").
		if len(row) > 0 && !strings.HasSuffix(row[0], "W53") {
			data = append(data, row)
		}
	})
	return data, nil
}
