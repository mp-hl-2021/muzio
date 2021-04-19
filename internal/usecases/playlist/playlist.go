package playlist

import (
	"github.com/mp-hl-2021/muzio/internal/domain"
	"github.com/mp-hl-2021/muzio/internal/domain/playlist"
)

type Playlist struct {
	Name    string
	Content []string
}

type Interface interface {
	CreatePlaylist(owner, name string, content []string) (string, error)
	GetPlaylistById(id string) (Playlist, error)
	UpdatePlayList(owner, id, name string, content []string) error
	DeletePlayList(owner, id string) error
}

type UseCases struct {
	PlaylistStorage playlist.Interface
}

func (u *UseCases) CreatePlaylist(owner, name string, content []string) (string, error) {
	p, err := u.PlaylistStorage.CreatePlaylist(owner, name, content)
	if err != nil {
		return "", err
	}
	return p.Id, nil
}

func (u *UseCases) GetPlaylistById(id string) (Playlist, error) {
	p, err := u.PlaylistStorage.GetPlaylistById(id)
	if err != nil {
		return Playlist{}, err
	}
	return Playlist{Name: p.Name, Content: p.Content}, nil
}

func (u *UseCases) UpdatePlayList(owner, id, name string, content []string) error {
	p, err := u.PlaylistStorage.GetPlaylistById(id)
	if err != nil {
		return err
	}
	if p.Owner != owner {
		return domain.ErrForbidden
	}
	return u.PlaylistStorage.UpdatePlaylist(id, name, content)
}

func (u *UseCases) DeletePlayList(owner, id string) error {
	p, err := u.PlaylistStorage.GetPlaylistById(id)
	if err != nil {
		return err
	}
	if p.Owner != owner {
		return domain.ErrForbidden
	}
	return u.PlaylistStorage.DeletePlaylist(id)
}
