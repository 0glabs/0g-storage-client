package parallel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	t      *testing.T
	result []int
}

func (f *foo) ParallelDo(routine, task int) (interface{}, error) {
	return task * task, nil
}

func (f *foo) ParallelCollect(result *Result) error {
	assert.Nil(f.t, result.err)
	assert.Equal(f.t, len(f.result), result.Task)
	assert.Equal(f.t, result.Task*result.Task, result.Value.(int))

	f.result = append(f.result, result.Value.(int))

	return nil
}

func TestSerial(t *testing.T) {
	f := foo{t, nil}

	tasks := 100

	err := Serial(&f, tasks, 4, 16)
	assert.Nil(t, err)
	assert.Equal(t, tasks, len(f.result))

	for i := 0; i < tasks; i++ {
		assert.Equal(t, i*i, f.result[i])
	}
}
