package gopool

import (
	"log"
	"testing"
)

type TestTask struct {
	t *testing.T
}

func (task *TestTask) Run() error {
	task.t.Log("added new task")
	return nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func TestFuncTask0(t *testing.T) {

	pool := New(5, Config{CollectStat: true})

	test := &TestTask{t: t}

	for n := 0; n < 10; n++ {
		t.Log("added new task")

		pool.Append(test)
	}

	pool.Close()
}

func TestFuncTask1(t *testing.T) {

	pool := New(100, Config{CollectStat: true})

	testData := make(chan int)

	go func() {
		var test []int
		for n := range testData {
			test = append(test, n)
		}
		t.Logf("%v\n", test)
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
		pool.Append(TaskFunc(x(n)))
	}

	pool.Close()
	close(testData)
}
