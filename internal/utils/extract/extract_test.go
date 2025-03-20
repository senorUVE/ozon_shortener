package extract_test

import (
	"errors"
	"testing"

	"ozon_shortener/internal/utils/extract" // подстройте путь под реальную структуру
)

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{
			name:    "valid URL with token",
			input:   "http://localhost:8080/AbCdE",
			want:    "AbCdE",
			wantErr: nil,
		},
		{
			name:    "valid URL with trailing slash (no token)",
			input:   "http://localhost:8080/",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
		{
			name:    "no scheme",
			input:   "localhost:8080/AbCdE",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
		{
			name:    "no host",
			input:   "http:///AbCdE",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
		{
			name:    "invalid URL parse",
			input:   "http://%zzzzz",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
		{
			name:    "path is slash only",
			input:   "http://localhost:8080/",
			want:    "",
			wantErr: extract.ErrInvalidShortUrl,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extract.ExtractToken(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ExtractToken(%q) error = %v; wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ExtractToken(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}
