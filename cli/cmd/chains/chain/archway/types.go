package archway

import (
	"encoding/json"

	"github.com/hugobyte/dive-core/cli/common"
)

type ArchwayServiceConfig struct {
	Cid            *string `json:"chain_id"`
	Key            *string `json:"key"`
	PublicGrpcPort *int    `json:"public_grpc"`
	PublicHttpPort *int    `json:"public_http"`
	PublicTcpPort  *int    `json:"public_tcp"`
	PublicRpcPort  *int    `json:"public_rpc"`
	Password       *string `json:"password"`
}

func (as *ArchwayServiceConfig) LoadDefaultConfig() {
	as.Cid = nil
	as.Key = nil
	as.Password = nil
	as.PublicGrpcPort = nil
	as.PublicHttpPort = nil
	as.PublicRpcPort = nil
	as.PublicTcpPort = nil
	as.Password = nil
}

func (as *ArchwayServiceConfig) EncodeToString() (string, error) {

	data, err := json.Marshal(as)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (as *ArchwayServiceConfig) LoadConfigFromFile(cliContext *common.Cli, filePath string) error {

	err := cliContext.FileHandler().ReadJson(filePath, as)
	if err != nil {
		return err
	}
	return nil
}
