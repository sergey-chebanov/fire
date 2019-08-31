package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sergey-chebanov/fire/vast/parser"
)

type timedURL struct {
	url string
	at  time.Time
}

var urlsToAppend = make(chan timedURL)

func vastHanlder(res *http.Response) (err error) {
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	//log.Printf("body: %s \n--------\n", body)

	if err != nil {
		return
	}

	events := parser.GetEvents(body)

	//log.Println("parsed ", events)

	if url, ok := (*events)["start"]; ok {
		urlsToAppend <- url
	} else {
		return errors.New("Got empty response (empty vast or not vast")
	}
	return
}

func eventHandler(res *http.Response) (err error) {
	defer res.Body.Close()
	return nil
}
