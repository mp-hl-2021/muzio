package entity

import "github.com/mp-hl-2021/muzio/internal/common"

type MusicalEntity struct {
	Id     string
	Artist string
	Album  string
	Track  string
	Links  []common.Link
}

type Interface interface {
	CreateMusicalEntity(artist, album, track string, links []common.Link) (MusicalEntity, error)
	GetMusicalEntityById(id string) (MusicalEntity, error)
	GetBatchToCheck(number int) ([]MusicalEntity, error)
	UpdateLinks(id string, links []common.Link) error
}
