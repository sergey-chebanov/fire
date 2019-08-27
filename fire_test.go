package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergey-chebanov/fire/gopool"
)

var pool = gopool.New(1, gopool.Config{CollectStat: true})

var ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, client")
}))

func init() {

}

type X struct {
	t *testing.T
}

func (x X) Run() error {
	t := x.t
	request := func() (err error) {
		res, err := http.Get("http://ya.ru")
		if err != nil {
			t.Error(err)
			return
		}

		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			t.Error(err)
			return
		}

		t.Logf("%s", greeting)
		return
	}
	return request()
}

func TestIt(t *testing.T) {
	for i := 0; i < 100; i++ {
		pool.Append(X{t})
	}

	pool.Close()
}
