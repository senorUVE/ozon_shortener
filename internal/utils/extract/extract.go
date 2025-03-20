package extract

import (
	"errors"
	"net/url"
	"path"
)

var ErrInvalidShortUrl = errors.New("invalid short url")

func ExtractToken(fullURL string) (string, error) {
	parsed, err := url.Parse(fullURL)
	if err != nil {
		return "", ErrInvalidShortUrl
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", ErrInvalidShortUrl
	}
	token := path.Base(parsed.Path)
	if token == "" || token == "/" {
		return "", ErrInvalidShortUrl
	}
	return token, nil
}
