package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sergey-chebanov/fire/gopool"
)

func parseFlags() (url string, concurrency int, N int, err error) {
	flag.StringVar(&url, "url", "", "URL to test")
	flag.IntVar(&concurrency, "concurrency", 10, "Number of threads")
	flag.IntVar(&N, "N", 100, "number of requests")

	flag.Parse()

	return
}

func main() {

	url, concurrency, N, err := parseFlags()

	if err != nil {
		fmt.Println(err)
		return
	}

	pool := gopool.New(concurrency, gopool.Config{CollectStat: true})

	request := func(url string) func() (err error) {
		requestInternal := func() (err error) {
			res, err := http.Get(url)
			if err != nil {
				return
			}

			greeting, err := ioutil.ReadAll(res.Body)
			_ = greeting

			res.Body.Close()

			if err != nil {
				return
			}
			return
		}
		return requestInternal
	}

	for i := 0; i < N; i++ {
		pool.Append(gopool.TaskFunc(request(url)))
	}

	pool.Close()

	fmt.Println(url)
}
