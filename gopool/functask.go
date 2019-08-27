package gopool

//FuncTask helps to make and
type FuncTask struct {
	Task func() error
	Err  error
}

func (task *FuncTask) run() error {
	task.Err = task.Task()
	return task.Err
}
