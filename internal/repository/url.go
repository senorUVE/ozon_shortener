package repository

import (
	"ozon_shortener/internal/repository/entity"

	sq "github.com/Masterminds/squirrel"

	daolib "ozon_shortener/dao"
)

const urlTable = "url"

type UrlQuery interface {
	InsertUrl(url entity.URL) error
	GetUrlByPK(id string) (entity.URL, error)
	GetUrlByToken(token string) (entity.URL, error)
	GetUrlByOriginal(originalUrl string) (entity.URL, error)
	GetByTokens(tokens []string) ([]entity.URL, error)
	UpdateURL(url entity.URL) error
}

type urlQuery struct {
	*daolib.BaseQuery
}

func (q *urlQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *urlQuery) InsertUrl(url entity.URL) error {
	query := q.QueryBuilder().Insert(urlTable).Columns(
		"original_url",
		"token",
	).
		Values(
			url.OriginalUrl,
			url.Token,
		)

	_, err := q.Runner().Execx(q.Context(), query)
	if err != nil {
		return err
	}
	return nil
}

func (q *urlQuery) GetUrlByPK(id string) (entity.URL, error) {
	query := q.QueryBuilder().Select(
		"id",
		"original_url",
		"token",
	).From(urlTable).Where(sq.Eq{
		"id": id,
	})
	var url entity.URL
	if err := q.Runner().Getx(q.Context(), &url, query); err != nil {
		return entity.URL{}, err
	}
	return url, nil
}

func (q *urlQuery) GetUrlByToken(token string) (entity.URL, error) {
	query := q.QueryBuilder().Select(
		"id",
		"original_url",
		"token",
	).From(urlTable).Where(sq.Eq{
		"token": token,
	})
	var url entity.URL
	if err := q.Runner().Getx(q.Context(), &url, query); err != nil {
		return entity.URL{}, err
	}
	return url, nil
}

func (q *urlQuery) GetUrlByOriginal(originalUrl string) (entity.URL, error) {
	query := q.QueryBuilder().Select(
		"id",
		"original_url",
		"token",
	).From(urlTable).Where(sq.Eq{
		"original_url": originalUrl,
	})
	var url entity.URL
	if err := q.Runner().Getx(q.Context(), &url, query); err != nil {
		return entity.URL{}, err
	}
	return url, nil
}

func (q *urlQuery) GetByTokens(tokens []string) ([]entity.URL, error) {
	query := q.QueryBuilder().Select(
		"id",
		"original_url",
		"token",
	).From(urlTable).Where(sq.Eq{
		"token": tokens,
	})

	var urls []entity.URL
	if err := q.Runner().Selectx(q.Context(), &urls, query); err != nil {
		return nil, err
	}
	return urls, nil
}

func (q *urlQuery) UpdateURL(url entity.URL) error {
	query := q.QueryBuilder().Update(urlTable).Set(
		"token", url.Token,
	).Where(sq.Eq{
		"id": url.Id,
	})
	_, err := q.Runner().Execx(q.Context(), query)
	return err
}
