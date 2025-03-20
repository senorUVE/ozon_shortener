package memory

import (
	"errors"
	"ozon_shortener/internal/repository"
	"ozon_shortener/internal/repository/entity"
)

type memoryUrlQuery struct {
	store *MemoryRepository
}

func NewMemoryUrlQuery(store *MemoryRepository) repository.UrlQuery {
	return &memoryUrlQuery{
		store: store,
	}
}

func (q *memoryUrlQuery) InsertUrl(u entity.URL) error {
	return q.store.InsertUrl(u)
}

func (q *memoryUrlQuery) GetUrlByOriginal(originalURL string) (entity.URL, error) {
	return q.store.GetUrlByOriginal(originalURL)
}

func (q *memoryUrlQuery) UpdateURL(u entity.URL) error {
	return q.store.UpdateURL(u)
}

func (q *memoryUrlQuery) GetByTokens(tokens []string) ([]entity.URL, error) {
	return q.store.GetByTokens(tokens)
}

func (q *memoryUrlQuery) GetUrlByPK(id string) (entity.URL, error) {
	return q.store.GetUrlByPK(id)
}

func (q *memoryUrlQuery) GetUrlByToken(token string) (entity.URL, error) {
	q.store.mutex.RLock()
	defer q.store.mutex.RUnlock()

	for _, u := range q.store.data {
		if u.Token == token {
			return u, nil
		}
	}
	return entity.URL{}, errors.New("not found")
}

func (q *memoryUrlQuery) InsertUrlReturning(u entity.URL) (entity.URL, error) {
	if err := q.store.InsertUrl(u); err != nil {
		return entity.URL{}, err
	}
	return q.store.GetUrlByOriginal(u.OriginalUrl)
}
