package parser

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParseGood(t *testing.T) {
	data, err := ioutil.ReadFile("vast.xml")
	if err != nil {
		t.Error("Could not read the file with sample vast")
	}
	v := Parse(data)
	t.Logf("Impression: %v\n", strings.TrimSpace(v.Impression))
	for _, tracking := range v.Creative[0].TrackingEvents {
		t.Logf("TrackingEvents: %v\n", strings.TrimSpace(tracking.URL))
	}

	t.Logf("Creative ID: %v\n", v.Creative[0].ID)
}

func TestParseErr(t *testing.T) {
	data, err := ioutil.ReadFile("empty_vast.xml")
	if err != nil {
		t.Error("Could not read the file with sample vast")
	}
	v := Parse(data)
	t.Logf("Impression: %v\n", v)
	if v.Error != "blah" {
		t.Error("not correct")
	}
}

func TestGetEvents(t *testing.T) {
	xmlSample, err := ioutil.ReadFile("vast.xml")
	if err != nil {
		t.Error(err)
	}

	events := GetEvents(xmlSample)

	for event, url := range *events {
		t.Logf("%s -- %s", event, url)
	}

	if _, ok := (*events)["impression"]; !ok {
		t.Error("\"Impression\" event must be included in events")
	}

	if _, ok := (*events)["start"]; !ok {
		t.Error("\"start\" event must be included in events")
	}

}

func TestGetEventsEmpty(t *testing.T) {
	xmlSample, err := ioutil.ReadFile("empty_vast.xml")
	if err != nil {
		t.Error(err)
	}

	events := GetEvents(xmlSample)

	for event, url := range *events {
		t.Logf("%s -- %s", event, url)
	}

	if _, ok := (*events)["error"]; !ok {
		t.Error("Events must contain \"error\" event")
	}

	if len(*events) != 1 {
		t.Error("Events must contain the only element (event)")
	}
}
