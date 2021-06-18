package linkchecker

import (
	"fmt"
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

func (c *LinkChecker) Start(workers int, interval time.Duration) {
	ch := make(chan entity.MusicalEntity, workers)

	go func() {
		for range time.Tick(interval) {
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
			for {
				e := <-ch
				nl := make([]common.Link, 0, len(e.Links))
				for _, l := range e.Links {
					isAvailable := true
					resp, err := http.Get(l.Url)
					if err != nil || resp.StatusCode != http.StatusOK{
						isAvailable = false
					}
					l.IsAvailable = isAvailable
					fmt.Printf("Link %s:%s isAvailable: %t\n", l.ServiceName, l.Url, l.IsAvailable)
					nl = append(nl, l)
				}
				_ = c.EntityStorage.UpdateLinks(e.Id, nl)
			}
		}()
	}
}
