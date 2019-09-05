package gopool

import (
	"github.com/sergey-chebanov/fire/stat"
)

//TaskFunc helps to make and
type TaskFunc struct {
	F        func() error
	TaskName string
}

//Run calls f()
func (f TaskFunc) Run() stat.Record {
	return stat.Record{Err: f.F(), Data: stat.Fields{"name": f.TaskName}}
}
