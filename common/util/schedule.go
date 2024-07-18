package util

import (
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
