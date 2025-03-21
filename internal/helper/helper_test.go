package helper_test

import (
	"github.com/ysonC/multi-stocks-download/internal/helper"
	"testing"
)

func TestCheckSpace(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		row  []string
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "Empty Row",
			row:  []string{},
			want: []string{},
		},
		{
			name: "Row with Empty String",
			row:  []string{"", "A", "B"},
			want: []string{"-", "A", "B"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.CheckSpace(tt.row)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CheckSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}
