package util

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// Schedule runs the specified `action` periodically.
func Schedule(action func() error, interval time.Duration, errorMessage string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		if err := action(); err != nil {
			logrus.WithError(err).Warn(errorMessage)
		}
	}
}

// ScheduleNow runs the specified `action` immediately and periodically.
func ScheduleNow(action func() error, interval time.Duration, errorMessage string) {
	if err := action(); err != nil {
		logrus.WithError(err).Error(errorMessage)
	}

	Schedule(action, interval, errorMessage)
}

// WaitUntil runs the given function within a time duration
func WaitUntil(fn func() error, timeout time.Duration) error {
	ch := make(chan error, 1)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		ch <- fn()
	}()

	select {
	case <-ctxTimeout.Done():
		return fmt.Errorf("wait until task timeout")
	case err := <-ch:
		return err
	}
}
