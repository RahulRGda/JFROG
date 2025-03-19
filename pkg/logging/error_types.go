package logging

type LoggingLevel string

const (
	INFO  LoggingLevel = "info"
	DEBUG              = "debug"
	WARN               = "warn"
	ERROR              = "error"
	FATAL              = "fatal"
)
