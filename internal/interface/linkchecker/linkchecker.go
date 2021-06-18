package linkchecker

import (
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/domain/entity"
	"net/http"
	"time"
)

type LinkChecker struct {
	EntityStorage entity.Interface
}

func New(e entity.Interface) *LinkChecker {
	return &LinkChecker{EntityStorage: e}
}

func (c *LinkChecker) Start(workers int) {
	ch := make(chan entity.MusicalEntity, workers)

	go func() {
		for range time.Tick(time.Second) {
			es, err := c.EntityStorage.GetBatchToCheck(workers)
			if err != nil {
				continue
			}
			for _, e := range es {
				ch <- e
			}
		}
	}()

	for i := 0; i < workers; i++ {
		go func() {
			e := <-ch
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
			c.EntityStorage.UpdateLinks(e.Id, nl)
		}()
	}
}
