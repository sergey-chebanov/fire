package gopool

import (
	"log"
	"sync"
	"time"

	"github.com/sergey-chebanov/fire/stat"
)

//Task is an abstract interface tasks for the Pool should comply to
type Task interface {
	Run() stat.Record
}

//Pool is a struct that holds everything needed for pool running
type Pool struct {
	sync.WaitGroup
	tasks chan Task
	//collecting stat
	collector stat.Collector
}

//New creates and initialize a new pool.
func New(N int, collector stat.Collector) *Pool {
	pool := &Pool{
		tasks:     make(chan Task),
		collector: collector,
	}

	runAndMeasure := func(run func() stat.Record) stat.Record {
		started := time.Now()
		stat := run()
		stat.Data["duration"] = int(time.Since(started))
		return stat
	}

	id := func(run func() stat.Record) stat.Record {
		return run()
	}

	//goroutines pool
	pool.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer pool.Done()
			for task := range pool.tasks {

				runAndMeasure := runAndMeasure

				if collector == nil {
					runAndMeasure = id
				}

				stat := runAndMeasure(task.Run)

				if collector != nil {
					collector.Collect(stat)
				}
			}
		}()
	}
	return pool
}

//Append schedule new task for running
func (pool *Pool) Append(task Task) {
	const t = 100
	select {
	case pool.tasks <- task:
	case <-time.After(t * time.Millisecond):
		log.Printf("All workers are busy. Can't append in expected time(%d ms). Task has been dropped", t)
	}
}

//Close stops the pool and free all resources
func (pool *Pool) Close() {
	close(pool.tasks)
	pool.Wait()

	if pool.collector != nil {
		pool.collector.Close()
	}
}
