package random_test

import (
	"strings"
	"testing"

	"ozon_shortener/internal/utils/random"
)

func TestGenerateRandomLength(t *testing.T) {
	length := 148
	s, err := random.GenerateRandom(length)
	if err != nil {
		t.Fatalf("GenerateRandom returned error: %v", err)
	}
	if len(s) != length {
		t.Errorf("expected length %d, got %d", length, len(s))
	}
}

func TestGenerateRandomAllowedCharacters(t *testing.T) {
	length := 1000
	s, err := random.GenerateRandom(length)
	if err != nil {
		t.Fatalf("GenerateRandom returned error: %v", err)
	}
	const allowed = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	for i, ch := range s {
		if !strings.ContainsRune(allowed, ch) {
			t.Errorf("character %q at index %d is not allowed", ch, i)
		}
	}
}

func TestGenerateRandomUniqueness(t *testing.T) {
	length := 16
	s1, err1 := random.GenerateRandom(length)
	s2, err2 := random.GenerateRandom(length)
	if err1 != nil || err2 != nil {
		t.Fatalf("GenerateRandom returned error: %v, %v", err1, err2)
	}
	if s1 == s2 {
		t.Errorf("expected two different random strings, got identical: %q", s1)
	}
}
