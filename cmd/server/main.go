package main

import (
	"flag"
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
	"github.com/mp-hl-2021/muzio/internal/interface/linkchecker"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/accountrepo"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/entityrepo"
	"github.com/mp-hl-2021/muzio/internal/interface/memory/playlistrepo"
	"github.com/mp-hl-2021/muzio/internal/service/auth"
	"github.com/mp-hl-2021/muzio/internal/usecases/account"
	"github.com/mp-hl-2021/muzio/internal/usecases/entity"
	"github.com/mp-hl-2021/muzio/internal/usecases/playlist"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := auth.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.NewMemory(),
		AuthToken: a,
	}
	entityUseCases := &entity.UseCases{
		EntityStorage: entityrepo.NewMemory(),
	}
	playlistUseCases := &playlist.UseCases{
		PlaylistStorage: playlistrepo.NewMemory(),
	}

	linkChecker := linkchecker.New(entityUseCases.EntityStorage)
	linkChecker.Start(10, 10 * time.Second)

	service := httpapi.NewApi(accountUseCases, entityUseCases, playlistUseCases)

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
