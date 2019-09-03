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
	"github.com/sergey-chebanov/fire/stat"
	"github.com/sergey-chebanov/fire/vast"
	"golang.org/x/time/rate"
)

var client *http.Client

func init() {
	//starting request dozer
	var err error
	client, err = connect()
	if err != nil {
		log.Panicf("%v: can't init client", err)
	}
}

type X struct {
	t   *testing.T
	URL string
}

func (x X) ID() string {
	return x.URL
}

func (x X) Run() error {
	t := x.t
	request := func() (err error) {
		res, err := client.Get(x.URL)
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

	pool := gopool.New(10, nil)

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

	collector := stat.New("sqlite:fire.db")
	pool := gopool.New(concurrency, collector)

	//starting stat collecting
	const (
		increase = iota
		decrease
		keep
	)
	changeRate := make(chan int)
	go func() {
		for stat := range collector.Completed() {
			haveErrors := false
			for err := range stat {
				if err != nil {
					haveErrors = true
					break
				}
			}

			if !haveErrors {
				changeRate <- increase
			}
			log.Println("Stat: ", stat)
		}
	}()

	limiter := rate.NewLimiter(100, 1)
	for i := 0; i < N; i++ {

		//waiting limit to request
		if err := limiter.Wait(context.Background()); err != nil {
			log.Panic(err)
			break
		}
		pool.Append(vast.Task{
			Client:  client,
			URL:     url,
			Handler: nil,
		})

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
