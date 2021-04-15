package main

import (
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
	"net/http"
	"time"
)

func main() {
	service := httpapi.NewApi()

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
