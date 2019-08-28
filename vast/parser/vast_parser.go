package parser

import (
	"encoding/xml"
	"fmt"
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
