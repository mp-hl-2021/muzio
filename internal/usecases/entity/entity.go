package entity

import (
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
)

type MusicalEntity struct {
	Artist string
	Album  string
	Track  string
	Links  []common.Link
}

type Interface interface {
	CreateMusicalEntity(artist, album, track string, links []common.Link) (string, error)
	GetMusicalEntityById(id string) (MusicalEntity, error)
}

type UseCases struct {
	EntityStorage entity.Interface
}

func (u *UseCases) CreateMusicalEntity(artist, album, track string, links []common.Link) (string, error) {
	e, err := u.EntityStorage.CreateMusicalEntity(artist, album, track, links)
	if err != nil {
		return "", err
	}
	return e.Id, nil
}

func (u *UseCases) GetMusicalEntityById(id string) (MusicalEntity, error) {
	e, err := u.EntityStorage.GetMusicalEntityById(id)
	if err != nil {
		return MusicalEntity{}, err
	}
	return MusicalEntity{Artist: e.Artist, Album: e.Album, Track: e.Track, Links: e.Links}, err
}
