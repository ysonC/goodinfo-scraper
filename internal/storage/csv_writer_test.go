package storage

import (
	"reflect"
	"testing"
)

func TestMergeCSVData(t *testing.T) {
	tests := []struct {
		name     string
		csv1     [][]string
		csv2     [][]string
		expected [][]string
		wantErr  bool
	}{
		{
			name: "Equal Rows",
			csv1: [][]string{
				{"A1", "B1"},
				{"A2", "B2"},
			},
			csv2: [][]string{
				{"X1", "Y1"},
				{"X2", "Y2"},
			},
			expected: [][]string{
				{"A1", "B1", "", "X1", "Y1"},
				{"A2", "B2", "", "X2", "Y2"},
			},
			wantErr: false,
		},
		{
			name: "CSV1 longer than CSV2",
			csv1: [][]string{
				{"A1", "B1"},
				{"A2", "B2"},
				{"A3", "B3"},
			},
			csv2: [][]string{
				{"X1", "Y1"},
			},
			expected: [][]string{
				{"A1", "B1", "", "X1", "Y1"},
				{"A2", "B2", "", "-"},
				{"A3", "B3", "", "-"},
			},
			wantErr: false,
		},
		{
			name:     "Both CSVs Empty",
			csv1:     [][]string{},
			csv2:     [][]string{},
			expected: nil,
			wantErr:  true,
		},
		{
			name: "One CSV Empty",
			csv1: [][]string{
				{"A1", "B1"},
				{"A2", "B2"},
			},
			csv2:     [][]string{},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeCSVData(tt.csv1, tt.csv2)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeCSVData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("MergeCSVData() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
