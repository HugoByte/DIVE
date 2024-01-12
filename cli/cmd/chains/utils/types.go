package utils

import (
	"encoding/json"
	"fmt"
	"slices"

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

	err = cs.IsEmpty()
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}
	return nil
}

func (cc *CosmosServiceConfig) IsEmpty() error {
	if cc.ChainID == nil || cc.Key == nil || cc.Password == nil ||
		cc.PublicGrpc == nil || cc.PublicHTTP == nil || cc.PublicTCP == nil || cc.PublicRPC == nil {
		return common.WrapMessageToErrorf(common.ErrEmptyFields, "Missing Fields In The Config File")
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

	availablePort, err := common.GetAvailablePort()
	if err != nil {
		return err
	}
	sc.PublicPort = availablePort

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
	err = sc.IsEmpty()
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}

	return nil
}

func (sc *IconServiceConfig) IsEmpty() error {
	if sc.Port == 0 || sc.PublicPort == 0 || sc.P2PListenAddress == "" || sc.P2PAddress == "" || sc.Cid == "" {
		return common.WrapMessageToErrorf(common.ErrEmptyFields, "Missing Fields In The Config File")
	}
	return nil
}

type HardhatServiceConfig struct {
	PublicPort int `json:"public_port"`
}

func (sc *HardhatServiceConfig) LoadDefaultConfig() error {
	availablePort, err := common.GetAvailablePort()
	if err != nil {
		return err
	}
	sc.PublicPort = availablePort
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
	NodeType   string `json:"node_type"`
	Prometheus bool   `json:"prometheus"`
	Ports      Ports  `json:"ports"`
}

type Ports struct {
	RPCPort        int `json:"rpc_port"`
	Lib2LibPort    int `json:"lib2lib_port"`
	PrometheusPort int `json:"prometheus_port,omitempty"`
}

type RelayChainConfig struct {
	Name  string       `json:"name,omitempty"`
	Nodes []NodeConfig `json:"nodes,omitempty"`
}

type ParaNodeConfig struct {
	Name  string       `json:"name"`
	Nodes []NodeConfig `json:"nodes"`
}

type PolkadotServiceConfig struct {
	ChainType  string           `json:"chain_type"`
	RelayChain RelayChainConfig `json:"relaychain"`
	Para       []ParaNodeConfig `json:"parachains"`
	Explorer   bool             `json:"explorer"`
}

func (pc *ParaNodeConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(pc)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
}

func (rc *RelayChainConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(rc)
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	return string(encodedBytes), nil
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

	err = sc.IsEmpty()
	if err != nil {
		return common.WrapMessageToError(err, "Failed To Load Configuration")
	}

	for i := range sc.RelayChain.Nodes {
		sc.RelayChain.Nodes[i].AssignPorts(sc.RelayChain.Nodes[i].Prometheus)
	}

	for _, parachain := range sc.Para {
		for i := range parachain.Nodes {
			parachain.Nodes[i].AssignPorts(parachain.Nodes[i].Prometheus)
		}
	}

	return nil
}

func (sc *PolkadotServiceConfig) LoadDefaultConfig() error {
	sc.ChainType = "local"
	sc.Explorer = false
	sc.RelayChain.Name = "rococo-local"
	sc.RelayChain.Nodes = []NodeConfig{
		{Name: "bob", NodeType: "validator", Prometheus: false},
		{Name: "alice", NodeType: "validator", Prometheus: false},
	}

	for i := range sc.RelayChain.Nodes {
		sc.RelayChain.Nodes[i].AssignPorts(sc.RelayChain.Nodes[i].Prometheus)
	}

	sc.Para = []ParaNodeConfig{}
	return nil
}

func (nc *NodeConfig) AssignPorts(prometheus bool) error {
	var rpcPort, lib2libPort, prometheusPort int
	var err error
	rpcPort, err = common.GetAvailablePort()
	if err != nil {
		return err
	}

	lib2libPort, err = common.GetAvailablePort()
	if err != nil {
		return err
	}
	if prometheus {
		prometheusPort, err = common.GetAvailablePort()
		if err != nil {
			return err
		}
	}
	nc.Ports = Ports{RPCPort: rpcPort, Lib2LibPort: lib2libPort, PrometheusPort: prometheusPort}
	return nil
}

func (psc *PolkadotServiceConfig) IsEmpty() error {

	if psc == nil || psc.ChainType == "" {
		return common.WrapMessageToError(common.ErrEmptyFields, "Missing Fields In PolkadotServiceConfig")
	}

	if err := psc.RelayChain.IsEmpty(); err != nil {
		return err
	}

	if psc.RelayChain.Name == "" && len(psc.RelayChain.Nodes) == 0 && len(psc.Para) == 0 {
		return common.WrapMessageToError(common.ErrEmptyFields, "Missing Fields In RelayChainConfig")
	}

	if len(psc.Para) == 0 {
		return nil
	} else {
		for _, para := range psc.Para {
			if err := para.IsEmpty(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rcc *RelayChainConfig) IsEmpty() error {

	if rcc.Name == "" && len(rcc.Nodes) == 0 {
		return nil
	}

	if rcc.Name == "" || len(rcc.Nodes) == 0 {
		return common.WrapMessageToError(common.ErrEmptyFields, "Missing Fields In RelayChainConfig")
	}

	for _, node := range rcc.Nodes {
		if err := node.IsEmpty(); err != nil {
			return err
		}
	}

	return nil
}

func (pnc *ParaNodeConfig) IsEmpty() error {

	if pnc == nil || pnc.Name == "" || len(pnc.Nodes) == 0 {
		return common.WrapMessageToError(common.ErrEmptyFields, "Missing Fields In ParaNodeConfig")
	}

	for _, node := range pnc.Nodes {
		if err := node.IsEmpty(); err != nil {
			return err
		}
	}

	return nil
}

func (nc *NodeConfig) IsEmpty() error {

	if nc == nil || nc.Name == "" || nc.NodeType == "" {
		return common.WrapMessageToError(common.ErrEmptyFields, "Missing Fields In NodeConfig")
	}

	return nil
}

func (sc *PolkadotServiceConfig) HasPrometheus() bool {
	// Check relay chain nodes
	if sc.RelayChain.Name != "" {
		for _, node := range sc.RelayChain.Nodes {
			if node.Prometheus {
				return true
			}
		}
	}

	// Check para nodes
	for _, paraNode := range sc.Para {
		for _, node := range paraNode.Nodes {
			if node.Prometheus {
				return true
			}
		}
	}

	return false
}

func (sc *PolkadotServiceConfig) ValidateConfig() error {
	var validChainTypes = []string{"local", "testnet", "mainnet"}
	var validRelayNodeType = []string{"validator", "full"}
	var validParaNodeType = []string{"collator", "full"}
	var invalidTestNetParaChains = []string{"parallel", "subzero"}

	if !slices.Contains(validChainTypes, sc.ChainType) {
		return fmt.Errorf("invalid Chain Type: %s", sc.ChainType)
	}

	if sc.ChainType == "local" {
		if sc.RelayChain.Name != "rococo-local" {
			return fmt.Errorf("invalid Chain Name for local: %s", sc.RelayChain.Name)
		}
		if len(sc.RelayChain.Nodes) < 2 {
			return fmt.Errorf("atleast two nodes required for Relay Chain Local")
		}
		for _, node := range sc.RelayChain.Nodes {
			if node.NodeType != "validator" {
				return fmt.Errorf("invalid Node Type for Relay Chain Local: %s", node.NodeType)
			}
		}
	} else {
		for _, node := range sc.RelayChain.Nodes {
			if !slices.Contains(validRelayNodeType, node.NodeType) {
				return fmt.Errorf("invalid Node Type for Relay Chain: %s", node.NodeType)
			}
		}
	}

	if sc.RelayChain.Name != "" {
		if sc.ChainType == "testnet" && !(sc.RelayChain.Name == "rococo" || sc.RelayChain.Name == "westend") {
			return fmt.Errorf("invalid Chain Name for testnet: %s", sc.RelayChain.Name)
		}
		if sc.ChainType == "mainnet" && !(sc.RelayChain.Name == "kusama" || sc.RelayChain.Name == "polkadot") {
			return fmt.Errorf("invalid Chain Name for mainnet: %s", sc.RelayChain.Name)
		}
	}

	if sc.ChainType == "testnet" {
		for _, paraChain := range sc.Para {
			if slices.Contains(invalidTestNetParaChains, paraChain.Name) {
				return fmt.Errorf("no Testnet for Para Chain: %s", paraChain.Name)
			}
		}
	}

	for _, paraChain := range sc.Para {
		for _, node := range paraChain.Nodes {
			if !slices.Contains(validParaNodeType, node.NodeType) {
				return fmt.Errorf("invalid Node Type for Para Chain: %s", node.NodeType)
			}
		}
	}

	return nil
}

func (sc *PolkadotServiceConfig) GetParamsForRelay() (string, error) {
	relay_nodes, err := sc.RelayChain.EncodeToString()
	if err != nil {
		return "", common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	if sc.ChainType == "local" {
		return fmt.Sprintf(`{"relaychain": %s}`, relay_nodes), nil
	} else {
		return fmt.Sprintf(`{"chain_type": "%s", "relaychain": %s}`, sc.ChainType, relay_nodes), nil
	}
}

func (sc *PolkadotServiceConfig) ConfigureMetrics() {
	for i := range sc.RelayChain.Nodes {
		sc.RelayChain.Nodes[i].Prometheus = true
	}
	if len(sc.Para) != 0 {
		for i := range sc.Para[0].Nodes {
			sc.Para[0].Nodes[i].Prometheus = true
		}
	}
}

func (sc *PolkadotServiceConfig) ConfigureFullNodes(network string) {

	if network == "testnet" {
		sc.RelayChain.Name = "rococo"
	} else if network == "mainnet" {
		sc.RelayChain.Name = "kusama"
	}

	sc.RelayChain.Nodes = []NodeConfig{}

	sc.RelayChain.Nodes = append(sc.RelayChain.Nodes, NodeConfig{
		Name:       "alice",
		NodeType:   "full",
		Prometheus: false,
	})
}
