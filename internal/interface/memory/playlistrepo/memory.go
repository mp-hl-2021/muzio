package playlistrepo

import (
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/playlist"
	"strconv"
	"sync"
)

type Memory struct {
	playlistsById map[string]playlist.Playlist
	nextId       uint64
	mu           *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		playlistsById: make(map[string]playlist.Playlist),
		mu:           &sync.Mutex{},
	}
}

func (m *Memory) CreatePlaylist(owner string, name string, content []string) (playlist.Playlist, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := playlist.Playlist{
		Id: strconv.FormatUint(m.nextId, 16),
		Owner: owner,
		Name: name,
		Content: content,
	}
	m.playlistsById[p.Id] = p
	m.nextId++
	return p, nil
}

func (m *Memory) GetPlaylistById(id string) (playlist.Playlist, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.playlistsById[id]
	if !ok {
		return p, domain.ErrNotFound
	}
	return p, nil
}

func (m *Memory) UpdatePlaylist(id string, name string, content []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.playlistsById[id]
	if !ok {
		return domain.ErrNotFound
	}
	p.Name = name
	p.Content = content
	m.playlistsById[id] = p
	return nil
}

func (m *Memory) DeletePlaylist(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.playlistsById[id]
	if !ok {
		return domain.ErrNotFound
	}
	delete(m.playlistsById, id)
	return nil
}
