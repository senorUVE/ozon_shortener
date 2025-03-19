package extract

import "errors"

var ErrInvalidShortUrl = errors.New("invalid short url")

func ExtractToken(domain, fullURL string) (string, error) {
	prefix := domain + "/"
	if len(fullURL) <= len(prefix) || fullURL[:len(prefix)] != prefix {
		return "", ErrInvalidShortUrl
	}
	return fullURL[len(prefix):], nil
}
