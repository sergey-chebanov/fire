package record

import (
	"fmt"
	"time"
)

type fields map[interface{}]interface{}

//Record is an atomic piece of information can be sent to the Collector. Can containt only Stringer
type Record struct {
	Err       error
	SessionID int64
	Start     time.Time
	Finish    time.Time
	URL       string
	data      fields
}

func New(err error) *Record {
	record := &Record{Err: err}
	return record
}

func (r *Record) With(key interface{}, value interface{}) *Record {
	if r.data == nil {
		r.data = fields{}
	}

	r.data[key] = value

	return r
}

func (r *Record) Value(key interface{}) interface{} {
	if i, ok := r.data[key]; ok {
		return i
	}
	return nil
}

func (r *Record) Int(key interface{}) (int, error) {
	if i, ok := r.data[key].(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("Can't i{} cast to int")
}

func (r *Record) String(key interface{}) (string, error) {
	if i, ok := r.data[key].(string); ok {
		return i, nil
	}
	return "", fmt.Errorf("Can't i{} cast to int")
}

func (r *Record) Duration() time.Duration {
	return r.Finish.Sub(r.Start)
}
