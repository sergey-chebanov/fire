package stat

import (
	"fmt"
	"log"
	"testing"
)

func TestNew(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	collector := New("sqlite:blahminor.db")
	defer collector.Close()

	collector.Collect(Record{Err: nil, Data: map[string]interface{}{"duration": 1, "id": "123123 12312312312"}})
	collector.Collect(Record{Err: nil, Data: map[string]interface{}{"duration": 2, "id": "123123 12312312312"}})
	collector.Collect(Record{Err: fmt.Errorf("oi"), Data: map[string]interface{}{"duration": 3, "id": "123123 12312312312"}})

	s, ok := <-collector.Completed()
	if !(ok && s[nil] == 2 && len(s) == 2) {
		t.Error("We must get 2 different events and 2 nil errors particulary")
	}

}
