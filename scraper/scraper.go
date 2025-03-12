package scraper

// Scraper defines the common behavior for any stock data scraper.
type Scraper interface {
	// Scrape retrieves data for the given stockNumber.
	// Some implementations might use startDate/endDate while others ignore them.
	Scrape(stockNumber, startDate, endDate string) ([][]string, error)
}
