package common

import (
	"os"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
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
	SetMessage(message string, color string)
	SetColor(color string)
	Start(message string)
	Stop(message string)
}

type Context interface {
	CheckSkippedInstructions()
	CleanAll()
	Clean(enclaveName string)
	CreateEnclave(enclaveName string)
	GetEnclaves() []string
	GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error)
	InitialiseKurtosisContext()
	StopServices()
	StopService()
}

type FileHandler interface {
	ReadFile(filePath string) ([]byte, error)
	ReadJson(filePath string, obj interface{}) (string, error)
	WriteFile(filePath string, data []byte) error
	WriteJson(filePath string, data interface{}) error
	GetPwd() string
	MkdirAll(dirPath string, permission string) error
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
}
