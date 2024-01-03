package common

import (
	"io"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// The diveLogger type is a struct that contains a logrus.Logger instance.
// log - The `log` is a pointer to an instance of the `logrus.Logger` struct. This
// logger is used for logging messages and events in the `diveLogger` struct.
type diveLogger struct {
	log *logrus.Logger
}

// The function `NewDiveLogger` creates a new instance of a logger that logs information and errors to
// separate files.
func NewDiveLogger(infoFilePath string, errorFilePath string) *diveLogger {

	log := logrus.New()

	log.SetOutput(io.Discard)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		PadLevelText:    true,
	})

	ditLogger := &lumberjack.Logger{
		// Log file absolute path, os agnostic
		Filename:  filepath.ToSlash(infoFilePath),
		LocalTime: true,
	}

	// Fork writing into two outputs
	ditWriter := io.MultiWriter(ditLogger)

	errorLogger := &lumberjack.Logger{
		Filename:  filepath.ToSlash(errorFilePath),
		LocalTime: true,
	}

	// Fork writing into two outputs
	errorWriter := io.MultiWriter(errorLogger)

	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  ditWriter,
			logrus.DebugLevel: ditWriter,
			logrus.TraceLevel: ditWriter,
			logrus.WarnLevel:  ditWriter,
			logrus.ErrorLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
		},
		&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	))

	return &diveLogger{log: log}
}

// The `SetErrorToStderr()` function sets the output of the logger to `os.Stderr`. This means that any
// log messages or errors logged using this logger will be printed to the standard error output.
func (d *diveLogger) SetErrorToStderr() {
	d.log.SetOutput(os.Stderr)
}

// The `SetOutputToStdout()` function sets the output of the logger to `os.Stdout`. This means that any
// log messages or errors logged using this logger will be printed to the standard output.
func (d *diveLogger) SetOutputToStdout() {
	d.log.SetOutput(os.Stdout)
}

// The `logWithFields` function is a helper function in the `diveLogger` struct. It is used to log
// messages with additional fields.
func (d *diveLogger) logWithFields(level logrus.Level, kind string, format string, args ...interface{}) {
	if d.log.IsLevelEnabled(level) {
		formatWithKind := " " + kind + "  " + format
		d.log.Logf(level, formatWithKind, args...)
	}
}

// The `Debug` function is a method of the `diveLogger` struct. It is used to log a debug message with
// the specified message string. It calls the `logWithFields` function, passing the log level
// `logrus.DebugLevel`, the kind "🐞 debug", and the message string as arguments. The `logWithFields`
// function adds the log level and kind as fields to the log entry and logs the message using the
// logrus logger.
func (d *diveLogger) Debug(message string) {
	d.logWithFields(logrus.DebugLevel, "🐞 debug", message)
}

// The `Info` function is a method of the `diveLogger` struct. It is used to log an informational
// message with the specified message string.
func (d *diveLogger) Info(message string) {
	d.logWithFields(logrus.InfoLevel, "ℹ️", message)
}

// The `Warn` function is a method of the `diveLogger` struct. It is used to log a warning message with
// the specified message string.
func (d *diveLogger) Warn(message string) {
	d.logWithFields(logrus.WarnLevel, "⚠️", message)
}

// The `Error` function is a method of the `diveLogger` struct. It is used to log an error message with
// the specified error code and error message.
func (d *diveLogger) Error(errorCode ErrorCode, errorMessage string) {
	d.logWithFields(logrus.ErrorLevel, "🛑", "Code:%d Error: %s", errorCode, errorMessage)
}

// The `Fatal` function is a method of the `diveLogger` struct. It is used to log an error message with
// the specified error code and error message, and then exit the program with a status code of 1.
func (d *diveLogger) Fatal(errorCode ErrorCode, errorMessage string) {
	d.logWithFields(logrus.FatalLevel, "💀", "Code:%d Error: %s", errorCode, errorMessage)
	d.log.Exit(1)
}

// The `Infof` function is a method of the `diveLogger` struct. It is used to log an informational
// message with the specified format string and arguments.
func (d *diveLogger) Infof(format string, args ...interface{}) {
	d.logWithFields(logrus.InfoLevel, "ℹ️", format, args...)
}

// The `Warnf` function is a method of the `diveLogger` struct. It is used to log a warning message
// with the specified format string and arguments. It calls the `logWithFields` function, passing the
// log level `logrus.WarnLevel`, the kind "⚠️", and the format string and arguments as arguments. The
// `logWithFields` function adds the log level and kind as fields to the log entry and logs the
// formatted message using the logrus logger.
func (d *diveLogger) Warnf(format string, args ...interface{}) {
	d.logWithFields(logrus.WarnLevel, "⚠️", format, args...)
}

// The `Debugf` function is a method of the `diveLogger` struct. It is used to log a debug message with
// the specified format string and arguments. It calls the `logWithFields` function, passing the log
// level `logrus.DebugLevel`, the kind "🐞", and the format string and arguments as arguments. The
// `logWithFields` function adds the log level and kind as fields to the log entry and logs the
// formatted message using the logrus logger.
func (d *diveLogger) Debugf(format string, args ...interface{}) {
	d.logWithFields(logrus.DebugLevel, "🐞", format, args...)
}

// The `Errorf` function is a method of the `diveLogger` struct. It is used to log an error message
// with the specified error code, format string, and arguments.
func (d *diveLogger) Errorf(errorCode ErrorCode, format string, args ...interface{}) {
	d.logWithFields(logrus.ErrorLevel, "🛑", format, args...)

}

// The `Fatalf` function is a method of the `diveLogger` struct. It is used to log an error message
// with the specified error code, format string, and arguments. It calls the `logWithFields` function,
// passing the log level `logrus.FatalLevel`, the kind "💀", and the format string and arguments as
// arguments. The `logWithFields` function adds the log level and kind as fields to the log entry and
// logs the formatted message using the logrus logger. After logging the message, it calls the `Exit`
// function of the logrus logger with a status code of 1, which causes the program to exit.
func (d *diveLogger) Fatalf(errorCode ErrorCode, format string, args ...interface{}) {
	d.logWithFields(logrus.FatalLevel, "💀", format, args...)
	d.log.Exit(1)
}
