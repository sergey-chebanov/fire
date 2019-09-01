package gopool

//TaskFunc helps to make and
type TaskFunc struct {
	F        func() error
	TaskName string
}

//Run calls f()
func (f TaskFunc) Run() error {
	return f.F()
}

//ID get function name
func (f TaskFunc) ID() string {
	return f.TaskName
}
