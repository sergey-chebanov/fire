package stat

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
)

const waitTimeout = 10 //ms

//Record is an atomic piece of information can be sent to the Collector. Can containt only Stringer
type Record struct {
	Err  error
	Data map[string]interface{}
}

type saver interface {
	save(Record)
}

type Collector interface {
	Collect(Record)
	Close()
	Completed() <-chan map[error]int
}

func New(path string) Collector {
	collector := &baseCollector{printTimings: true}
	collector.collectStats(&sqlite{})
	return collector
}

type baseCollector struct {
	events chan Record
	sync.WaitGroup
	completed    chan map[error]int
	printTimings bool
}

func (bc *baseCollector) Completed() <-chan map[error]int {
	return bc.completed
}

func (bc *baseCollector) Close() {
	close(bc.events)
	bc.Wait()
	close(bc.completed)
}

func (bc *baseCollector) Collect(r Record) {
	if bc.events == nil {
		return
	}

	//??? Or use buffered chan
	select {
	case bc.events <- r:
	case <-time.After(waitTimeout * time.Millisecond):
		log.Println("runStat dropped. Waiting too long")
	}

}

func (bc *baseCollector) collectStats(saver saver) {

	bc.events = make(chan Record)
	bc.completed = make(chan map[error]int, 10) //10 seconds buffer

	completed := map[error]int{}
	durations := []float64{}

	resetStat := func() {
		completed = map[error]int{}
		durations = []float64{}
	}

	sendStat := func() {
		if bc.printTimings && len(durations) > 0 {
			buf := &bytes.Buffer{}
			for _, p := range []float64{50, 75, 90, 99} {
				if average, err := stats.Percentile(durations, p); err == nil {
					fmt.Fprintf(buf, "%.2f%% -- %v; ", p, time.Duration(average))
					//log.Printf("Timing: %.2f%%- %f percentile", p, average)
				} else {
					log.Println(err)
				}
			}
			if buf.Len() > 0 {
				log.Printf("Timings: %s", buf)
			}
		}

		//try to send stat
		select {
		case bc.completed <- completed:
		default:
			log.Println("completed stat has been dropped")
		}
	}

	tick := time.Tick(time.Second)
	bc.Add(1)
	go func() {
		defer bc.Done()

		log.Println("stat receiver started")

		for {
			stop, statReady := false, false

			select {
			case stat, ok := <-bc.events:
				if !ok {
					statReady, stop = true, true
					break
				}

				//TODO: save to persistant DB
				saver.save(stat)

				completed[stat.Err]++

				if stat.Err != nil {
					log.Printf("error: %s, %T", stat.Err, stat.Err)
				}

				if bc.printTimings {
					if duration, err := stat.int("duration"); err == nil {
						durations = append(durations, float64(duration))
					}
				}

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

	return
}

type sqlite struct {
	path string
}

func (col *sqlite) save(r Record) {
	//log.Println(r.Err, r.Data)
}

func (r Record) Set(name string, stringer fmt.Stringer) {
	r.Data[name] = stringer
}

func (r Record) int(name string) (int, error) {
	if i, ok := r.Data[name].(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("Can't i{} cast to int")
}
