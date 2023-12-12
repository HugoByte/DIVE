package common

import (
	"encoding/json"
)

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

// The `Decode` function is a method of the `DiveServiceResponse` struct. It takes a byte slice
// `responseData` as input and attempts to decode it into a `DiveServiceResponse` object.
func (dive *DiveServiceResponse) Decode(responseData []byte) (*DiveServiceResponse, error) {

	err := json.Unmarshal(responseData, &dive)
	if err != nil {
		return nil, WrapMessageToError(ErrDataUnMarshall, err.Error())
	}
	return dive, nil
}

// The `EncodeToString` function is a method of the `DiveServiceResponse` struct. It encodes the
// `DiveServiceResponse` object into a JSON string representation.
func (dive *DiveServiceResponse) EncodeToString() (string, error) {

	encodedBytes, err := json.Marshal(dive)
	if err != nil {
		return "", WrapMessageToError(ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
}

type Services map[string]*DiveServiceResponse

// The EnclaveInfo type represents information about an enclave, including its name, UUID, short UUID,
// creation time, and status.
// @property {string} Name - The name of the enclave.
// @property {string} Uuid - The Uuid property is a unique identifier for the enclave. It is used to
// distinguish one enclave from another.
// @property {string} ShortUuid - The ShortUuid property is a shortened version of the Uuid property.
// It is typically used to provide a more concise representation of the unique identifier for the
// Enclave.
// @property {string} CreatedTime - The CreatedTime property in the EnclaveInfo struct represents the
// timestamp when the enclave was created. It is a string that typically follows a specific date and
// time format, such as "YYYY-MM-DD HH:MM:SS".
// @property {string} Status - The "Status" property in the EnclaveInfo struct represents the current
// status of the enclave. It can have different values depending on the implementation, but some common
// values could be "active", "inactive", "error", or "unavailable".
type EnclaveInfo struct {
	Name        string
	Uuid        string
	ShortUuid   string
	CreatedTime string
	Status      string
}

// The ConfigLoader interface defines methods for loading default configurations and configurations
// from a file.
// @property {error} LoadDefaultConfig - This method is responsible for loading the default
// configuration. It does not take any arguments and returns an error if there is any issue loading the
// default configuration.
// @property {error} LoadConfigFromFile - This method is responsible for loading a configuration file
// from a given file path. It takes two parameters: a cliContext object, which represents the current
// command-line interface context, and a filePath string, which represents the path to the
// configuration file. The method returns an error if there is any issue loading the
type ConfigLoader interface {
	LoadDefaultConfig() error
	LoadConfigFromFile(cliContext *Cli, filePath string) error
}
