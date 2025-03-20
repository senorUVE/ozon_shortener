package bijection_test

import (
	"fmt"
	"testing"

	"ozon_shortener/internal/utils/bijection"
)

func TestConvertNumberToKey(t *testing.T) {
	tests := []struct {
		number int64
		want   string
	}{
		{0, ""},
		{1, "1"},
		{9, "9"},
		{10, "a"},
		{35, "z"},
		{36, "10"},
		{37, "11"},
		{100, "2s"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("number=%d", tt.number), func(t *testing.T) {
			got := bijection.ConvertNumberToKey(tt.number)
			if got != tt.want {
				t.Errorf("ConvertNumberToKey(%d) = %q; want %q", tt.number, got, tt.want)
			}
		})
	}
}

func TestConvertKeyToNumber(t *testing.T) {
	tests := []struct {
		key  string
		want int64
	}{
		{"", 0},
		{"1", 1},
		{"9", 9},
		{"a", 10},
		{"z", 35},
		{"10", 36},
		{"2s", 100},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("key=%s", tt.key), func(t *testing.T) {
			got := bijection.ConvertKeyToNumber(tt.key)
			if got != tt.want {
				t.Errorf("ConvertKeyToNumber(%q) = %d; want %d", tt.key, got, tt.want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	for i := int64(0); i < 2000; i++ {
		key := bijection.ConvertNumberToKey(i)
		got := bijection.ConvertKeyToNumber(key)
		if got != i {
			t.Fatalf("RoundTrip failed for %d -> %q -> %d", i, key, got)
		}
	}
}
