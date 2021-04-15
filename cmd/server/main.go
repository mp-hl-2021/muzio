package main

import (
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/accountrepo"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/entityrepo"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/playlistrepo"
	"github.com/mp-hl-2021/muzio/internal/usecases/account"
	"github.com/mp-hl-2021/muzio/internal/usecases/entity"
	"github.com/mp-hl-2021/muzio/internal/usecases/playlist"
	"net/http"
	"time"
)

func main() {
	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.NewMemory(),
	}
	entityUseCases := &entity.UseCases{
		EntityStorage: entityrepo.NewMemory(),
	}
	playlistUseCases := &playlist.UseCases{
		PlaylistStorage: playlistrepo.NewMemory(),
	}

	service := httpapi.NewApi(accountUseCases, entityUseCases, playlistUseCases)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
