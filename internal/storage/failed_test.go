package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSaveAndLoadFailedStocks(t *testing.T) {
	dir := t.TempDir()
	original := []string{"2330", "0050"}

	if err := SaveFailedStocks(dir, original); err != nil {
		t.Fatalf("SaveFailedStocks returned error: %v", err)
	}

	loaded, err := LoadFailedStocks(dir)
	if err != nil {
		t.Fatalf("LoadFailedStocks returned error: %v", err)
	}

	if !reflect.DeepEqual(original, loaded) {
		t.Fatalf("Loaded stocks mismatch, expected %v got %v", original, loaded)
	}

	if err := SaveFailedStocks(dir, []string{}); err != nil {
		t.Fatalf("SaveFailedStocks (empty) returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, failedFileName)); err == nil {
		t.Fatalf("expected failed file to be removed on empty save")
	}

	emptyLoaded, err := LoadFailedStocks(dir)
	if err != nil {
		t.Fatalf("LoadFailedStocks returned error after delete: %v", err)
	}
	if len(emptyLoaded) != 0 {
		t.Fatalf("expected empty slice after delete, got %v", emptyLoaded)
	}
}
