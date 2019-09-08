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

	if saver, err := saver.New("clickhouse", "http://127.0.0.1:8123"); err == nil {

		collector := New(saver)

		//rec1 := &record.Record{}
		rec1 := record.New(nil)
		rec1.Start = time.Now()
		rec1.SessionID = 590
		rec1.Finish = time.Now().Add(100)
		rec1.URL = "http://yandex.ru"

		collector.Collect(rec1)

		rec2 := record.New(nil)
		rec2.SessionID = 590
		rec2.Start = time.Now()
		rec2.Finish = time.Now().Add(100)
		rec2.URL = "http://yandex.ru"

		collector.Collect(rec2)

		rec3 := record.New(fmt.Errorf("Test error"))
		rec3.SessionID = 590
		rec3.Start = time.Now()
		rec3.Finish = time.Now().Add(100)
		rec3.URL = "http://yandex.ru"
		collector.Collect(rec3)

		s, ok := <-collector.Completed()
		if !(ok && s[nil] == 2 && len(s) == 2) {
			t.Error("We must get 2 different events and 2 nil errors particulary")
		}

		collector.Close()
	} else {
		t.Errorf("Clickhouse init failed: %s", err)
		return
	}
}
