package api

import (
	"testing"
	"time"
)

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		date     string
		expected string
	}{
		{"2024-11-07", "2024-11-04"},
		{"2024-11-10", "2024-11-04"},
		{"2024-11-03", "2024-10-28"},
		{"2024-11-01", "2024-10-28"},
	}

	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			inputDate, err := time.Parse(ApiDateFormat, tt.date)
			if err != nil {
				t.Fatalf("Failed to parse date: %v", err)
			}

			got := StartOfWeek(inputDate)

			expectedDate, err := time.Parse(ApiDateFormat, tt.expected)
			if err != nil {
				t.Fatalf("Failed to parse expected date: %v", err)
			}

			if !got.Equal(expectedDate) {
				t.Errorf("StartOfWeek(%v) = %v; want %v", inputDate, got, expectedDate)
			}
		})
	}
}
