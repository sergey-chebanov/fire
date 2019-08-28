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
