package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/go-github/github"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

type DiveserviceResponse struct {
	ServiceName     string `json:"service_name"`
	PublicEndpoint  string `json:"endpoint_public"`
	PrivateEndpoint string `json:"endpoint"`
	KeyPassword     string `json:"keypassword"`
	KeystorePath    string `json:"keystore_path"`
	Network         string `json:"network"`
	NetworkName     string `json:"network_name"`
	NetworkId       string `json:"nid"`
}

func (dive *DiveserviceResponse) Decode(responseData []byte) (*DiveserviceResponse, error) {

	err := json.Unmarshal(responseData, &dive)
	if err != nil {
		return nil, err
	}
	return dive, nil
}
func (dive *DiveserviceResponse) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(dive)
	if err != nil {
		return "", nil
	}

	return string(encodedBytes), nil
}
func (dive *DiveserviceResponse) WriteDiveResponse(diveContext *DiveContext) error {

	serialisedData, err := dive.EncodeToString()

	if err != nil {
		return err
	}

	return WriteToFile(serialisedData)
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

// This function will fetch the latest version from HugoByte/Dive repo
func GetLatestVersion() string {

	// Repo Name
	repo := "DIVE"
	owner := "HugoByte"

	// Create a new github client
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Print the release version.
	return release.GetName()
}

type DiveContext struct {
	Ctx             context.Context
	KurtosisContext *kurtosis_context.KurtosisContext
	Log             *logrus.Logger
	spinner         *spinner.Spinner
}

func NewDiveContext() *DiveContext {

	ctx := context.Background()
	log := logrus.New()

	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	if _, err := os.Stat(pwd + DiveLogDirectory); os.IsNotExist(err) {
		err := os.Mkdir(pwd+DiveLogDirectory, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	diwLogfile, err := os.OpenFile(pwd+DiveLogDirectory+DiveDiwLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	errorLogfile, err := os.OpenFile(pwd+DiveLogDirectory+DiveErorLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(io.Discard) // Send all logs to nowhere by default

	log.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: errorLogfile,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: diwLogfile,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})
	logFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log.SetFormatter(logFormatter)

	spinner := spinner.New(spinner.CharSets[80], 100*time.Millisecond, spinner.WithWriter(os.Stderr))

	return &DiveContext{Ctx: ctx, Log: log, spinner: spinner}
}

func (diveContext *DiveContext) GetEnclaveContext() (*enclaves.EnclaveContext, error) {

	_, err := diveContext.KurtosisContext.GetEnclave(diveContext.Ctx, DiveEnclave)
	if err != nil {
		enclaveCtx, err := diveContext.KurtosisContext.CreateEnclave(diveContext.Ctx, DiveEnclave, false)
		if err != nil {
			return nil, err

		}
		return enclaveCtx, nil
	}
	enclaveCtx, err := diveContext.KurtosisContext.GetEnclaveContext(diveContext.Ctx, DiveEnclave)

	if err != nil {
		return nil, err
	}
	return enclaveCtx, nil
}

func ReadConfigFile(filePath string) ([]byte, error) {

	file, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}
func WriteToFile(data string) error {

	pwd, err := os.Getwd()

	if err != nil {
		return err
	}

	file, err := os.OpenFile(pwd+"/dive.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)

	if err != nil {
		return err
	}

	return nil
}

// To get names of running enclaves, returns empty string if no enclaves
func (diveContext *DiveContext) GetEnclaves() string {
	enclaves, err := diveContext.KurtosisContext.GetEnclaves(diveContext.Ctx)
	if err != nil {
		logrus.Errorf("Getting Enclaves failed with error:  %v", err)
	}
	enclaveMap := enclaves.GetEnclavesByName()
	for _, enclaveInfoList := range enclaveMap {
		for _, enclaveInfo := range enclaveInfoList {
			return enclaveInfo.GetName()
		}
	}
	return ""
}

// Funstionality to clean the enclaves
func (diveContext *DiveContext) Clean() {
	diveContext.Log.SetOutput(os.Stdout)
	diveContext.Log.Info("Successfully connected to kurtosis engine...")
	diveContext.Log.Info("Initializing cleaning process...")
	// shouldCleanAll set to true as default for beta release.
	enclaves, err := diveContext.KurtosisContext.Clean(diveContext.Ctx, true)
	if err != nil {
		diveContext.Log.SetOutput(os.Stderr)
		diveContext.Log.Errorf("Failed cleaning with error: %v", err)
	}

	// Assuming only one enclave is running for beta release
	diveContext.Log.Infof("Successfully destroyed and cleaned enclave %s", enclaves[0].Name)
}

func (diveContext *DiveContext) FatalError(message, err string) {

	diveContext.Log.Fatalf("%s : %s", message, err)
}

func (diveContext *DiveContext) Info(message string) {

	diveContext.Log.Infoln(message)
}

func (diveContext *DiveContext) StartSpinner(message string) {

	diveContext.spinner.Suffix = message
	diveContext.spinner.Color("green")

	diveContext.spinner.Start()
}

func (diveContext *DiveContext) SetSpinnerMessage(message string) {

	diveContext.spinner.Suffix = message
}

func (diveContext *DiveContext) StopSpinner(message string) {
	c := color.New(color.FgCyan).Add(color.Underline)
	diveContext.spinner.FinalMSG = c.Sprintln(message)
	diveContext.spinner.Stop()

}

func (diveContext *DiveContext) GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]bool, error) {
	if DiveLogs {
		diveContext.spinner.Stop()
		diveContext.Log.SetFormatter(&logrus.TextFormatter{})
		diveContext.Log.SetOutput(os.Stdout)
	}
	var serializedOutputObj string
	skippedInstruction := map[string]bool{}
	for executionResponse := range response {

		if executionResponse.GetError() != nil {

			return "", nil, errors.New(executionResponse.GetError().String())
		}

		diveContext.Log.Info(executionResponse.String())

		if executionResponse.GetInstruction().GetIsSkipped() {
			skippedInstruction[executionResponse.GetInstruction().GetExecutableInstruction()] = executionResponse.GetInstruction().GetIsSkipped()
		}

		runFinishedEvent := executionResponse.GetRunFinishedEvent()

		if runFinishedEvent != nil {

			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()

			} else {
				return "", nil, errors.New(executionResponse.GetError().String())
			}

		} else {
			diveContext.spinner.Color("blue")
			if executionResponse.GetProgressInfo() != nil {
				c := color.New(color.FgGreen)
				diveContext.spinner.Suffix = c.Sprintf(strings.ReplaceAll(executionResponse.GetProgressInfo().String(), "current_step_info:", " "))

			}
		}

	}

	return serializedOutputObj, skippedInstruction, nil

}

func ValidateCmdArgs(args []string, cmd string) {
	if len(args) != 0 {
		logrus.Fatalf("Invalid Usage of command. Find cmd %s", cmd)

	}
}

func (diveContext *DiveContext) Error(err string) {
	diveContext.Log.Error(err)
}

func (diveContext *DiveContext) InitKurtosisContext() {
	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		log.Fatal("The Kurtosis Engine Server is unavailable and is probably not running; you will need to start it using the Kurtosis CLI before you can create a connection to it")

	}
	diveContext.KurtosisContext = kurtosisContext
}

func (diveContext *DiveContext) CheckInstructionSkipped(instuctions map[string]bool, message string) {

	if len(instuctions) != 0 {

		diveContext.StopSpinner(message)
		os.Exit(0)
	}
}
