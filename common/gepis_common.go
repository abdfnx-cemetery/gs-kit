package gepis_common

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type gepisLogger struct {
	name string
	logger *logrus.Entry
}

var gepisVersion string = "unknown"

func newgepisLogger(name string) *gepisLogger {
	newLogger := logrus.New()
	newLogger.SetOutput(os.Stdout)

	dl := &gepisLogger{
		name: name,
		logger: newLogger.WithFields(logrus.Fields{
			logFieldScope: name,
			logFieldType:  LogTypeLog,
		}),
	}

	dl.EnableJSONOutput(defaultJSONOutput)

	return dl
}

func (l *gepisLogger) EnableJSONOutput(enabled bool) {
	var formatter logrus.Formatter

	fieldMap := logrus.FieldMap{
		logrus.FieldKeyTime:  logFieldTimeStamp,
		logrus.FieldKeyLevel: logFieldLevel,
		logrus.FieldKeyMsg:   logFieldMessage,
	}

	hostname, _ := os.Hostname()
	l.logger.Data = logrus.Fields{
		logFieldScope:    l.logger.Data[logFieldScope],
		logFieldType:     LogTypeLog,
		logFieldInstance: hostname,
		logFieldgepisVer:  gepisVersion,
	}

	if enabled {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	}

	l.logger.Logger.SetFormatter(formatter)
}

func (l *gepisLogger) SetAppID(id string) {
	l.logger = l.logger.WithField(logFieldAppID, id)
}

func toLogrusLevel(lvl LogLevel) logrus.Level {
	l, _ := logrus.ParseLevel(string(lvl))
	return l
}

func (l *gepisLogger) SetOutputLevel(outputLevel LogLevel) {
	l.logger.Logger.SetLevel(toLogrusLevel(outputLevel))
}

func (l *gepisLogger) WithLogType(logType string) Logger {
	return &gepisLogger{
		name:   l.name,
		logger: l.logger.WithField(logFieldType, logType),
	}
}

func (l *gepisLogger) Info(args ...interface{}) {
	l.logger.Log(logrus.InfoLevel, args...)
}

func (l *gepisLogger) Infof(format string, args ...interface{}) {
	l.logger.Logf(logrus.InfoLevel, format, args...)
}

func (l *gepisLogger) Debug(args ...interface{}) {
	l.logger.Log(logrus.DebugLevel, args...)
}

func (l *gepisLogger) Debugf(format string, args ...interface{}) {
	l.logger.Logf(logrus.DebugLevel, format, args...)
}

func (l *gepisLogger) Warn(args ...interface{}) {
	l.logger.Log(logrus.WarnLevel, args...)
}

func (l *gepisLogger) Warnf(format string, args ...interface{}) {
	l.logger.Logf(logrus.WarnLevel, format, args...)
}

func (l *gepisLogger) Error(args ...interface{}) {
	l.logger.Log(logrus.ErrorLevel, args...)
}

func (l *gepisLogger) Errorf(format string, args ...interface{}) {
	l.logger.Logf(logrus.ErrorLevel, format, args...)
}

func (l *gepisLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *gepisLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
