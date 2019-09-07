package saver

import (
	"fmt"
	"testing"
	"time"

	"github.com/sergey-chebanov/fire/stat/record"
)

func Test_clickhouseSaver_save(t *testing.T) {
	ch, err := New("clickhouse:http://127.0.0.1:9000?debug=true")

	if err != nil {
		t.Errorf("can't init clickhouse connection: %s", err)
	}

	rec := record.Record{}
	rec.Err = fmt.Errorf("nothing")
	rec.Data = record.Fields{}
	rec.Data["sessionID"] = time.Now().UnixNano()
	rec.Data["started"] = time.Now().UnixNano()
	rec.Data["finished"] = time.Now().UnixNano() + 100
	rec.Data["url"] = "http://yandex.ru"

	ch.Save([]*record.Record{&rec, &rec, &rec})
}
