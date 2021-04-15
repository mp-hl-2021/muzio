package playlist

type Playlist struct {
	Id      string
	Owner   string
	Content []string // ids of MusicalEntity-es
}

type Interface interface {
	CreatePlaylist(owner string, content []string) (Playlist, error)
	GetPlaylistById(id string) (Playlist, error)
	UpdatePlayList(id string, content []string) error
	DeletePlayList(id string) error
}
