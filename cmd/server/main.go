package main

import (
	"flag"
	"fmt"
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
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
	privateKeyPath := flag.String("private_key", "app.rsa", "file path")
	publicKeyPath := flag.String("public_key", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	a, err := auth.NewJwt(privateKeyBytes, publicKeyBytes, 100 * time.Minute)
	if err != nil {
		fmt.Println(err)
		return
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

	service := httpapi.NewApi(accountUseCases, entityUseCases, playlistUseCases)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
