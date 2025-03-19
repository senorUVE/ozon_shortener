package memory

import (
	"context"
	"errors"

	"ozon_shortener/internal/repository"

	daolib "github.com/WantBeASleep/med_ml_lib/dao"
)

type MemoryDAO struct {
	store *MemoryRepository
}

func NewMemoryDAO(store *MemoryRepository) repository.DAO {
	return &MemoryDAO{
		store: store,
	}
}

func (m *MemoryDAO) BeginTx(ctx context.Context, opts ...daolib.TxOption) (context.Context, error) {
	return ctx, errors.New("transactions not supported in memory storage")
}

func (m *MemoryDAO) RollbackTx(ctx context.Context) error {
	return nil
}

func (m *MemoryDAO) CommitTx(ctx context.Context) error {
	return nil
}

func (m *MemoryDAO) NewRepo(ctx context.Context, querier daolib.BaseQuerySetter) {
	// no-op
}

func (m *MemoryDAO) NewUrlQuery(ctx context.Context) repository.UrlQuery {
	return NewMemoryUrlQuery(m.store)
}
