package vast

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type TimedURL struct {
	url string
	at  time.Time
}

type Handler interface {
	handle(res *http.Response) error
}

type VastHandler struct {
	URLsToAppend chan TimedURL
}

func (vh VastHandler) handle(res *http.Response) (err error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	events := getEvents(body)

	if res.StatusCode%100 != 2 {
		return fmt.Errorf("VAST Handler got not 2xx HTTP status ")
	}

	const quartile = 10 * time.Second

	now := time.Now()
	if url, ok := (*events)["start"]; ok {
		vh.URLsToAppend <- TimedURL{url, now}
		vh.URLsToAppend <- TimedURL{(*events)["firstQuartile"], now.Add(quartile)}
		vh.URLsToAppend <- TimedURL{(*events)["midpoint"], now.Add(2 * quartile)}
		vh.URLsToAppend <- TimedURL{(*events)["thirdQuartile"], now.Add(3 * quartile)}
		vh.URLsToAppend <- TimedURL{(*events)["complete"], now.Add(4 * quartile)}
	} else {
		return errors.New("Got empty response (empty vast or not vast")
	}

	return
}

type EventHandler struct{}

func (vh EventHandler) handle(res *http.Response) (err error) {
	defer res.Body.Close()
	log.Println("Event Handler HTTP status = %d", res.StatusCode)
	if res.StatusCode%100 != 2 {
		return fmt.Errorf("Event Handler got not 2xx HTTP status ")
	}
	return nil
}

func MakeRequest(client *http.Client, url string, handler Handler) func() error {
	request := func() (err error) {
		res, err := client.Get(url)
		if err != nil {
			return
		}

		err = handler.handle(res)

		return
	}
	return request
}

type EventsMap map[string]string

var eventsRegexp = regexp.MustCompile(`(?m)<Impression>\s*<\!\[CDATA\[(.*?)\]\]>\s*</Impression>` +
	`|<Tracking event="(.*?)">\s*<\!\[CDATA\[(.*?)\]\]>\s*</Tracking>` +
	`|<Error>(.*?)</Error>`)

func getEvents(xml []byte) *EventsMap {

	events := EventsMap{}

	for _, event := range eventsRegexp.FindAllSubmatch(xml, -1) {
		switch {
		case len(event[1]) != 0:
			events["impression"] = string(event[1])
		case len(event[2]) != 0:
			events[string(event[2])] = string(event[3])
		case len(event[4]) != 0:
			events["error"] = string(event[4])
		}
	}
	return &events
}
