package vast

import (
	"fmt"
	"net/http"
)

//Task is an implementation of gopool.Task interfaces. It should be used for gopool.Pool as a Task
type Task struct {
	Client  *http.Client
	URL     string
	Handler ResponseHandler
}

//Run makes an HTTP request via the client and handle response with Handler
func (task Task) Run() (err error) {
	res, err := task.Client.Get(task.URL)
	if err != nil {
		return fmt.Errorf("%s: Request Failed", err)
	}
	defer res.Body.Close()

	if task.Handler != nil {
		err = task.Handler.handle(res)
	}

	return
}

func (task Task) ID() string {
	return task.URL
}
