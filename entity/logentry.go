package entity

import "time"

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelError LogLevel = "ERROR"
	LogLevelPanic LogLevel = "PANIC"
	LogLevelFatal LogLevel = "FATAL"
	LogLevelDebug LogLevel = "DEBUG"
)

type LogEntry struct {
	Level     LogLevel
	Timestamp time.Time
	Source    string
	Message   string
}
