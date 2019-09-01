package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergey-chebanov/fire/gopool"
	"github.com/sergey-chebanov/fire/vast"
	"golang.org/x/time/rate"
)

type X struct {
	t   *testing.T
	URL string
}

func (x X) Run() error {
	t := x.t
	request := func() (err error) {
		res, err := http.Get(x.URL)
		if err != nil {
			t.Error(err)
			return
		}

		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			t.Error(err)
			return
		}
		_ = greeting
		//t.Logf("%s", greeting)
		return
	}
	return request()
}

func TestIt(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	pool := gopool.New(10, gopool.Config{CollectStat: true})

	for i := 0; i < 100; i++ {
		//go X{t}.Run()
		pool.Append(X{t, ts.URL})
	}

	pool.Close()
}

func TestSimpleRequests(t *testing.T) {
	sample, err := ioutil.ReadFile("vast/samples/empty_vast.xml")
	if err != nil {
		t.Error(err)
	}

	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sample)
	}))

	url, concurrency, N := ts.URL, 10, 100

	pool := gopool.New(concurrency, gopool.Config{CollectStat: true})

	//starting stat collecting
	const (
		increase = iota
		decrease
		keep
	)
	changeRate := make(chan int)
	go func() {
		for stat := range pool.Stat {
			if stat.Errors == 0 {
				changeRate <- increase
			}
			log.Println("Stat: ", stat)
		}
	}()

	//starting request dozer
	client, err := connect()
	if err != nil {
		t.Fatalf("%v: can't init client", err)
	}
	limiter := rate.NewLimiter(100, 1)
	for i := 0; i < N; i++ {

		//waiting limit to request
		if err := limiter.Wait(context.Background()); err != nil {
			log.Panic(err)
			break
		}
		pool.Append(gopool.TaskFunc(vast.MakeRequest(client, url, nil)))

		//check if we should increase rate
		select {
		case change := <-changeRate:
			if change == increase {
				limiter.SetLimit(limiter.Limit() + 10)
				log.Println("Limit: ", limiter.Limit())
			}
		default:
		}

	}

	pool.Close()
}
