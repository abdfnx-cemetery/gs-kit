package gepis_common

import (
	"strings"
	"sync"
)

const (
	LogTypeLog = "log"
	LogTypeRequest = "request"
	logFieldTimeStamp = "time"
	logFieldLevel     = "level"
	logFieldType      = "type"
	logFieldScope     = "scope"
	logFieldMessage   = "msg"
	logFieldInstance  = "instance"
	logFieldGepisVer   = "ver"
	logFieldAppID     = "app_id"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel LogLevel = "info"
	WarnLevel LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	UndefinedLevel LogLevel = "undefined"
)

var (
	globalLoggers     = map[string]Logger{}
	globalLoggersLock = sync.RWMutex{}
)

type Logger interface {
	EnableJSONOutput(enabled bool)
	SetAppID(id string)
	SetOutputLevel(outputLevel LogLevel)
	WithLogType(logType string) Logger
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func toLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
		case "debug":
			return DebugLevel
		case "info":
			return InfoLevel
		case "warn":
			return WarnLevel
		case "error":
			return ErrorLevel
		case "fatal":
			return FatalLevel
	}

	return UndefinedLevel
}

func NewLogger(name string) Logger {
	globalLoggersLock.Lock()
	defer globalLoggersLock.Unlock()

	logger, ok := globalLoggers[name]
	if !ok {
		logger = newGepisLogger(name)
		globalLoggers[name] = logger
	}

	return logger
}

func getLoggers() map[string]Logger {
	globalLoggersLock.RLock()
	defer globalLoggersLock.RUnlock()

	l := map[string]Logger{}
	for k, v := range globalLoggers {
		l[k] = v
	}

	return l
}
