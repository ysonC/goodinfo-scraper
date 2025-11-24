package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const failedFileName = "failed.txt"

// SaveFailedStocks writes the failed stock list to failed.txt inside the provided directory.
// An empty list removes any existing file to avoid rerunning stale failures.
func SaveFailedStocks(dir string, stocks []string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	filePath := filepath.Join(dir, failedFileName)

	if len(stocks) == 0 {
		err := os.Remove(filePath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	}

	data := strings.Join(stocks, "\n")
	return os.WriteFile(filePath, []byte(data), 0o644)
}

// LoadFailedStocks reads failed.txt inside the provided directory.
// If the file does not exist, it returns an empty slice without error.
func LoadFailedStocks(dir string) ([]string, error) {
	filePath := filepath.Join(dir, failedFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	var stocks []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			stocks = append(stocks, trimmed)
		}
	}
	return stocks, nil
}
