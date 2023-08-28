package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/google/go-github/github"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/natefinch/lumberjack"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type DiveserviceResponse struct {
	ServiceName     string `json:"service_name,omitempty"`
	PublicEndpoint  string `json:"endpoint_public,omitempty"`
	PrivateEndpoint string `json:"endpoint,omitempty"`
	KeyPassword     string `json:"keypassword,omitempty"`
	KeystorePath    string `json:"keystore_path,omitempty"`
	Network         string `json:"network,omitempty"`
	NetworkName     string `json:"network_name,omitempty"`
	NetworkId       string `json:"nid,omitempty"`
	ChainId         string `json:"chain_id,omitempty"`
	ChainKey        string `json:"chain_key,omitempty"`
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

func ValidateCmdArgs(diveContext *DiveContext, args []string, cmd string) {
	if len(args) != 0 {

		diveContext.FatalError("Invalid Usage of command", cmd)

	}
}

func setupLogger() *logrus.Logger {
	pwd, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	log := logrus.New()

	log.SetOutput(io.Discard)
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceColors:     true,
		PadLevelText:    true,
	})

	ditFilePath := pwd + DiveLogDirectory + DiveDitLogFile
	errorFilePath := pwd + DiveLogDirectory + DiveErorLogFile

	ditLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:  filepath.ToSlash(ditFilePath),
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
			logrus.ErrorLevel: errorWriter,
			logrus.FatalLevel: errorWriter,
			logrus.WarnLevel:  errorWriter,
		},
		&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
	))

	return log
}

type Services map[string]*DiveserviceResponse

func WriteToServiceFile(serviceName string, data DiveserviceResponse) error {

	pwd, err := getPwd()
	if err != nil {
		return err
	}

	serviceFile := fmt.Sprintf("%s/%s", pwd, ServiceFilePath)

	file, err := os.OpenFile(serviceFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer file.Close()

	jsonDataFromFile, err := ReadServiceJsonFile()

	if err != nil {
		return err
	}

	if len(jsonDataFromFile) == 0 {
		jsonDataFromFile = Services{}
	}

	var dataToWrite []byte

	_, ok := jsonDataFromFile[serviceName]
	if !ok {
		jsonDataFromFile[serviceName] = &data
		dataToWrite, err = json.Marshal(jsonDataFromFile)
		if err != nil {
			return err
		}
	}

	_, err = file.Write(dataToWrite)

	if err != nil {
		return err
	}

	return nil
}

func ReadServiceJsonFile() (Services, error) {

	services := Services{}

	pwd, err := getPwd()
	if err != nil {
		return nil, err
	}
	serviceFile := fmt.Sprintf("%s/%s", pwd, ServiceFilePath)

	jsonFile, _ := os.ReadFile(serviceFile)

	if len(jsonFile) == 0 {
		return nil, nil
	}
	json.Unmarshal(jsonFile, &services)

	return services, nil

}

func getPwd() (string, error) {
	pwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return pwd, nil
}
