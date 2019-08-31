package parser

import (
	"encoding/xml"
	"fmt"
	"regexp"
)

//VAST used for unmarshaling VAST XML
type VAST struct {
	Error      string `xml:"Error"`
	Impression string `xml:"Ad>InLine>Impression"`
	Creative   []struct {
		ID             int `xml:"id,attr"`
		TrackingEvents []struct {
			Type string `xml:"event,attr"`
			URL  string `xml:",chardata"`
		} `xml:"Linear>TrackingEvents>Tracking"`
	} `xml:"Ad>InLine>Creatives>Creative"`
}

//Parse create a VAST from string
func Parse(data []byte) *VAST {

	v := VAST{}

	err := xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil
	}

	return &v
}

type EventsMap map[string]string

var re = regexp.MustCompile(`(?m)<Impression>\s*<\!\[CDATA\[(.*?)\]\]>\s*</Impression>` +
	`|<Tracking event="(.*?)">\s*<\!\[CDATA\[(.*?)\]\]>\s*</Tracking>` +
	`|<Error>(.*?)</Error>`)

func GetEvents(xml []byte) *EventsMap {

	events := EventsMap{}

	for _, event := range re.FindAllSubmatch(xml, -1) {
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
