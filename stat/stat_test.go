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

		rec1 := &record.Record{}
		rec1.Data = record.Fields{}
		rec1.Data["sessionID"] = 588
		rec1.Data["started"] = time.Now().UnixNano()
		rec1.Data["finished"] = time.Now().UnixNano() + 100
		rec1.Data["url"] = "http://yandex.ru"

		collector.Collect(rec1)

		rec2 := &record.Record{}
		rec2.Data = record.Fields{}
		rec2.Data["sessionID"] = 588
		rec2.Data["started"] = time.Now().UnixNano()
		rec2.Data["finished"] = time.Now().UnixNano() + 100
		rec2.Data["url"] = "http://yandex.ru"
		collector.Collect(rec2)

		rec3 := &record.Record{}
		rec3.Data = record.Fields{}
		rec3.Data["sessionID"] = 588
		rec3.Data["started"] = time.Now().UnixNano()
		rec3.Data["finished"] = time.Now().UnixNano() + 100
		rec3.Data["url"] = "http://yandex.ru"
		rec3.Err = fmt.Errorf("nothing")
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
