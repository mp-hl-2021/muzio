package entity

type Link struct {
	ServiceName string
	Url         string
}

type MusicalEntity struct {
	Id     string
	Artist string
	Album  string
	Track  string
	Links  []Link
}

type Interface interface {
	CreateMusicalEntity(artist string, album string, track string) (MusicalEntity, error)
	GetMusicalEntityById(id string) (MusicalEntity, error)
}
