package stat

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/sergey-chebanov/fire/stat/record"
	"github.com/sergey-chebanov/fire/stat/saver"
)

func TestNew(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if saver, err := saver.New("clickhouse:http://127.0.0.1:9000?debug=true"); err == nil {

		collector := New(saver)
		defer collector.Close()

		//rec1 := &record.Record{}
		rec1 := record.New(nil).
			With("sessionID", 588).
			With("started", time.Now().UnixNano()).
			With("finished", time.Now().UnixNano()+100).
			With("url", "http://yandex.ru")

		collector.Collect(rec1)

		rec2 := record.New(nil).
			With("sessionID", 588).
			With("started", time.Now().UnixNano()).
			With("finished", time.Now().UnixNano()+100).
			With("url", "http://yandex.ru")

		collector.Collect(rec2)

		rec3 := record.New(fmt.Errorf("Test error")).
			With("sessionID", 588).
			With("started", time.Now().UnixNano()).
			With("finished", time.Now().UnixNano()+100).
			With("url", "http://yandex.ru")
		collector.Collect(rec3)

		s, ok := <-collector.Completed()
		if !(ok && s[nil] == 2 && len(s) == 2) {
			t.Error("We must get 2 different events and 2 nil errors particulary")
		}
	} else {
		t.Errorf("Clickhouse init failed: %s", err)
		return
	}
}
