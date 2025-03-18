package repository

import (
	"context"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
	"github.com/jmoiron/sqlx"
)

type DAO interface {
	daolib.DAO
	NewUrlQuery(ctx context.Context) UrlQuery
}

type dao struct {
	daolib.DAO
}

func NewRepository(psql *sqlx.DB) DAO {
	return &dao{DAO: daolib.NewDao(psql)}
}

func (d *dao) NewUrlQuery(ctx context.Context) UrlQuery {
	urlQuery := &urlQuery{}
	d.NewRepo(ctx, urlQuery)

	return urlQuery
}
