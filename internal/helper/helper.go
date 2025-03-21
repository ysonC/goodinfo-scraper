package helper

import (
	"strings"
)

func CheckSpace(row []string) ([]string, error) {
	for i, v := range row {
		if strings.TrimSpace(v) == "" {
			row[i] = "-"
		}
	}
	return row, nil
}
