package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/sergey-chebanov/fire/gopool"
	"github.com/sergey-chebanov/fire/vast"
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
		Timeout:   1 * time.Second,
	}

}

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var url string
	flag.StringVar(&url, "url", "", "URL to test")
	var concurrency int
	flag.IntVar(&concurrency, "concurrency", 10, "Number of threads")
	var N int
	flag.IntVar(&N, "N", 0, "number of requests")
	var duration int
	flag.IntVar(&duration, "duration", 10, "test duration in seconds")
	var rateLimit int
	flag.IntVar(&rateLimit, "rate-limit", 100, "maximum RPM rate")

	flag.Parse()

	pool := gopool.New(concurrency, gopool.Config{CollectStat: true})

	//for non-blocking stat checking in the loop below we have to have blocking loop here
	statRepeater := make(chan gopool.Stat)
	go func() {
		for stat := range pool.Stat {
			log.Println("Stat: ", stat)
			statRepeater <- stat
		}
	}()

	//open connectio
	client := connect()

	eventURLs := make(chan vast.TimedURL)

	//start routine for appending new urls
	go func() {
		for timedURL := range eventURLs {
			pool.Append(gopool.TaskFunc(
				vast.MakeRequest(client, string(timedURL.URL), nil)))
		}
	}()

	//starting request dozer
	const limitStart = 2
	limiter := rate.NewLimiter(limitStart, 1)
	limitChange := float64(rateLimit-limitStart) / float64(duration) * 2

	timeIsUp := time.After(time.Duration(duration) * time.Second)

	for {
		//waiting limit to request
		if err := limiter.Wait(context.Background()); err != nil {
			log.Panic(err)
			break
		}
		pool.Append(gopool.TaskFunc(
			vast.MakeRequest(client, url, vast.Handler{URLsToAppend: eventURLs})))

		//check if we should increase rate
		select {
		case <-timeIsUp:
			log.Println("Time's Up")
			return
		case stat := <-statRepeater:
			log.Println("Limit: ", limiter.Limit())

			rateLimit := rate.Limit(rateLimit)

			newLimit := rate.Limit(float64(limiter.Limit()) + limitChange)

			if newLimit > rateLimit {
				newLimit = rateLimit
			}

			if newLimit < limitStart {
				newLimit = limitStart
			}

			if stat.Errors == 0 {
				limiter.SetLimit(newLimit)
			}
		default:
		}
	}
	pool.Close()
	time.Sleep(time.Second)
}
