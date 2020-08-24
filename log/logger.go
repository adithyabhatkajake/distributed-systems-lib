package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var (
	// std is the name of the standard logger in stdlib `log`
	logger = logrus.New()
)

// Formatter renames logrus Formatter
type Formatter logrus.Formatter

// Level refers to the LogLevel
// It is one of PanicLevel, FatalLevel, ErrorLevel,
// WarnLevel, InfoLevel, DebugLevel, TraceLevel
type Level uint32

const (
	// PanicLevel errorlevel. Logs.
	// Use log.SetLevel(log.PanicLevel) to set it.
	PanicLevel = logrus.PanicLevel
	// FatalLevel errorlevel. Logs.
	// Use log.SetLevel(log.FatalLevel) to set it.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel errorlevel. Logs.
	// Use log.SetLevel(log.ErrorLevel) to set it.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel errorlevel. Logs.
	// Use log.SetLevel(log.WarnLevel) to set it.
	WarnLevel = logrus.WarnLevel
	// InfoLevel errorlevel. Logs.
	// Use log.SetLevel(log.InfoLevel) to set it.
	InfoLevel = logrus.InfoLevel
	// DebugLevel errorlevel. Logs.
	// Use log.SetLevel(log.DebugLevel) to set it.
	DebugLevel = logrus.DebugLevel
	// TraceLevel errorlevel. Logs.
	// Use log.SetLevel(log.TraceLevel) to set it.
	TraceLevel = logrus.TraceLevel
)

// SetOutput sets the standard logger output.
func SetOutput(out io.Writer) {
	logger.SetOutput(out)
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter Formatter) {
	logger.SetFormatter(formatter)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	logger.SetReportCaller(include)
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	return logger.GetLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level logrus.Level) bool {
	return logger.IsLevelEnabled(level)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	logger.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	logger.Warning(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	logger.Traceln(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	logger.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}
