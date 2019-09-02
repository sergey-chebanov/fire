package stat

import "log"

type KeyValue struct {
	Name  string
	Value string
}

type Collector interface {
	AddBulk([]KeyValue)
	Close()
}

func New(path string) Collector {
	return &sqliteCollector{}
}

type sqliteCollector struct {
	path string
}

func (col *sqliteCollector) AddBulk(KVs []KeyValue) {
	log.Println(KVs)
}

//Close wait till all events are saved and close the Collector
func (col *sqliteCollector) Close() {
}
