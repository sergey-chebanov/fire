package saver

import (
	"fmt"
	"testing"
	"time"

	"github.com/sergey-chebanov/fire/stat/record"
)

func Test_clickhouseSaver_save(t *testing.T) {
	ch, err := New("clickhouse", "http://127.0.0.1:9000")

	if err != nil {
		t.Errorf("can't init clickhouse connection: %s", err)
	}

	rec := record.New(fmt.Errorf("nothing"))
	rec.SessionID = 1234
	rec.Start = time.Now()
	rec.Finish = time.Now()
	rec.URL = "http://yandex.ru"

	ch.Save([]*record.Record{rec, rec, rec})

	ch.Close()
}
