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
	err  error
	data map[string]interface{}
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

	select {
	case bc.events <- r:
	case <-time.After(waitTimeout * time.Millisecond):
		log.Println("runStat dropped. Waiting too long")
	}

}

func (bc *baseCollector) collectStats(saver saver) {

	bc.events = make(chan Record)
	bc.completed = make(chan map[error]int)

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
					fmt.Fprintf(buf, "Timing: %.2f%% -- %f percentile", p, average)
					//log.Printf("Timing: %.2f%%- %f percentile", p, average)
				} else {
					log.Println(err)
				}
			}
			log.Print(buf)
		}

		//log.Println("Stat: ", time.Duration(average), completedStat)

		//try to send stat
		select {
		case bc.completed <- completed:
		case <-time.After(waitTimeout * time.Millisecond):
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

				completed[stat.err]++

				if stat.err != nil {
					log.Printf("error: %s, %T", stat.err, stat.err)
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
	log.Println(r.err, r.data)
}

func (r Record) Set(name string, stringer fmt.Stringer) {
	r.data[name] = stringer
}

func (r Record) int(name string) (int, error) {
	if i, ok := r.data[name].(int); ok {
		return i, nil
	}
	return 0, nil
}
