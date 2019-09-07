package record

import "fmt"

type Fields map[string]interface{}

//Record is an atomic piece of information can be sent to the Collector. Can containt only Stringer
type Record struct {
	Err  error
	Data Fields
}

func (r *Record) With(key string, value interface{}) {
	if r.Data == nil {
		r.Data = Fields{}
	}

	r.Data[key] = value
}

func (r Record) Set(name string, stringer fmt.Stringer) {
	r.Data[name] = stringer
}

func (r Record) Int(name string) (int, error) {
	if i, ok := r.Data[name].(int); ok {
		return i, nil
	}
	return 0, fmt.Errorf("Can't i{} cast to int")
}
