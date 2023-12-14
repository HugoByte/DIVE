package utils

import (
	"encoding/json"

	"github.com/hugobyte/dive-core/cli/common"
)

type CosmosServiceConfig struct {
	ChainID    *string `json:"chain_id"`
	ChainName  *string `json:"chain_name"`
	Key        *string `json:"key"`
	Password   *string `json:"password"`
	PublicGrpc *int    `json:"public_grpc"`
	PublicHTTP *int    `json:"public_http"`
	PublicTCP  *int    `json:"public_tcp"`
	PublicRPC  *int    `json:"public_rpc"`
}

func (cs *CosmosServiceConfig) LoadDefaultConfig() error {
	cs.ChainID = nil
	cs.Key = nil
	cs.Password = nil
	publicGrpc, err := common.GetAvailablePort()
	if err != nil {
		return common.WrapMessageToError(err, "error getting available gRPC port")
	}
	cs.PublicGrpc = &publicGrpc

	publicHTTP, err := common.GetAvailablePort()
	if err != nil {
		return common.WrapMessageToError(err, "error getting available HTTP port")
	}
	cs.PublicHTTP = &publicHTTP

	publicRPC, err := common.GetAvailablePort()
	if err != nil {
		return common.WrapMessageToError(err, "error getting available Rpc port")
	}
	cs.PublicRPC = &publicRPC

	publicTCP, err := common.GetAvailablePort()
	if err != nil {
		return common.WrapMessageToError(err, "error getting available Tcp port")
	}
	cs.PublicTCP = &publicTCP

	return nil
}

func (cs *CosmosServiceConfig) EncodeToString() (string, error) {

	data, err := json.Marshal(cs)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(data), nil
}

func (cs *CosmosServiceConfig) LoadConfigFromFile(cliContext *common.Cli, filePath string) error {

	err := cliContext.FileHandler().ReadJson(filePath, cs)
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}
	return nil
}

type IconServiceConfig struct {
	Port             int    `json:"private_port"`
	PublicPort       int    `json:"public_port"`
	P2PListenAddress string `json:"p2p_listen_address"`
	P2PAddress       string `json:"p2p_address"`
	Cid              string `json:"cid"`
}

func (sc *IconServiceConfig) LoadDefaultConfig() error {
	sc.Port = 9080
	sc.P2PListenAddress = "7080"
	sc.P2PAddress = "8080"
	sc.Cid = "0xacbc4e"

	if common.CheckPort(8090) {
		sc.PublicPort = 8090
	} else {
		availablePort, err := common.GetAvailablePort()
		if err != nil {
			return err
		}
		sc.PublicPort = availablePort
	}

	return nil

}

func (sc *IconServiceConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(sc)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
}

func (sc *IconServiceConfig) LoadConfigFromFile(cliContext *common.Cli, filePath string) error {
	err := cliContext.FileHandler().ReadJson(filePath, sc)
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}
	return nil
}

type HardhatServiceConfig struct {
	PublicPort int `json:"public_port"`
}

func (sc *HardhatServiceConfig) LoadDefaultConfig() error {
	sc.PublicPort = 8545
	if common.CheckPort(8545) {
		sc.PublicPort = 8545
	} else {
		availablePort, err := common.GetAvailablePort()
		if err != nil {
			return err
		}
		sc.PublicPort = availablePort
	}

	return nil
}

func (sc *HardhatServiceConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(sc)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
}

// This code is for polkadot config file

type NodeConfig struct {
	Name       string `json:"name"`
	NodeType   string `json:"node-type"`
	Port       int    `json:"port"`
	Prometheus bool   `json:"prometheus"`
}

type RelayChainConfig struct {
	Name  string       `json:"name"`
	Nodes []NodeConfig `json:"nodes"`
}

type ParaNodeConfig struct {
	Name  string       `json:"name"`
	Nodes []NodeConfig `json:"nodes"`
}

type PolkadotServiceConfig struct {
	ChainType  string           `json:"chain-type"`
	RelayChain RelayChainConfig `json:"relaychain"`
	Para       []ParaNodeConfig `json:"para"`
	Explorer   bool             `json:"explorer"`
}

func (sc *PolkadotServiceConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(sc)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
}

func (sc *PolkadotServiceConfig) LoadConfigFromFile(cliContext *common.Cli, filePath string) error {
	err := cliContext.FileHandler().ReadJson(filePath, sc)
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}
	return nil
}

func (sc *PolkadotServiceConfig) LoadDefaultConfig() error {
	sc.ChainType = "local"
	sc.Explorer = false
	sc.RelayChain.Name = "rococo-local"
	sc.RelayChain.Nodes = []NodeConfig{
		{Name: "bob", NodeType: "full", Port: 9944, Prometheus: false},
		{Name: "alice", NodeType: "validator", Port: 9955, Prometheus: false},
	}
	sc.Para = []ParaNodeConfig{
		{
			Name: "",
			Nodes: []NodeConfig{
				{Name: "alice", NodeType: "full", Prometheus: false},
				{Name: "bob", NodeType: "collator", Prometheus: false},
			},
		},
	}

	return nil
}
