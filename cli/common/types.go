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

type Services map[string]*DiveServiceResponse

type EnclaveInfo struct {
	Name      string
	Uuid      string
	ShortUuid string
}

type ConfigLoader interface {
	LoadDefaultConfig()
	LoadConfigFromFile(cliContext *Cli, filePath string) error
}
