package playlist

type Playlist struct {
	Id      string
	Owner   string
	Name    string
	Content []string // ids of MusicalEntity-es
}

type Interface interface {
	CreatePlaylist(owner, name string, content []string) (Playlist, error)
	GetPlaylistById(id string) (Playlist, error)
	UpdatePlaylist(id, name string, content []string) error
	DeletePlaylist(id string) error
}
