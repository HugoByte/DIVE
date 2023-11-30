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
		// Log file abbsolute path, os agnostic
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
func (d *diveLogger) Debug(message string) {

	d.log.WithFields(logrus.Fields{
		"level": "🐞 debug",
	}).Debug(message)

}
func (d *diveLogger) Info(message string) {
	d.log.WithFields(logrus.Fields{
		"level": "ℹ️ info",
	}).Info(message)

}
func (d *diveLogger) Warn(message string) {
	d.log.WithFields(logrus.Fields{
		"level": "⚠️ warn",
	}).Warn(message)

}
func (d *diveLogger) Error(errorCode ErrorCode, errorMessage string) {
	d.log.WithFields(logrus.Fields{
		"level":      "🛑 error",
		"error_code": errorCode,
	}).Error(errorMessage)

}
func (d *diveLogger) Fatal(errorCode ErrorCode, errorMessage string) {
	d.log.WithFields(logrus.Fields{
		"level":      "💀 fatal",
		"error_code": errorCode,
	}).Fatal(errorMessage)

}
func (d *diveLogger) Infof(message string) {
	d.log.WithFields(logrus.Fields{
		"level": "ℹ️ info",
	}).Infof("%s", message)

}
func (d *diveLogger) Warnf(message string) {
	d.log.WithFields(logrus.Fields{
		"level": "⚠️ warn",
	}).Warnf("%s", message)

}
func (d *diveLogger) Debugf(message string) {

	d.log.WithFields(logrus.Fields{
		"level": "🐞 debug",
	}).Debugf("%s", message)

}
func (d *diveLogger) Errorf(errorCode ErrorCode, errorMessage string) {
	d.log.WithFields(logrus.Fields{
		"level":      "🛑 error",
		"error_code": errorCode,
	}).Errorf("%s", errorMessage)

}
func (d *diveLogger) Fatalf(errorCode ErrorCode, errorMessage string) {
	d.log.WithFields(logrus.Fields{
		"level":      "💀 fatal",
		"error_code": errorCode,
	}).Fatalf("%s", errorMessage)

}
