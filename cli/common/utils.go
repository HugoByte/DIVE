package common

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/stacktrace"
)

// The function "ValidateArgs" checks if the given arguments are empty and returns an error if they are
// not.
func ValidateArgs(args []string) error {
	if len(args) != 0 {

		return ErrInvalidCommand

	}
	return nil
}

// The function writes service response data to a JSON file.
func WriteServiceResponseData(serviceName string, data DiveServiceResponse, cliContext *Cli, fileName string) error {
	var jsonDataFromFile = Services{}
	err := cliContext.FileHandler().ReadJson(fileName, &jsonDataFromFile)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Read %s", fileName)
	}

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = &data

	}
	err = cliContext.FileHandler().WriteJson(fileName, jsonDataFromFile)
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write %s", fileName)
	}

	return nil
}

// The function writes bridge response data to a JSON file.
func WriteBridgeResponseData(serviceName string, data string, cliContext *Cli, fileName string) error {
	bridgeResponse := DiveBridgeResponse{}
	diveBridgeResponse, err := bridgeResponse.Decode([]byte(data))
	if err != nil {
		return WrapMessageToErrorf(ErrDataUnMarshall, "Failed To Unmarshall data")
	}

	var jsonDataFromFile = BridgeServices{}
	err = cliContext.FileHandler().ReadJson(fileName, &jsonDataFromFile)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Read %s", fileName)
	}

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = diveBridgeResponse
	}

	err = cliContext.FileHandler().WriteJson(fileName, jsonDataFromFile)
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Write %s", fileName)
	}

	return nil
}

// The OpenFile function opens a file using the appropriate command based on the operating system.
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
		return WrapMessageToError(ErrUnsupportedOS, stacktrace.NewError("Unsupported operating system").Error())
	}
	command := exec.Command(args[0], args[1:]...)
	if err := command.Start(); err != nil {
		return WrapMessageToError(ErrOpenFile, stacktrace.Propagate(err, "An error occurred while opening '%v'", URL).Error())
	}
	return nil
}

// The function `LoadConfig` loads a configuration either from a default source or from a specified
// file path.
func LoadConfig(cliContext *Cli, config ConfigLoader, filePath string) error {
	if filePath == "" {
		if err := config.LoadDefaultConfig(); err != nil {
			return err
		}
	} else {
		err := config.LoadConfigFromFile(cliContext, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// This function returns a StarlarkRunConfig object with the provided parameters.
func GetStarlarkRunConfig(params string, relativePathToMainFile string, mainFunctionName string) *starlark_run_config.StarlarkRunConfig {

	starlarkConfig := &starlark_run_config.StarlarkRunConfig{
		RelativePathToMainFile:   relativePathToMainFile,
		MainFunctionName:         mainFunctionName,
		DryRun:                   DiveDryRun,
		SerializedParams:         params,
		Parallelism:              DiveDefaultParallelism,
		ExperimentalFeatureFlags: []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{},
	}
	return starlarkConfig
}

// The function `GetSerializedData` retrieves serialized data, service information, skipped
// instructions, and any errors from a given response channel.
func GetSerializedData(cliContext *Cli, response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error) {

	var serializedOutputObj string
	services := map[string]string{}

	skippedInstruction := map[string]bool{}
	for executionResponse := range response {

		if strings.Contains(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), "added with service") {
			res1 := strings.Split(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), " ")
			serviceName := res1[1][1 : len(res1[1])-1]
			serviceUUID := res1[len(res1)-1][1 : len(res1[len(res1)-1])-1]
			services[serviceName] = serviceUUID
		}

		cliContext.log.Info(executionResponse.String())

		if executionResponse.GetInstruction().GetIsSkipped() {
			skippedInstruction[executionResponse.GetInstruction().GetExecutableInstruction()] = executionResponse.GetInstruction().GetIsSkipped()
			break
		}

		if executionResponse.GetError() != nil {

			return "", services, nil, WrapMessageToError(ErrStarlarkResponse, executionResponse.GetError().String())

		}

		runFinishedEvent := executionResponse.GetRunFinishedEvent()

		if runFinishedEvent != nil {

			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()

			} else {
				return "", services, nil, WrapMessageToError(ErrStarlarkResponse, executionResponse.GetError().String())
			}

		} else {
			cliContext.spinner.SetColor("blue")
			if executionResponse.GetProgressInfo() != nil {

				cliContext.spinner.SetSuffixMessage(strings.ReplaceAll(executionResponse.GetProgressInfo().String(), "current_step_info:", " "), "fgGreen")

			}
		}

	}

	return serializedOutputObj, services, skippedInstruction, nil
}

// Check if a port is available
func CheckPort(port int) bool {

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	defer ln.Close()
	return true
}

// The function `GetAvailablePort` generates a random port number between 1024 and 65535 and checks if
// it is available by calling the `CheckPort` function, and returns the first available port or an
// error if no port is available.
func GetAvailablePort() (int, error) {

	// Check random ports in the range 1024-65535
	for i := 0; i < 1000; i++ {
		port := rand.Intn(64511) + 1024
		if CheckPort(port) {
			return port, nil
		}
	}

	return 0, ErrPortAllocation
}
