package stat

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/sergey-chebanov/fire/stat/record"
	"github.com/sergey-chebanov/fire/stat/saver"
)

const waitTimeout = 10 //ms

//Collector collects atomic stats from workers with Collector.Collect
//and provide basic completion/error stat every second with Collector.Completed chan
type Collector interface {
	Collect(*record.Record)
	Close()
	Completed() <-chan map[error]int
}

//New constructs a new collector with provided saver interface.
func New(s saver.Interface) Collector {
	collector := &baseCollector{printTimings: true}
	collector.collectStats(s)
	return collector
}

type baseCollector struct {
	events chan *record.Record
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

func (bc *baseCollector) Collect(r *record.Record) {
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

func (bc *baseCollector) collectStats(saver saver.Interface) {

	bc.events = make(chan *record.Record)
	bc.completed = make(chan map[error]int, 10) //10 seconds buffer

	records := []*record.Record{}

	sendStat := func() {
		completed := map[error]int{}
		durations := make([]float64, 0, len(records))

		for _, rec := range records {
			completed[rec.Err]++
			if bc.printTimings {
				durations = append(durations, float64(rec.Duration()))
			}
		}

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

		if saver != nil {
			saver.Save(records)
		}
		records = []*record.Record{}

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

				records = append(records, stat)

				if stat.Err != nil {
					log.Printf("error: %s, %T", stat.Err, stat.Err)
				}

			case <-tick:
				statReady = true
			}

			if statReady {
				sendStat()
			}

			if stop {
				log.Println("stat receiver stopped; waiting for savers")
				if saver != nil {
					saver.Close()
				}
				log.Println("savers stopped")
				break
			}
		}
	}()

	return
}
