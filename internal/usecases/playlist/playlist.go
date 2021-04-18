package playlist

import "github.com/mp-hl-2021/muzio/internal/domain/playlist"

//type Playlist struct {
//	Id      string
//	Owner   string
//	Name    string
//	Content []string // ids of MusicalEntity-es
//}

type Interface interface {
	CreatePlaylist(owner, name string, content []string) (string, error)
	GetPlaylistById(id string) (playlist.Playlist, error)
	UpdatePlayList(/* owner, */ id, name string, content []string) error
	DeletePlayList(/* owner, */ id string) error
}

type UseCases struct {
	PlaylistStorage playlist.Interface
}

func (u *UseCases) CreatePlaylist(owner, name string, content []string) (string, error) {
	// TODO: Auth
	p, err := u.PlaylistStorage.CreatePlaylist(owner, name, content)
	if err != nil {
		return "", err
	}
	return p.Id, nil
}

func (u *UseCases) GetPlaylistById(id string) (playlist.Playlist, error) {
	//p, err := u.PlaylistStorage.GetPlaylistById(id)
	//if err != nil {
	//	return Playlist{}, err
	//}
	//return Playlist{Name: p.Name, Content: p.Content}, nil
	return u.PlaylistStorage.GetPlaylistById(id)
}

func (u *UseCases) UpdatePlayList(/* owner, */ id, name string, content []string) error {
	// TODO: Auth
	return u.PlaylistStorage.UpdatePlaylist(id, name, content)
}

func (u *UseCases) DeletePlayList(/* owner, */ id string) error {
	// TODO: Auth
	return u.PlaylistStorage.DeletePlaylist(id)
}
