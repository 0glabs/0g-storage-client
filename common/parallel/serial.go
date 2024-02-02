package parallel

import (
	"context"
	"sync"
)

func Serial(parallelizable Interface, tasks, routines, window int) error {
	if tasks == 0 {
		return nil
	}

	if routines == 0 {
		routines = 1
	}

	if routines > tasks {
		routines = tasks
	}

	channelLen := max(routines, window)
	if window > 0 && window < routines {
		window = routines
	}

	taskCh := make(chan int, channelLen)
	defer close(taskCh)
	resultCh := make(chan *Result, channelLen)
	defer close(resultCh)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	// start routines to do tasks
	for i := 0; i < routines; i++ {
		wg.Add(1)
		go work(ctx, i, parallelizable, taskCh, resultCh, &wg)
	}

	err := collect(parallelizable, taskCh, resultCh, tasks, channelLen, window > 0)

	// notify all routines to terminate
	cancel()

	// wait for termination for all routines
	wg.Wait()

	return err
}

func work(ctx context.Context, routine int, parallelizable Interface, taskCh <-chan int, resultCh chan<- *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-taskCh:
			val, err := parallelizable.ParallelDo(routine, task)
			resultCh <- &Result{routine, task, val, err}
			if err != nil {
				return
			}
		}
	}
}

func collect(parallelizable Interface, taskCh chan<- int, resultCh <-chan *Result, tasks, channelLen int, hasWindow bool) error {
	// if hasWindow = true, channelLen == window, fill window first
	// if hasWindow = false, channelLen = routines
	for i := 0; i < channelLen && i < tasks; i++ {
		taskCh <- i
	}

	var next, cnt int
	cache := map[int]*Result{}

	for result := range resultCh {
		if result.err != nil {
			return result.err
		}

		cache[result.Task] = result

		// handle task in sequence
		for cache[next] != nil {
			if err := parallelizable.ParallelCollect(cache[next]); err != nil {
				return err
			}

			// dispatch new task
			if hasWindow {
				if newTask := next + channelLen; newTask < tasks {
					taskCh <- newTask
				}
			}

			// clear cache and move window forward
			delete(cache, next)
			next++
		}
		if !hasWindow {
			if newTask := cnt + channelLen; newTask < tasks {
				taskCh <- newTask
			}
		}
		cnt += 1
		if next >= tasks {
			break
		}
	}

	return nil
}
