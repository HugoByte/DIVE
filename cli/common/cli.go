package common

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	cliContext *Cli

	initOnce sync.Once
)

type Cli struct {
	log         Logger
	spinner     Spinner
	context     Context
	fileHandler FileHandler
}

func initCli() (*Cli, error) {
	fileHandler := NewDiveFileHandler()

	pwd, err := fileHandler.GetPwd()

	if err != nil {
		return nil, WrapMessageToError(err, "Failed To Initialize CLi")
	}
	errorLogFilePath := filepath.Join(pwd, DiveLogDirectory, DiveErrorLogFile)
	infoLogFilePath := filepath.Join(pwd, DiveLogDirectory, DiveDitLogFile)

	return &Cli{
		log:         NewDiveLogger(infoLogFilePath, errorLogFilePath),
		spinner:     NewDiveSpinner(),
		context:     NewDiveContext1(),
		fileHandler: fileHandler,
	}, nil
}

func GetCli() *Cli {

	var err error
	initOnce.Do(func() {
		cliContext, err = initCli()
	})

	if err != nil {
		fmt.Println("Failed to get CLI context:", err)
		os.Exit(1)
	}

	return cliContext
}

func GetCliWithKurtosisContext() *Cli {

	_, err := cliContext.Context().GetKurtosisContext()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cliContext
}

func (c *Cli) Logger() Logger {
	return c.log
}

func (c *Cli) Spinner() Spinner {
	return c.spinner
}

func (c *Cli) Context() Context {
	return c.context
}

func (c *Cli) FileHandler() FileHandler {
	return c.fileHandler
}

func (c *Cli) Errorf(format string, err error, args ...interface{}) {

	c.spinner.Stop()
	c.log.SetErrorToStderr()
	actualError, _ := CoderOf(err)
	c.log.Errorf(actualError.ErrorCode(), "%s. message: %s", actualError.Error(), err.Error())
}

func (c *Cli) Fatalf(format string, err error, args ...interface{}) {
	c.spinner.Stop()
	c.log.SetErrorToStderr()
	actualError, _ := CoderOf(err)
	c.log.Fatalf(actualError.ErrorCode(), "%s. message: %s", actualError.Error(), err.Error())
}

func (c *Cli) Error(err error) {
	c.spinner.Stop()

	c.log.SetErrorToStderr()
	actualError, _ := CoderOf(err)
	c.log.Error(actualError.ErrorCode(), fmt.Sprintf("%s. message: %s", actualError.Error(), err.Error()))
}

func (c *Cli) Fatal(err error) {
	c.spinner.Stop()
	c.log.SetErrorToStderr()
	actualError, _ := CoderOf(err)
	c.log.Fatal(actualError.ErrorCode(), fmt.Sprintf("%s. message: %s", actualError.Error(), err.Error()))
}

func (c *Cli) Info(message string) {
	c.log.Info(message)
}
func (c *Cli) Infof(format string, args ...interface{}) {
	c.log.Infof(format, args...)
}

func (c *Cli) Warn(message string) {
	c.log.Warn(message)
}
func (c *Cli) Warnf(format string, args ...interface{}) {
	c.log.Warnf(format, args...)
}

func (c *Cli) Debug(message string) {
	c.log.Debug(message)
}
func (c *Cli) Debugf(format string, args ...interface{}) {
	c.log.Debugf(format, args...)
}

func (c *Cli) StartSpinnerIfNotVerbose(message string, verbose bool) {
	if verbose {
		c.log.SetOutputToStdout()
		c.log.Info(message)

	} else {
		c.spinner.StartWithMessage(message, "green")
	}
}

func (c *Cli) StopSpinnerIfNotVerbose(message string, verbose bool) {
	if verbose {
		c.log.SetOutputToStdout()
		c.log.Info(message)
	} else {
		c.spinner.StopWithMessage(message)
	}
}
