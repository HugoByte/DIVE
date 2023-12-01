package common

import (
	"context"
	"io/fs"
	"os"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/spf13/cobra"
)

type Logger interface {
	SetErrorToStderr()
	SetOutputToStdout()
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(errorCode ErrorCode, errorMessage string)
	Fatal(errorCode ErrorCode, errorMessage string)
	Infof(message string)
	Warnf(message string)
	Debugf(message string)
	Errorf(errorCode ErrorCode, errorMessage string)
	Fatalf(errorCode ErrorCode, errorMessage string)
}

type Spinner interface {
	SetSuffixMessage(message, color string)
	SetPrefixMessage(message string)
	SetColor(color string)
	Start(color string)
	StartWithMessage(message, color string)
	Stop()
	StopWithMessage(message string)
}

type Context interface {
	GetContext() context.Context
	GetKurtosisContext() (*kurtosis_context.KurtosisContext, error)
	GetEnclaves() ([]EnclaveInfo, error)
	GetEnclaveContext(enclaveName string) (*enclaves.EnclaveContext, error)
	CleanEnclaves() ([]*EnclaveInfo, error)
	CleanEnclaveByName(enclaveName string) error
	CheckSkippedInstructions()
	StopService(serviceName string) error
	StopServices(enclaveName string) error
	RemoveServices(enclaveName string) error
	RemoveService(serviceName string, enclaveName string) error
	CreateEnclave(enclaveName string) (*enclaves.EnclaveContext, error)
	GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error)
}

type FileHandler interface {
	ReadFile(filePath string) ([]byte, error)
	ReadJson(fileName string, obj interface{}) error
	ReadAppFile(fileName string) ([]byte, error)
	WriteFile(fileName string, data []byte) error
	WriteJson(fileName string, data interface{}) error
	WriteAppFile(fileName string, data []byte) error
	GetPwd() (string, error)
	GetHomeDir() (string, error)
	MkdirAll(dirPath string, permission fs.FileMode) error
	OpenFile(filePath string, fileOpenMode string, permission int) (*os.File, error)
}

// CommandBuilder is an interface for building a Cobra command.
type CommandBuilder interface {
	// AddCommand adds a subcommand to the command.
	AddCommand(cmd *cobra.Command) CommandBuilder

	// Add Persistent Bool Flag
	AddBoolPersistentFlag(boolV *bool, name string, value bool, usage string) CommandBuilder

	// Add Persistent Bool Flag with Short hand
	AddBoolPersistentFlagWithShortHand(boolV *bool, name string, value bool, usage string, shorthand string) CommandBuilder

	// Add Persistent String Flag
	AddStringPersistentFlag(stringV *string, name string, value string, usage string) CommandBuilder

	// Add Persistent String Flag with Short hand
	AddStringPersistentFlagWithShortHand(stringV *string, name string, shorthand string, value string, usage string) CommandBuilder

	// Add StringFlag adds a string flag to the command that persists
	AddStringFlag(stringV *string, name string, value string, usage string) CommandBuilder

	// Add StringFlag adds a string flag to the command that persists with short hand
	AddStringFlagWithShortHand(stringV *string, name string, shorthand string, value string, usage string) CommandBuilder

	// Add BooFlag adds a boolean flag to the command that persists
	AddBoolFlag(boolV *bool, name string, value bool, usage string) CommandBuilder

	AddBoolFlagWithShortHand(boolV *bool, name string, shorthand string, value bool, usage string) CommandBuilder

	// Build constructs and returns the Cobra command.
	Build() *cobra.Command

	// SetUse sets the Use field of the command.
	SetUse(use string) CommandBuilder

	// SetShort sets the Short field of the command.
	SetShort(short string) CommandBuilder

	// SetLong sets the Long field of the command.
	SetLong(long string) CommandBuilder

	// SetRun sets the Run field of the command.
	SetRun(run func(cmd *cobra.Command, args []string)) CommandBuilder

	ToggleHelpCommand(enable bool) CommandBuilder

	SetRunE(run func(cmd *cobra.Command, args []string) error) CommandBuilder
}
