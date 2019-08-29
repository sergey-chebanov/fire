package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"golang.org/x/time/rate"

	"github.com/sergey-chebanov/fire/gopool"
)

func parseFlags() (url string, concurrency int, N int, err error) {
	flag.StringVar(&url, "url", "", "URL to test")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of threads")
	flag.IntVar(&N, "N", 100, "number of requests")

	flag.Parse()

	return
}

func connect() *http.Client {
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	tr := http.DefaultTransport.(*http.Transport)
	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == "ads.nsc-lab.io:443" {
			addr = "85.235.174.22:443"
		}
		if addr == "ads.nsc-lab.io:80" {
			addr = "85.235.174.22:80"
		}
		return dialer.DialContext(ctx, network, addr)
	}
	tr.MaxIdleConns = 500
	tr.MaxIdleConnsPerHost = 500

	return &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

}

func main() {

	var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	url, concurrency, N, err := parseFlags()

	url = ts.URL

	if err != nil {
		fmt.Println(err)
		return
	}

	pool := gopool.New(concurrency, gopool.Config{CollectStat: true})

	client := connect()

	request := func(url string) func() (err error) {
		requestInternal := func() (err error) {
			res, err := client.Get(url)
			if err != nil {
				return
			}

			body, err := ioutil.ReadAll(res.Body)

			//fmt.Printf("%s\n", body[:30])
			_ = body

			res.Body.Close()

			if err != nil {
				return
			}
			return
		}
		return requestInternal
	}

	const (
		increase = iota
		decrease
		keep
	)

	limiter := rate.NewLimiter(100, 1)
	changeRate := make(chan int)

	go func() {
		for stat := range pool.Stat {
			if stat.Errors == 0 {
				changeRate <- increase
			}
			log.Println("Stat: ", stat)
		}
	}()

	for i := 0; i < N; i++ {
		if err := limiter.Wait(context.Background()); err != nil {
			log.Panic(err)
			break
		}
		pool.Append(gopool.TaskFunc(request(url)))

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

	fmt.Println(url)
}
