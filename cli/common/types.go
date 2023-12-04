package common

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"time"

	"github.com/kurtosis-tech/stacktrace"
)

var lastChecked time.Time
var latestVersion = ""

type DiveServiceResponse struct {
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

func (dive *DiveServiceResponse) Decode(responseData []byte) (*DiveServiceResponse, error) {

	err := json.Unmarshal(responseData, &dive)
	if err != nil {
		return nil, err
	}
	return dive, nil
}
func (dive *DiveServiceResponse) EncodeToString() (string, error) {

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
// func GetLatestVersion() string {

// 	// Repo Name
// 	repo := "DIVE"
// 	owner := "HugoByte"
// 	userHomeDir, err := os.UserHomeDir()

// 	if err != nil {
// 		fmt.Println(err)
// 		return ""
// 	}
// 	cachedFile := filepath.Join(userHomeDir, "/.dive/version_cache.txt")

// 	if time.Since(lastChecked).Hours() > 1 {
// 		cachedVersion, err := ReadConfigFile(cachedFile)
// 		fmt.Println("here ")

// 		if err == nil && string(cachedVersion) != "" {
// 			latestVersion = string(cachedVersion)
// 			fmt.Println("here 1")
// 		} else {
// 			fmt.Println("here 2")
// 			client := github.NewClient(nil)
// 			release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
// 			if err != nil {
// 				fmt.Println(err)
// 				return ""
// 			}

// 			latestVersion = release.GetName()
// 			writeCache(cachedFile, latestVersion)
// 		}
// 		lastChecked = time.Now()

// 	}

// 	return latestVersion
// }

type Services map[string]*DiveServiceResponse

type EnclaveInfo struct {
	Name      string
	Uuid      string
	ShortUuid string
}
