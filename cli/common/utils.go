package common

import (
	"os/exec"
	"runtime"

	"github.com/kurtosis-tech/stacktrace"
)

func ValidateArgs(args []string) error {
	if len(args) != 0 {

		return Errorc(InvalidCommandError, "Invalid Usage Of Command Arguments")

	}
	return nil
}

func WriteServiceResponseData(serviceName string, data DiveServiceResponse, cliContext *Cli) error {
	var jsonDataFromFile = Services{}
	err := cliContext.FileHandler().ReadJson("services.json", &jsonDataFromFile)

	if err != nil {
		return err
	}

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = &data

	}
	err = cliContext.FileHandler().WriteJson("services.json", jsonDataFromFile)
	if err != nil {
		return err
	}

	return nil
}

func OpenFile(URL string) error {
	var args []string
	switch runtime.GOOS {
	case linuxOSName:
		args = []string{openFileLinuxCommandName, URL}
	case macOSName:
		args = []string{openFileMacCommandName, URL}
	case windowsOSName:
		args = []string{openFileWindowsCommandName, openFileWindowsCommandFirstArgumentDefault, URL}
	default:
		return stacktrace.NewError("Unsupported operating system")
	}
	command := exec.Command(args[0], args[1:]...)
	if err := command.Start(); err != nil {
		return stacktrace.Propagate(err, "An error occurred while opening '%v'", URL)
	}
	return nil
}
