package common

import (
	"context"
	"io/fs"
	"os"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/spf13/cobra"
)

// Logger represents a generic logging interface with various log levels and formatting options.
type Logger interface {
	// SetErrorToStderr configures the logger to output error messages to stderr.
	SetErrorToStderr()

	// SetOutputToStdout configures the logger to output messages to stdout.
	SetOutputToStdout()

	// Debug logs a debug message.
	Debug(message string)

	// Info logs an informational message.
	Info(message string)

	// Warn logs a warning message.
	Warn(message string)

	// Error logs an error with an error code and a corresponding error message.
	Error(errorCode ErrorCode, errorMessage string)

	// Fatal logs a fatal error with an error code and a corresponding error message,
	// and then exits the program.
	Fatal(errorCode ErrorCode, errorMessage string)

	// Infof logs a formatted informational message.
	Infof(format string, args ...interface{})

	// Warnf logs a formatted warning message.
	Warnf(format string, args ...interface{})

	// Debugf logs a formatted debug message.
	Debugf(format string, args ...interface{})

	// Errorf logs a formatted error with an error code and a corresponding error message.
	Errorf(errorCode ErrorCode, format string, args ...interface{})

	// Fatalf logs a formatted fatal error with an error code and a corresponding error message,
	// and then exits the program.
	Fatalf(errorCode ErrorCode, format string, args ...interface{})
}

// Spinner is an interface for managing and controlling a terminal spinner.
type Spinner interface {
	// SetSuffixMessage sets a suffix message to be displayed beside the spinner.
	SetSuffixMessage(message, color string)

	// SetPrefixMessage sets a prefix message to be displayed beside the spinner.
	SetPrefixMessage(message string)

	// SetColor sets the color of the spinner.
	SetColor(color string)

	// Start starts the spinner with the specified color.
	Start(color string)

	// StartWithMessage starts the spinner with a specified message and color.
	StartWithMessage(message, color string)

	// Stop stops the spinner.
	Stop()

	// StopWithMessage stops the spinner and displays a final message.
	StopWithMessage(message string)
}

// Context represents a context for managing and interacting with enclaves and services.
type Context interface {
	// GetContext returns the underlying context.Context.
	GetContext() context.Context

	// GetKurtosisContext returns the Kurtosis context, including configuration and utility functions.
	GetKurtosisContext() (*kurtosis_context.KurtosisContext, error)

	// GetEnclaves retrieves information about all enclaves currently running.
	GetEnclaves() ([]EnclaveInfo, error)

	// GetEnclaveContext retrieves the context of a specific enclave by its name.
	GetEnclaveContext(enclaveName string) (*enclaves.EnclaveContext, error)

	// CleanEnclaves stops and cleans up all running enclaves.
	CleanEnclaves() ([]*EnclaveInfo, error)

	// CleanEnclaveByName stops and cleans up a specific enclave by its name.
	CleanEnclaveByName(enclaveName string) error

	// CheckSkippedInstructions checks if specific instructions are skipped in the current kurtosis run context.
	CheckSkippedInstructions(instructions map[string]bool) bool

	// StopService stops a specific service within an enclave by name.
	StopService(serviceName string, enclaveName string) error

	// StopServices stops all services within a specific enclave.
	StopServices(enclaveName string) error

	// RemoveServices stops and removes all services within a specific enclave.
	RemoveServices(enclaveName string) error

	// RemoveService stops and removes a specific service within an enclave by name.
	RemoveService(serviceName string, enclaveName string) error

	// RemoveServicesByServiceNames stops and removes services within an enclave based on a map of service names.
	RemoveServicesByServiceNames(services map[string]string, enclaveName string) error

	// CreateEnclave creates a new enclave with the specified name and returns its context.
	CreateEnclave(enclaveName string) (*enclaves.EnclaveContext, error)

	// Get the short UUID of the given enclave.
	GetShortUuid(enclaveName string) (string, error)

	// Exit terminates the execution of the context with the given status code.
	Exit(statusCode int)
}

// FileHandler defines methods for handling file-related operations.
type FileHandler interface {
	// ReadFile reads the contents of a file specified by the filePath.
	ReadFile(filePath string) ([]byte, error)

	// ReadJson reads the contents of a JSON file specified by the fileName
	// and unmarshals it into the provided object (obj).
	ReadJson(fileName string, obj interface{}) error

	// ReadAppFile reads the contents of an application-specific file specified by the fileName.
	ReadAppFile(fileName string) ([]byte, error)

	// WriteFile writes the provided data to a file specified by the fileName.
	WriteFile(fileName string, data []byte) error

	// WriteJson writes the provided data, marshaled as JSON, to a file specified by the fileName.
	WriteJson(fileName string, data interface{}) error

	// WriteAppFile writes the provided data to an application-specific file specified by the fileName.
	WriteAppFile(fileName string, data []byte) error

	// GetPwd returns the current working directory.
	GetPwd() (string, error)

	// GetHomeDir returns the home directory of the user.
	GetHomeDir() (string, error)

	// MkdirAll creates a directory along with any necessary parents,
	// and sets the specified permission.
	MkdirAll(dirPath string, permission fs.FileMode) error

	// OpenFile opens a file with the specified fileOpenMode and permission.
	OpenFile(filePath string, fileOpenMode string, permission int) (*os.File, error)

	// RemoveFile removes the file specified by the fileName.
	RemoveFile(fileName string) error

	// RemoveFiles removes multiple files specified by the fileNames.
	RemoveFiles(fileNames []string) error

	// RemoveDir removes the output directories by enclaveName.
	RemoveDir(enclaveName string) error

	// RemoveDir removes all the output directories.
	RemoveAllDir() error

	// GetAppDirPathOrAppFilePath returns the path to the application directory or a specific file within it
	// based on the provided fileName.
	GetAppDirPathOrAppFilePath(fileName string) (string, error)
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

	// Add StringSliceFlag adds a slice of string flag to the command that persists with short hand
	AddStringSliceFlagWithShortHand(stringV *[]string, name string, shorthand string, value []string, usage string) CommandBuilder

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

	// ToggleHelpCommand enables or disables the automatic generation of a help command.
	ToggleHelpCommand(enable bool) CommandBuilder

	// SetRunE sets the RunE field of the command.
	SetRunE(run func(cmd *cobra.Command, args []string) error) CommandBuilder

	// MarkFlagsAsRequired marks multiple flags as required for the command.
	MarkFlagsAsRequired(flags []string) CommandBuilder

	// MarkFlagRequired marks a flag as required for the command.
	MarkFlagRequired(flag string) CommandBuilder

	// AddBoolFlagP adds a boolean flag with shorthand to the command.
	AddBoolFlagP(name string, shorthand string, value bool, usage string) CommandBuilder
}
