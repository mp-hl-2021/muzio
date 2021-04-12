package entity

import (
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
)

type Link struct {
	ServiceName string
	Url         string
}

type MusicalEntity struct {
	Artist string
	Album  string
	Track  string
	Links  []Link
}

type Interface interface {
	CreateMusicalEntity(artist, album, track string, links []Link) (string, error)
	GetMusicalEntityById(id string) (MusicalEntity, error)
}

type UseCases struct {
	EntityStorage entity.Interface
}

func (u *UseCases) CreateMusicalEntity(artist, album, track string, links []Link) (string, error) {
	nl := make([]entity.Link, 0, len(links))
	for _, l := range links {
		nl = append(nl, entity.Link{
			ServiceName: l.ServiceName,
			Url: l.Url,
		})
	}
	e, err := u.EntityStorage.CreateMusicalEntity(artist, album, track, nl)
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
	nl := make([]Link, 0, len(e.Links))
	for _, l := range e.Links {
		nl = append(nl, Link{
			ServiceName: l.ServiceName,
			Url: l.Url,
		})
	}
	return MusicalEntity{Artist: e.Artist, Album: e.Album, Track: e.Track, Links: nl}, err
}
