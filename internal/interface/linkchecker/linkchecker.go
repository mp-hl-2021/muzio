package linkchecker

import (
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/usecases/entity"
	"net/http"
)

type LinkChecker struct {
	MusicalEntityUseCases entity.Interface
	IdsToCheckChannel     <-chan string
}

func New(e entity.Interface, c <-chan string) *LinkChecker {
	return &LinkChecker{MusicalEntityUseCases: e, IdsToCheckChannel: c}
}

func (c *LinkChecker) CheckMusicalEntities() {
	go func() {
		for eid := range c.IdsToCheckChannel {
			e, err := c.MusicalEntityUseCases.GetMusicalEntityById(eid)
			if err != nil {
				continue
			}
			nl := make([]common.Link, 0, len(e.Links))
			for _, l := range e.Links {
				isAvailable := true
				_, err := http.Get(l.Url)
				if err != nil {
					isAvailable = false
				}
				l.IsAvailable = isAvailable
				nl = append(nl, l)
			}
			_ = c.MusicalEntityUseCases.UpdateLinks(eid, nl)
		}
	}()
}
