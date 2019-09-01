package vast

import (
	"fmt"
	"log"
	"net/http"
)

type EventHandler struct{}

func (vh EventHandler) handle(res *http.Response) (err error) {
	log.Printf("Event Handler HTTP status = %d", res.StatusCode)
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("Event Handler got not 2xx HTTP status ")
	}
	return
}
