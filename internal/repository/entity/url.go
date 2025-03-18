package entity

import "ozon_shortener/internal/domain"

type URL struct {
	Id          string `db:"id"`
	OriginalUrl string `db:"original_url"`
	Token       string `db:"token"`
}

func (u URL) ToDomain() domain.URL {
	return domain.URL{
		Id:          u.Id,
		OriginalUrl: u.OriginalUrl,
		Token:       u.Token,
	}
}

func (URL) FromDomain(u domain.URL) URL {
	return URL{
		Id:          u.Id,
		OriginalUrl: u.OriginalUrl,
		Token:       u.Token,
	}
}
