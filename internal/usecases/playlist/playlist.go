package playlist

import "github.com/mp-hl-2021/muzio/internal/domain/playlist"

type Playlist struct {
	Id      string
	Content []string
}

type Interface interface {
	CreatePlaylist(owner string, content []string) (Playlist, error)
	GetPlaylistById(id string) (Playlist, error)
	UpdatePlayList(owner string, id string, content []string) error
	DeletePlayList(owner string, id string) error
}

type UseCases struct {
	PlaylistStorage playlist.Interface
}

func (u *UseCases) CreatePlaylist(owner string, content []string) (Playlist, error) {
	panic("implement me")
}

func (u *UseCases) GetPlaylistById(id string) (Playlist, error) {
	panic("implement me")
}

func (u *UseCases) UpdatePlayList(owner string, id string, content []string) error {
	panic("implement me")
}

func (u *UseCases) DeletePlayList(owner string, id string) error {
	panic("implement me")
}
