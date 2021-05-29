package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mp-hl-2021/muzio/internal/interface/httpapi"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var config = struct {
	address          string
	concurrencyLevel int
}{}

func init() {
	adress := flag.String("address", "http://localhost:8080", "muzio address")
	concurrencyLvl := flag.Int("concurrency", 50, "amount of concurrent requests")
	flag.Parse()

	config.address = *adress
	config.concurrencyLevel = *concurrencyLvl
}

func main() {

	//fmt.Println("URL:>", config.address)
	//
	//ctx, _ := context.WithCancel(context.Background())
	//
	//login := "a" + gofakeit.DigitN(8)
	//pass := gofakeit.Password(true, true, true, false, false, 16)
	//
	//c := client{
	//	c: http.Client{
	//		Timeout: 10 * time.Second,
	//	},
	//}
	//
	//err := c.createAccount(ctx, login, pass)
	//err = c.loginInto(ctx, login, pass)
	//
	//fmt.Println(err)

	////var jsonStr = []byte(`{"login":"asdf123", "password":"qwer234556qwert"}`)
	//var jsonStr = []byte(
	//	fmt.Sprintf("{\"login\":\"%s\", \"password\":\"%s\"}", login, pass))
	//req, err := http.NewRequest("POST", config.address + "/signup", bytes.NewBuffer(jsonStr))
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer func() {
		signal.Stop(ch)
		cancel()
	}()

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	c := client{
		c: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	var wg sync.WaitGroup
	wg.Add(config.concurrencyLevel)

	for i := 0; i < config.concurrencyLevel; i++ {
		go func(i int) {
			err := accountCreator(ctx, c)
			fmt.Printf("worker %d finished, err: %v\n", i, err)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println("All workers have been finished")

}

func accountCreator(ctx context.Context, c client) error {
	for {
		select {
		default:
			login := "a" + gofakeit.DigitN(8)
			pass := gofakeit.Password(true, true, true, false, false, 16)
			err := c.createAccount(ctx, login, pass)
			if err != nil {
				fmt.Println("Account creation failed:", err)
			}
			err = c.loginInto(ctx, login, pass)
			if err != nil {
				fmt.Println("Account login failed:", err)
			}
		case <- ctx.Done():
			fmt.Println("leaving worker")
			return ctx.Err()
		}
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

func (c client) loginInto(ctx context.Context, login string, pass string) error {
	body := httpapi.PostSignupRequestModel{Login: login, Password: pass}

	s, err := json.Marshal(body)
	if err != nil{
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/signin", bytes.NewReader(s))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")
	response, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Account creation failed: %v", response.Status)
	}
	return nil
}
