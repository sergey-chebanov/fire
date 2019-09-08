package saver

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/kshvakov/clickhouse" //using via database/sql
	//_ "github.com/mailru/go-clickhouse" //using via database/sql
	"github.com/sergey-chebanov/fire/stat/record"
)

type clickhouseSaver struct {
	sync.WaitGroup
	connect *sql.DB
}

func init() {
	Constructors["clickhouse"] = newClickhouseSaver
}

func newClickhouseSaver(arguments string) (Interface, error) {
	ch := new(clickhouseSaver)

	//arguments = "http://127.0.0.1:9000?debug=true"

	connect, err := sql.Open("clickhouse", arguments)
	if err != nil {
		return nil, fmt.Errorf("init new clickhouse connection: %s", err)
	}
	ch.connect = connect

	// if err := ch.connect.Ping(); err != nil {
	// 	fmt.Println(err)
	// 	return nil, fmt.Errorf("ping failed: %s", err)
	// }

	_, err = ch.connect.Exec(`
			CREATE TABLE IF NOT EXISTS performance.tests 
			(
				sessionId Int64,
				started Int64,
				finished Int64,
				url String,
				err String
			) ENGINE = MergeTree()
			ORDER BY (started, finished)			
	`)

	if err != nil {
		return nil, fmt.Errorf("Can't create DB: %s", err)
	}

	return ch, nil
}

func (ch *clickhouseSaver) Save(recs []*record.Record) {
	ch.Add(1)
	go func(recs []*record.Record) {
		defer ch.Done()

		if len(recs) == 0 {
			log.Printf("Got empty records list -- skip")
			return
		}
		// log.Printf("Got records: %+v", recs)

		if ch == nil {
			log.Fatal("clickhouse connection is not inited")
		}

		var (
			tx, _   = ch.connect.Begin()
			stmt, _ = tx.Prepare(`
			INSERT INTO performance.tests (
				sessionId, started, finished, url, err
			) VALUES (
				?, ?, ?, ?, ?
			)`)
		)
		defer stmt.Close()

		for _, rec := range recs {
			// log.Printf("adding rec: %v", *rec)
			if _, err := stmt.Exec(
				rec.SessionID,
				rec.Start.UnixNano(),
				rec.Finish.UnixNano(),
				rec.URL,
				fmt.Sprint(rec.Err),
			); err != nil {
				log.Fatal(err)
			}

		}

		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

		log.Printf("saved %d", len(recs))
	}(recs)
}

func (ch *clickhouseSaver) Close() {
	ch.Wait()
	ch.connect.Close()
}
