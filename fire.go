package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/time/rate"

	"github.com/sergey-chebanov/fire/gopool"
	"github.com/sergey-chebanov/fire/vast"
)

func connect() (client *http.Client, err error) {
	dialer := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	tr := http.DefaultTransport.(*http.Transport)
	if err = http2.ConfigureTransport(tr); err != nil {
		log.Printf("%v: Error while initilizing http2 client", err)
		return
	}
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

	client = &http.Client{
		Transport: tr,
		Timeout:   1 * time.Second,
	}
	return
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
			statRepeater <- stat
		}
	}()

	//open connectio
	client, err := connect()
	if err != nil {
		log.Panicf("%v: can't init client", err)
	}

	eventURLs := make(chan vast.TimedURL)

	//start routine for appending new urls
	go func() {
		for timedURL := range eventURLs {
			pool.Append(vast.Task{
				Client:  client,
				URL:     string(timedURL.URL),
				Handler: nil})
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
		pool.Append(vast.Task{
			Client:  client,
			URL:     url,
			Handler: vast.Handler{URLsToAppend: eventURLs}})

		//check if we should increase rate
		select {
		case <-timeIsUp:
			log.Println("Time's Up")
			return
		case stat := <-statRepeater:
			log.Printf("VAST Requests Limit: %.2f RPS -- Stat: %v", limiter.Limit(), stat)

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
