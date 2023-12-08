package common

import (
	"io"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type diveLogger struct {
	log *logrus.Logger
}

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

func (d *diveLogger) SetErrorToStderr() {
	d.log.SetOutput(os.Stderr)
}
func (d *diveLogger) SetOutputToStdout() {
	d.log.SetOutput(os.Stdout)
}

func (d *diveLogger) logWithFields(level logrus.Level, kind string, format string, args ...interface{}) {
	if d.log.IsLevelEnabled(level) {
		d.log.WithFields(logrus.Fields{level.String(): kind}).Logf(level, format, args...)
	}
}

func (d *diveLogger) Debug(message string) {
	d.logWithFields(logrus.DebugLevel, "🐞 debug", message)
}

func (d *diveLogger) Info(message string) {
	d.logWithFields(logrus.InfoLevel, "ℹ️", message)
}

func (d *diveLogger) Warn(message string) {
	d.logWithFields(logrus.WarnLevel, "⚠️", message)
}

func (d *diveLogger) Error(errorCode ErrorCode, errorMessage string) {
	d.logWithFields(logrus.ErrorLevel, "🛑", "Code:%d Error: %s", errorCode, errorMessage)
}

func (d *diveLogger) Fatal(errorCode ErrorCode, errorMessage string) {
	d.logWithFields(logrus.FatalLevel, "💀", "%s", errorMessage)
	d.log.Exit(1)
}

func (d *diveLogger) Infof(format string, args ...interface{}) {
	d.logWithFields(logrus.InfoLevel, "ℹ️", format, args...)
}

func (d *diveLogger) Warnf(format string, args ...interface{}) {
	d.logWithFields(logrus.WarnLevel, "⚠️", format, args...)
}

func (d *diveLogger) Debugf(format string, args ...interface{}) {
	d.logWithFields(logrus.DebugLevel, "🐞", format, args...)
}

func (d *diveLogger) Errorf(errorCode ErrorCode, format string, args ...interface{}) {
	d.logWithFields(logrus.ErrorLevel, "🛑", format, args...)

}

func (d *diveLogger) Fatalf(errorCode ErrorCode, format string, args ...interface{}) {
	d.logWithFields(logrus.FatalLevel, "💀", format, args...)
	d.log.Exit(1)
}
