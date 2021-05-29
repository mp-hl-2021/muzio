package entityrepo

import (
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
	"strconv"
	"sync"
)

type Memory struct {
	entitiesById map[string]entity.MusicalEntity
	nextId       uint64
	mu           *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		entitiesById: make(map[string]entity.MusicalEntity),
		mu:           &sync.Mutex{},
	}
}

func (m *Memory) CreateMusicalEntity(artist string, album string, track string, links []common.Link) (entity.MusicalEntity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e := entity.MusicalEntity{
		Id: strconv.FormatUint(m.nextId, 16),
		Artist: artist,
		Album: album,
		Track: track,
		Links: links,
	}
	m.entitiesById[e.Id] = e
	m.nextId++
	return e, nil
}

func (m *Memory) GetMusicalEntityById(id string) (entity.MusicalEntity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entitiesById[id]
	if !ok {
		return e, domain.ErrNotFound
	}
	return e, nil
}

func (m *Memory) UpdateLinks(id string, links []common.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.entitiesById[id]
	if !ok {
		return domain.ErrNotFound
	}
	e.Links = links
	m.entitiesById[id] = e
	return nil
}
