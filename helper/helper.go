package helper

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// extractFullTableData parses the table HTML without skipping the header.
func ExtractFullTableData(html string) ([][]string, error) {
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
		row = checkSpace(row)
		if len(row) > 0 {
			data = append(data, row)
		}
	})
	return data, nil
}

// extractTableData wraps the HTML in a table tag and uses goquery to extract rows.
func ExtractTableData(html string, maxColumns int, skipHeader bool) ([][]string, error) {
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

func checkSpace(row []string) []string {
	for i, v := range row {
		if strings.TrimSpace(v) == "" {
			row[i] = "-"
		}
	}
	return row
}
