package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	parser "github.com/sergey-chebanov/fire/vast/parser"
)

var p = fmt.Println

type firestarter struct {
	sync.WaitGroup
	ctx context.Context

	RPS    int
	incRPS int
	maxRPS int

	client *http.Client

	pollerNum int
	parserNum int

	statStatus chan string
	statErrors chan error
	/*statTimings chan struct {
		d    time.Duration
		what string
	}*/
}

func newFirestarter(ctx context.Context) *firestarter {
	f := firestarter{
		ctx: ctx,
		RPS: 10, maxRPS: 500, incRPS: 1,
		pollerNum: 500, parserNum: 10}

	f.statStatus = make(chan string)
	f.statErrors = make(chan error)
	f.stats()

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

	f.client = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

	return &f
}

func (f *firestarter) dozer(url string, out chan<- string) error {
	timer := time.Tick(1000 * time.Millisecond)

	go func() {
		f.Add(1)
		defer f.Done()
	Loop:
		for {
			for i := 0; i < f.RPS; i++ {
				select {
				case out <- url:
				case <-f.ctx.Done():
					return
				case <-timer:
					if i < f.RPS {
						f.RPS -= f.incRPS
						f.statErrors <- errors.New("Too fasts")
						//p("Too fast. Only ", i)
					}
					continue Loop
				}
			}
			<-timer
			f.RPS += f.incRPS
			if f.RPS > f.maxRPS {
				f.RPS = f.maxRPS
			}
			p("RPS: ", f.RPS)
		}
	}()

	return nil
}

func (f *firestarter) poller(urls <-chan string, bodies chan<- []byte) error {

	for i := 0; i < f.pollerNum; i++ {
		go func() {
			f.Add(1)
			defer f.Done()

			for {
				select {
				case <-f.ctx.Done():
					return
				case url := <-urls:
					req, err := http.NewRequest("GET", url, nil)
					resp, err := f.client.Do(req)
					if err != nil {
						f.statErrors <- err
						//p(err)
						resp = nil
						continue
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						f.statErrors <- err
						//p(err)
					}
					/*
						if len(body) < 3000 {
							p("corrupted VAST")
						}
					*/
					if bodies != nil {
						bodies <- body
					}
					f.statStatus <- resp.Status
				}
			}
		}()
	}
	return nil
}

func (f *firestarter) parseVast(vast <-chan []byte, eventsUrls chan<- string) error {
	for i := 0; i < f.parserNum; i++ {
		go func() {
			f.Add(1)
			defer f.Done()

			for {
				select {
				case <-f.ctx.Done():
					return
				case body := <-vast:
					vastP := parser.Parse(body)
					if vastP.Impression == "" {
						f.statErrors <- errors.New("Empty VAST")
						//p("Epmty vast")
					}
					eventsUrls <- vastP.Impression
					for _, event := range vastP.Creative[0].TrackingEvents {
						switch event.Type {
						case "start", "creativeView", "firstQuartile":
							eventsUrls <- event.URL
						}
					}
				}
			}
		}()
	}
	return nil
}

func (f *firestarter) stats() error {

	stats := make(map[string]int)
	//errors := make(map[string]int)
	errors := 0

	go func() {
		f.Add(1)
		defer f.Done()

		ticker := time.Tick(1000 * time.Millisecond)
		for {
			select {
			case <-f.ctx.Done():
			case <-f.statErrors:
				errors++
			case status := <-f.statStatus:
				stats[status]++
			case <-ticker:
				p(stats)
				p(errors)
			}
		}
	}()

	return nil
}

func main() {
	p(os.Args)
	url := flag.String("url", "https://ads.nsc-lab.io/ads/vast/get?out=3&opc=57baef95b53c4b5d6c57d5e8&cid=100&cnt=1&dur=0&adid=&appid=&site_id=171&zone_id=5871&pid=5848", "url to test")
	p(*url)
	f := newFirestarter(context.Background())
	urls := make(chan string, 10)
	f.dozer(*url, urls)
	bodies := make(chan []byte, 10)
	f.poller(urls, bodies)
	events := make(chan string, 10)
	f.parseVast(bodies, events)

	f.poller(events, nil)

	f.Wait()
}
