package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/go-github/github"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
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
func (dive *DiveserviceResponse) WriteDiveResponse(diveContext *DiveContext) {

	serialisedData, err := dive.EncodeToString()

	if err != nil {
		diveContext.FatalError("Failed To Serialzed Data", err.Error())
	}

	WriteToFile(serialisedData)
}

func GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) string {

	var serializedOutputObj string

	for executionResponseLine := range response {

		runFinishedEvent := executionResponseLine.GetRunFinishedEvent()

		if runFinishedEvent == nil {

		} else {

			if runFinishedEvent.GetIsRunSuccessful() {

				serializedOutputObj = runFinishedEvent.GetSerializedOutput()
			} else {
				logrus.Fatal("Starlark run Fails")
			}
		}
	}

	return serializedOutputObj

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
	log             *logrus.Logger
}

func NewDiveContext() *DiveContext {

	ctx := context.Background()

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		logrus.Fatal("The Kurtosis Engine Server is unavailable and is probably not running; you will need to start it using the Kurtosis CLI before you can create a connection to it")

	}
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return &DiveContext{Ctx: ctx, KurtosisContext: kurtosisContext, log: logrus.New()}
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
func WriteToFile(data string) {
	file, err := os.Create("dive.json")
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(data)
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
	logrus.Info("Successfully connected to kurtosis engine...")
	logrus.Info("Initializing cleaning process...")

	// shouldCleanAll set to true as default for beta release.
	enclaves, err := diveContext.KurtosisContext.Clean(diveContext.Ctx, true)
	if err != nil {
		logrus.Errorf("Failed cleaning with error: %v", err)
	}

	// Assuming only one enclave is running for beta release
	logrus.Infof("Successfully destroyed and cleaned enclave %s", enclaves[0].Name)
}

func (diveContext *DiveContext) FatalError(message, err string) {

	diveContext.log.Fatalf("%s : %s", message, err)
}

func (diveContext *DiveContext) Info(message string) {

	diveContext.log.Info(message)
}
