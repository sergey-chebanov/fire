package vast

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type ResponseHandler interface {
	handle(res *http.Response) error
}

type Handler struct {
	URLsToAppend chan TimedURL
}

type TimedURL struct {
	URL string
	At  time.Time
}

func (vh Handler) handle(res *http.Response) (err error) {
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("%s: VAST Handler can't read body from request", err)
	}

	events := getEvents(body)

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("Got %d: VAST Handler got not 2xx HTTP status ", res.StatusCode)
	}

	const quartile = 10 * time.Second

	now := time.Now()
	if url, ok := events["start"]; ok {
		vh.URLsToAppend <- TimedURL{url, now}
		vh.URLsToAppend <- TimedURL{events["firstQuartile"], now.Add(quartile)}
		vh.URLsToAppend <- TimedURL{events["midpoint"], now.Add(2 * quartile)}
		vh.URLsToAppend <- TimedURL{events["thirdQuartile"], now.Add(3 * quartile)}
		vh.URLsToAppend <- TimedURL{events["complete"], now.Add(4 * quartile)}
	} else {
		return fmt.Errorf("got empty response (empty vast or not vast). body='%s'", body)
	}

	return
}

var eventsRegexp = regexp.MustCompile(`(?m)<Impression>\s*<\!\[CDATA\[(.*?)\]\]>\s*</Impression>` +
	`|<Tracking event="(.*?)">\s*<\!\[CDATA\[(.*?)\]\]>\s*</Tracking>` +
	`|<Error>(.*?)</Error>`)

func getEvents(xml []byte) map[string]string {

	events := map[string]string{}

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
	return events
}
