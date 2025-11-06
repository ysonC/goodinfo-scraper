package scraper

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"

	"github.com/ysonC/multi-stocks-download/internal/helper"
)

// BaseScraper encapsulates shared browser and page logic.
type BaseScraper struct {
	pw *playwright.Playwright
}

// NewBaseScraper returns a new BaseScraper.
func NewBaseScraper(pw *playwright.Playwright) *BaseScraper {
	return &BaseScraper{pw: pw}
}

// fetchHTML launches a browser, navigates to the URL, waits for the table,
// and returns the inner HTML of the table element.
func (b *BaseScraper) fetchHTML(url string) (string, error) {
	browser, err := b.pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-setuid-sandbox"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("failed to create page: %w", err)
	}

	if _, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		return "", fmt.Errorf("failed to goto URL: %w", err)
	}

	tableLocator := page.Locator("#tblDetail")
	if err := tableLocator.WaitFor(); err != nil {
		log.Printf("Warning: table not found")
	}

	// Now wait until the table has at least 200 rows.
	// 	if _, err := page.WaitForFunction(`() => {
	//     return document.querySelector("#tblDetail").querySelectorAll("tr").length >= 200;
	// }`, nil); err != nil {
	// 		return "", fmt.Errorf("timeout waiting for table data: %w", err)
	// 	}
	//

	html, err := tableLocator.InnerHTML()
	if err != nil {
		return "", fmt.Errorf("failed to get table HTML: %w", err)
	}
	return html, nil
}

// extractFullTableData parses the table HTML without skipping the header.
func (b *BaseScraper) extractFullTableData(html string) ([][]string, error) {
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
		row, err = helper.CheckSpace(row)
		if err != nil {
			return
		}
		if len(row) > 0 {
			data = append(data, row)
		}
	})
	return data, nil
}

// extractTableData wraps the HTML in a table tag and uses goquery to extract rows.
func (b *BaseScraper) extractTableData(
	html string,
	maxColumns int,
	skipHeader bool,
) ([][]string, error) {
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
		// if len(row) > 0 && !strings.HasSuffix(row[0], "W53") {
		if len(row) > 0 {
			data = append(data, row)
		}
	})
	return data, nil
}
