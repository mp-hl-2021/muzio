package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mp-hl-2021/muzio/internal/common"
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
	"io"
	"net/http"
	"time"
)

var config = struct {
	address          string
}{}

func init() {
	adress := flag.String("address", "http://localhost:8080", "muzio address")
	flag.Parse()

	config.address = *adress
}

func main() {

	ctx, _ := context.WithCancel(context.Background())

	login := "a" + gofakeit.DigitN(8)
	pass := gofakeit.Password(true, true, true, false, false, 16)

	c := client{
		c: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	err := c.createAccount(ctx, login, pass)
	if err != nil {
		panic(err)
	}

	_, err = c.loginInto(ctx, login, pass)
	if err != nil {
		panic(err)
	}

	lnks := [...]string{
		"https://music.yandex.ru/album/9358192/track/60532413",
		"https://music.yandex.ru/album/945335243523452435234524523452452452452345245",
		"https://music.yandex.ru/album/9358192/track/60532413",
		"https://music.yandex.ru/album/9358192/track/60531620",
		"https://music.yandex.ru/album/9358192/track/6053162083459209348560928374502837456082374560972364059",
		"длывоаждфывваждофиывадоооывафждлаофжыдлвлтаждфоыит",
	}

	ids := []string{}

	for _, l := range lnks {

		links := []common.Link{{
			ServiceName: "yandex",
			Url:         l,
		}}

		id, err := c.addSong(ctx, links)
		if err != nil {
			panic(err)
		}

		ids = append(ids, id)
	}

	time.Sleep(20 * time.Second)

	for _, id := range ids {
		song, err := c.getSong(ctx, id)
		if err != nil {
			panic(err)
		}

		fmt.Println(song)
	}

}

type client struct {
	c http.Client
}

func (c client) createAccount(ctx context.Context, login string, pass string) error {
	body := httpapi.PostSignupRequestModel{Login: login, Password: pass}

	s, err := json.Marshal(body)
	if err != nil{
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/signup", bytes.NewReader(s))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")
	response, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("Account creation failed: %v", response.Status)
	}
	return nil
}

func (c client) loginInto(ctx context.Context, login string, pass string) ([]byte, error) {
	body := httpapi.PostSignupRequestModel{Login: login, Password: pass}

	s, err := json.Marshal(body)
	if err != nil{
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/signin", bytes.NewReader(s))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json")
	response, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Login failed: %v", response.Status)
	}

	tkn, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return tkn, nil
}

func (c client) addSong (ctx context.Context, links []common.Link) (string, error) {
	track := "c" + gofakeit.DigitN(8)
	fmt.Println(track)
	body := httpapi.GetMusicalEntityResponseModel{
		Artist: "a" + gofakeit.DigitN(8),
		Album:  "b" + gofakeit.DigitN(8),
		Track:  track,
		Links:  links,
	}

	s, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/drop/music", bytes.NewReader(s))
	if err != nil {
		return "", err
	}
	response, err := c.c.Do(req)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	var res httpapi.PostMusicalEntityResponseModel

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Adding song failed: %v", response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("Decoding failed: %v", response.Status)
	}

	return res.Id, nil
}

func (c client) getSong (ctx context.Context, id string) (*httpapi.GetMusicalEntityResponseModel, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.address+"/music/" + id, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var res httpapi.GetMusicalEntityResponseModel

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Getting failed: %v", response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("Decoding failed: %v", response.Status)
	}

	return &res, nil
}