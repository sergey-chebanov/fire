package gopool

import (
	"log"
	"testing"

	"github.com/sergey-chebanov/fire/stat"
	"github.com/sergey-chebanov/fire/stat/record"
	"github.com/sergey-chebanov/fire/stat/saver"
)

type TestTask struct {
	t *testing.T
}

func (task *TestTask) Run() record.Record {
	task.t.Log("added new task")
	return record.Record{Err: nil}
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func TestFuncTask0(t *testing.T) {

	pool := New(5, nil)

	test := &TestTask{t: t}

	for n := 0; n < 10; n++ {
		t.Log("added new task")

		pool.Append(test)
	}

	pool.Close()
}

func TestFuncTask1(t *testing.T) {

	saver, err := saver.New("clickhouse:http://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Panicf("Can't init saver %s", err)
	}
	pool := New(100, stat.New(saver))

	testData := make(chan int)

	go func() {
		var test []int
		for n := range testData {
			test = append(test, n)
		}
	}()

	x := func(N int) func() error {
		return func() error {
			t.Logf("N: %d\n", N)
			testData <- N
			return nil
		}
	}

	for n := 0; n < 20000; n++ {
		t.Log("added new task")
		pool.Append(TaskFunc{F: x(n)})
	}

	pool.Close()
	close(testData)
}
