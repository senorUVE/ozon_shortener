package repository

import (
	"ozon_shortener/internal/repository/entity"

	sq "github.com/Masterminds/squirrel"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
)

const urlTable = "url"

type UrlQuery interface {
	InsertUrl(url entity.URL) error
	GetUrlByPK(id string) (entity.URL, error)
	GetUrlByToken(token string) (entity.URL, error)
}

type urlQuery struct {
	*daolib.BaseQuery
}

func (q *urlQuery) SetBaseQuery(baseQuery *daolib.BaseQuery) {
	q.BaseQuery = baseQuery
}

func (q *urlQuery) InsertUrl(url entity.URL) error {
	query := q.QueryBuilder().Insert(urlTable).Columns(
		"id",
		"original_url",
		"token",
	).
		Values(
			url.Id,
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
