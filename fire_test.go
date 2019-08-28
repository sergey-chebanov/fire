package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergey-chebanov/fire/gopool"
)

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
		res, err := http.Get(ts.URL)
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
		_ = greeting
		//t.Logf("%s", greeting)
		return
	}
	return request()
}

var pool = gopool.New(10, gopool.Config{CollectStat: true})

func TestIt(t *testing.T) {
	//time.Sleep(1000 * time.Millisecond)

	for i := 0; i < 100; i++ {
		//go X{t}.Run()
		pool.Append(X{t})
	}

	pool.Close()
}
