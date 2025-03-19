package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"ozon_shortener/internal/repository/entity"
)

const (
	maxEntries  = 10000
	coldStorage = "cold_data.json"
)

type MemoryRepository struct {
	data  map[string]entity.URL
	mutex sync.RWMutex
	auto  int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]entity.URL),
	}
}

func (m *MemoryRepository) InsertUrl(u entity.URL) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.auto++
	u.Id = fmt.Sprintf("%d", m.auto)
	if len(m.data) >= maxEntries {
		go m.flushToDisk()
	}

	m.data[u.OriginalUrl] = u
	return nil
}

func (m *MemoryRepository) GetUrlByOriginal(originalURL string) (entity.URL, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if url, ok := m.data[originalURL]; ok {
		return url, nil
	}

	coldData, err := m.loadFromDisk()
	if err == nil {
		if url, ok := coldData[originalURL]; ok {
			return url, nil
		}
	}

	return entity.URL{}, errors.New("not found")
}

func (m *MemoryRepository) UpdateURL(u entity.URL) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.data[u.OriginalUrl]; !ok {
		return errors.New("url not found")
	}

	m.data[u.OriginalUrl] = u
	return nil
}

func (m *MemoryRepository) GetByTokens(tokens []string) ([]entity.URL, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var results []entity.URL
	for _, u := range m.data {
		for _, token := range tokens {
			if u.Token == token {
				results = append(results, u)
			}
		}
	}

	coldData, err := m.loadFromDisk()
	if err == nil {
		for _, u := range coldData {
			for _, token := range tokens {
				if u.Token == token {
					results = append(results, u)
				}
			}
		}
	}

	return results, nil
}

func (m *MemoryRepository) flushToDisk() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	dataToFlush := make(map[string]entity.URL)
	counter := 0
	for k, v := range m.data {
		dataToFlush[k] = v
		delete(m.data, k)
		counter++
		if counter >= maxEntries/2 {
			break
		}
	}

	file, err := os.Create(coldStorage)
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(dataToFlush)
}

func (m *MemoryRepository) GetUrlByPK(id string) (entity.URL, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, u := range m.data {
		if u.Id == id {
			return u, nil
		}
	}

	coldData, err := m.loadFromDisk()
	if err == nil {
		for _, u := range coldData {
			if u.Id == id {
				return u, nil
			}
		}
	}

	return entity.URL{}, errors.New("not found")
}

func (m *MemoryRepository) loadFromDisk() (map[string]entity.URL, error) {
	file, err := os.Open(coldStorage)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var coldData map[string]entity.URL
	err = json.NewDecoder(file).Decode(&coldData)
	if err != nil {
		return nil, err
	}

	return coldData, nil
}
