package scraper

import (
	"fmt"
	"log"

	"github.com/playwright-community/playwright-go"
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
		WaitUntil: playwright.WaitUntilStateLoad,
	}); err != nil {
		return "", fmt.Errorf("failed to goto URL: %w", err)
	}

	tableLocator := page.Locator("#tblDetail")
	if err := tableLocator.WaitFor(); err != nil {
		log.Printf("Warning: table not found")
	}

	// Now wait until the table has at least 200 rows.
	if _, err := page.WaitForFunction(`() => {
    return document.querySelector("#tblDetail").querySelectorAll("tr").length >= 200;
}`, nil); err != nil {
		return "", fmt.Errorf("timeout waiting for table data: %w", err)
	}

	html, err := tableLocator.InnerHTML()
	if err != nil {
		return "", fmt.Errorf("failed to get table HTML: %w", err)
	}
	return html, nil
}
