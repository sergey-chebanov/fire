package gopool

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//Task is an abstract interface tasks for the Pool should comply to
type Task interface {
	Run() error
}

//Pool is a struct that holds everything needed for pool running
type Pool struct {
	sync.WaitGroup
	tasks chan Task
	//collecting stat
	completed chan runStat
	Stat      chan Stat
	statClose func()
}

//Config is a holder of setting for pool initization. See WithStat.
type Config struct {
	CollectStat bool
}

//Stat is a record of stats
type Stat struct {
	Completed       int
	AverageDuration time.Duration
}

func (pool *Pool) collectStats() {

	pool.completed = make(chan runStat)
	pool.Stat = make(chan Stat)

	var (
		completedStat int
		errorStat     int
		durationStat  []time.Duration
	)

	resetStat := func() {
		completedStat, errorStat = 0, 0
		durationStat = durationStat[:0]
	}

	sendStat := func() {
		average := time.Duration(0)
		for _, dur := range durationStat {
			average += dur
		}
		if len(durationStat) > 0 {
			average /= time.Duration(len(durationStat))
		}

		log.Println("Stat: ", average, completedStat)

		//try to send stat
		select {
		case pool.Stat <- Stat{completedStat, time.Duration(average)}:
		default:
		}
	}

	tick := time.Tick(time.Second)

	stopped := sync.WaitGroup{}
	stopped.Add(1)

	//stats
	go func() {
		defer stopped.Done()
		log.Println("stat receiver started")

		for {
			stop, statReady := false, false

			select {
			case stat, ok := <-pool.completed:
				if !ok {
					statReady, stop = true, true
					break
				}

				if stat.err == nil {
					completedStat++
				} else {
					fmt.Println(stat.err)
					errorStat++
				}

				durationStat = append(durationStat, stat.duration)

			case <-tick:
				statReady = true
			}

			if statReady {
				sendStat()
				resetStat()
			}

			if stop {
				log.Println("stat receiver stopped")
				break
			}
		}
	}()

	pool.statClose = func() {
		close(pool.completed)
		stopped.Wait()
	}

	return
}

type runStat struct {
	err      error
	duration time.Duration
}

//New creates and initialize a new pool.
func New(N int, config Config) *Pool {
	pool := &Pool{
		tasks: make(chan Task),
	}

	if config.CollectStat {
		pool.collectStats()
	}

	//TODO remove it -- just to test
	runAndMeasure := func(run func() error) (err error, dur time.Duration) {
		started := time.Now()
		err = run()
		dur = time.Since(started)
		return
	}

	//goroutines pool
	pool.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer pool.Done()
			for task := range pool.tasks {

				err, dur := runAndMeasure(task.Run)

				if err != nil {
					fmt.Println(err)
				}

				//should it be wrapped in go?
				select {
				case pool.completed <- runStat{err: err, duration: dur}:
				case <-time.After(10 * time.Millisecond):
					log.Println("runStat dropped. Waiting too long")
				}
			}
		}()
	}
	return pool
}

//Append schedule new task for running
func (pool *Pool) Append(task Task) {
	pool.tasks <- task
}

//Close stops the pool and free all resources
func (pool *Pool) Close() {
	close(pool.tasks)
	pool.Wait()

	if pool.statClose != nil {
		pool.statClose()
	}
}
