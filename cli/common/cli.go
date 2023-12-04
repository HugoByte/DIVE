package common

import (
	"path/filepath"
	"sync"
)

var (
	CliContext *Cli

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

func GetCli() (*Cli, error) {
	var err error
	initOnce.Do(func() {
		CliContext, err = initCli()
	})

	if err != nil {
		return nil, WrapMessageToError(err, "Failed To Retrieve Cli")
	}

	return CliContext, nil
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
