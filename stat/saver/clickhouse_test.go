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

	rec := record.New(fmt.Errorf("nothing")).
		With("sessionID", time.Now().UnixNano()).
		With("started", time.Now().UnixNano()).
		With("finished", time.Now().UnixNano()).
		With("sessionID", time.Now().UnixNano()).
		With("url", "http://yandex.ru")

	ch.Save([]*record.Record{rec, rec, rec})
}
