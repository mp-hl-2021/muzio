package entityrepo

import (
	"container/heap"
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
	"strconv"
	"sync"
)

type Memory struct {
	entitiesById  map[string]entity.MusicalEntity
	entitiesQueue MusicalEntityHeap
	nextId        uint64
	nextCheckedAt uint64
	mu            *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		entitiesById:  make(map[string]entity.MusicalEntity),
		entitiesQueue: make([]MusicalEntityItem, 0),
		mu:            &sync.Mutex{},
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
	heap.Push(&m.entitiesQueue, MusicalEntityItem{e.Id, 0})
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

func (m *Memory) GetBatchToCheck(number int) ([]entity.MusicalEntity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]entity.MusicalEntity, 0, number)
	for i := 0; i < number; i++ {
		if m.entitiesQueue.Len() == 0 {
			break
		}
		nextEntity := heap.Pop(&m.entitiesQueue).(MusicalEntityItem)
		e, ok := m.entitiesById[nextEntity.id]
		if !ok { // skip if entity was deleted
			i--
			continue
		}
		result = append(result, e)
	}

	// update check time
	m.nextCheckedAt++
	for _, r := range result {
		heap.Push(&m.entitiesQueue, MusicalEntityItem{r.Id, m.nextCheckedAt})
	}

	return result, nil
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
