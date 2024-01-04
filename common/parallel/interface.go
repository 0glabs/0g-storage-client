package parallel

type Result struct {
	Routine int
	Task    int
	Value   interface{}
	err     error
}

type Interface interface {
	ParallelDo(routine, task int) (interface{}, error)
	ParallelCollect(result *Result) error
}
