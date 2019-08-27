package gopool

//TaskFunc helps to make and
type TaskFunc func() error

//Run calls f()
func (f TaskFunc) Run() error {
	return f()
}
