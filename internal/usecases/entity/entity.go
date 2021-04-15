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
	CreateMusicalEntity(artist string, album string, track string) (MusicalEntity, error)
	GetMusicalEntityById(id string) (MusicalEntity, error)
}

type UseCases struct {
	EntityStorage entity.Interface
}

func (u *UseCases) CreateMusicalEntity(artist string, album string, track string,) (MusicalEntity, error) {
	panic("implement me")
}

func (u *UseCases) GetMusicalEntityById(id string) (MusicalEntity, error) {
	panic("implement me")
}