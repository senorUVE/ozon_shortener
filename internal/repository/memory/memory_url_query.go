package memory

import (
	"errors"
	"ozon_shortener/internal/repository"
	"ozon_shortener/internal/repository/entity"
)

// memoryUrlQuery реализует интерфейс repository.UrlQuery для in‑memory хранилища.
type memoryUrlQuery struct {
	store *MemoryRepository
}

// NewMemoryUrlQuery создаёт новый in‑memory UrlQuery, используя предоставленное хранилище.
func NewMemoryUrlQuery(store *MemoryRepository) repository.UrlQuery {
	return &memoryUrlQuery{
		store: store,
	}
}

// InsertUrl вставляет новую запись.
func (q *memoryUrlQuery) InsertUrl(u entity.URL) error {
	return q.store.InsertUrl(u)
}

// GetUrlByOriginal ищет запись по оригинальному URL.
func (q *memoryUrlQuery) GetUrlByOriginal(originalURL string) (entity.URL, error) {
	return q.store.GetUrlByOriginal(originalURL)
}

// UpdateURL обновляет запись.
func (q *memoryUrlQuery) UpdateURL(u entity.URL) error {
	return q.store.UpdateURL(u)
}

// GetByTokens возвращает записи, у которых поле Token совпадает с одним из переданных.
func (q *memoryUrlQuery) GetByTokens(tokens []string) ([]entity.URL, error) {
	return q.store.GetByTokens(tokens)
}

// GetUrlByPK ищет запись по автоинкрементному идентификатору.
func (q *memoryUrlQuery) GetUrlByPK(id string) (entity.URL, error) {
	return q.store.GetUrlByPK(id)
}

// GetUrlByToken ищет запись по полю Token.
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
