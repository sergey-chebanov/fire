package record

import (
	"fmt"
)

type fields map[interface{}]interface{}

//Record is an atomic piece of information can be sent to the Collector. Can containt only Stringer
type Record struct {
	Err  error
	Data fields
}

func New(err error) *Record {
	return &Record{Err: err}
}

func (r *Record) With(key interface{}, value interface{}) *Record {
	if r.Data == nil {
		r.Data = fields{}
	}

	r.Data[key] = value

	return r
}

func (r Record) Set(name string, stringer fmt.Stringer) {
	r.Data[name] = stringer
}

func (r *Record) Value(key interface{}) interface{} {
	if i, ok := r.Data[key]; ok {
		return i
	}
	return nil
}

func (r Record) Int(key interface{}) (int, error) {
	if i, ok := r.Data[key].(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("Can't i{} cast to int")
}

func (r Record) String(key interface{}) (string, error) {
	if i, ok := r.Data[key].(string); ok {
		return i, nil
	}
	return "", fmt.Errorf("Can't i{} cast to int")
}
