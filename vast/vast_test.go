package vast

import (
	"io/ioutil"
	"testing"
)

func TestGetEvents(t *testing.T) {
	xmlSample, err := ioutil.ReadFile("samples/vast.xml")
	if err != nil {
		t.Error(err)
	}

	events := getEvents(xmlSample)

	for event, url := range events {
		t.Logf("%s -- %s", event, url)
	}

	if _, ok := events["impression"]; !ok {
		t.Error("\"Impression\" event must be included in events")
	}

	if _, ok := events["start"]; !ok {
		t.Error("\"start\" event must be included in events")
	}
}

func TestGetEventsEmpty(t *testing.T) {
	xmlSample, err := ioutil.ReadFile("samples/empty_vast.xml")
	if err != nil {
		t.Error(err)
	}

	events := getEvents(xmlSample)

	for event, url := range events {
		t.Logf("%s -- %s", event, url)
	}

	if _, ok := events["error"]; !ok {
		t.Error("Events must contain \"error\" event")
	}

	if len(events) != 1 {
		t.Error("Events must contain the only element (event)")
	}
}
