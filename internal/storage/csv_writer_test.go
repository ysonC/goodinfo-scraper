package storage

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a mock CSV file.
func createMockCSVFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func TestCombineAllCSVInFolder(t *testing.T) {
	tests := []struct {
		name       string // description of this test case
		folderPath string
		wantErr    bool
	}{
		{
			name:       "Valid CSV files exist",
			folderPath: "testdata/valid", // Will be created in test setup
			wantErr:    false,
		},
		{
			name:       "Missing required CSV file",
			folderPath: "testdata/missing", // Will be created with missing file
			wantErr:    true,
		},
		{
			name:       "Invalid folder path",
			folderPath: "testdata/nonexistent",
			wantErr:    true,
		},
	}

	// Setup test directories
	os.MkdirAll("testdata/valid", 0755)
	os.MkdirAll("testdata/missing", 0755)

	// Create valid CSV files
	validFiles := map[string]string{
		"per.csv":            "date,value\n2024-01-01,10",
		"stockdata.csv":      "date,stock\n2024-01-01,100",
		"monthlyrevenue.csv": "month,revenue\n2024-01,5000",
		"cashflow.csv":       "month,cashflow\n2024-01,2000",
	}

	for name, content := range validFiles {
		createMockCSVFile(filepath.Join("testdata/valid", name), content)
	}

	// Create a test folder with a missing file
	createMockCSVFile(filepath.Join("testdata/missing", "per.csv"), "date,value\n2024-01-01,10")
	createMockCSVFile(
		filepath.Join("testdata/missing", "stockdata.csv"),
		"date,stock\n2024-01-01,100",
	)

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := CombineAllCSVInFolder(tt.folderPath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CombineAllCSVInFolder() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CombineAllCSVInFolder() succeeded unexpectedly")
			}
		})
	}

	// Cleanup
	os.RemoveAll("testdata")
}
