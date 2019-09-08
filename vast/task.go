package vast

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/sergey-chebanov/fire/stat/record"
)

//Task is an implementation of gopool.Task interfaces. It should be used for gopool.Pool as a Task
type Task struct {
	Client  *http.Client
	URL     string
	Handler ResponseHandler
}

//Run makes an HTTP request via the client and handle response with Handler
func (task Task) Run() (rec *record.Record) {
	rec = record.New(nil)
	rec.URL = task.URL

	resp, err := task.Client.Get(task.URL)
	if err != nil {
		rec.Err = fmt.Errorf("%s: Request Failed", err)
		return
	}
	defer resp.Body.Close()

	if task.Handler != nil {
		rec.Err = task.Handler.handle(resp)
		rec.With("handler", reflect.TypeOf(task.Handler))
	}

	return
}
