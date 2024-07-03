package parallel

import "context"

type Result struct {
	Routine int
	Task    int
	Value   interface{}
	err     error
}

type Interface interface {
	ParallelDo(ctx context.Context, routine, task int) (interface{}, error)
	ParallelCollect(result *Result) error
}
