package util

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Reminder is used for time consuming operations to remind user about progress.
type Reminder struct {
	start    time.Time     // start time since last warn
	interval time.Duration // interval to warn once
	level    logrus.Level  // log level to remind in general
}

// NewReminder returns a new Reminder instance.
//
// `level`: log level to remind in general.
//
// `interval`: interval to remind in warning level.
func NewReminder(level logrus.Level, interval time.Duration) *Reminder {
	if level < logrus.InfoLevel {
		panic("invalid log level to remind in general")
	}

	return &Reminder{
		start:    time.Now(),
		interval: interval,
		level:    level,
	}
}

// RemindWith reminds about specified `message` along with `key` and `value`.
func (reminder *Reminder) RemindWith(message string, key string, value interface{}) {
	reminder.Remind(message, logrus.Fields{key: value})
}

// Remind reminds about specified `message` and optional `fields`.
func (reminder *Reminder) Remind(message string, fields ...logrus.Fields) {
	if time.Since(reminder.start) > reminder.interval {
		reminder.remind(logrus.WarnLevel, message, fields...)
		reminder.start = time.Now()
	} else if logrus.IsLevelEnabled(reminder.level) {
		reminder.remind(reminder.level, message, fields...)
	}
}

func (reminder *Reminder) remind(level logrus.Level, message string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		logrus.WithFields(fields[0]).Log(level, message)
	} else {
		logrus.StandardLogger().Log(level, message)
	}
}
