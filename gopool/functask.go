package gopool

import (
	"github.com/sergey-chebanov/fire/stat/record"
)

//TaskFunc helps to make and
type TaskFunc struct {
	F        func() error
	TaskName string
}

//Run calls f()
func (f TaskFunc) Run() record.Record {
	return record.Record{Err: f.F(), Data: record.Fields{"name": f.TaskName}}
}
