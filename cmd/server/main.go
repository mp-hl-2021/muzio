package main

import (
	"net/http"
	"time"
)

func main() {
	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
