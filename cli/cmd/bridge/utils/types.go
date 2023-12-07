package utils

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hugobyte/dive-core/cli/common"
)

const (
	IconChain    = "icon"
	EthChain     = "eth"
	HardhatChain = "hardhat"
	ArchwayChain = "archway"
	NeutronChain = "neutron"
)

var supportedChainsForBtp = []string{IconChain, EthChain, HardhatChain}
var supportedChainsForIbc = []string{ArchwayChain, NeutronChain, IconChain}

type Chains struct {
	ChainA            string
	ChainB            string
	ChainAServiceName string
	ChainBServiceName string
	Bridge            string
}

func InitChains(chainA, chainB, serviceA, serviceB string, bridge bool) *Chains {
	return &Chains{

		ChainA:            strings.ToLower(chainA),
		ChainB:            strings.ToLower(chainB),
		ChainAServiceName: serviceA,
		ChainBServiceName: serviceB,
		Bridge:            strconv.FormatBool(bridge),
	}
}

func (c *Chains) AreChainsIcon() bool {
	return (c.ChainA == "icon" && c.ChainB == "icon")
}

func (chains *Chains) GetParams() string {
	return fmt.Sprintf(`{"src_chain": "%s", "dst_chain": "%s", "bridge":"%s"}`, chains.ChainA, chains.ChainB, chains.Bridge)
}
func (chains *Chains) GetIbcRelayParams() string {

	return fmt.Sprintf(`{"src_chain": "%s", "dst_chain": "%s"}`, chains.ChainA, chains.ChainB)
}

func (chains *Chains) GetServicesResponse(cli *common.Cli) (string, string, error) {

	var serviceConfig = common.Services{}

	err := cli.FileHandler().ReadJson("services.json", &serviceConfig)

	if err != nil {
		return "", "", err
	}

	chainAServiceResponse, OK := serviceConfig[chains.ChainAServiceName]
	if !OK {
		return "", "", fmt.Errorf("service name not found")
	}
	chainBServiceResponse, OK := serviceConfig[chains.ChainBServiceName]
	if !OK {
		return "", "", fmt.Errorf("service name not found")
	}

	srcChainServiceResponse, err := chainAServiceResponse.EncodeToString()
	if err != nil {
		return "", "", err
	}
	dstChainServiceResponse, err := chainBServiceResponse.EncodeToString()

	if err != nil {
		return "", "", err
	}

	return srcChainServiceResponse, dstChainServiceResponse, nil
}

func (chains *Chains) CheckForBtpSupportedChains() error {
	if !slices.Contains(supportedChainsForBtp, chains.ChainA) {
		return fmt.Errorf("invalid Chain: %s", chains.ChainA)
	}
	if !slices.Contains(supportedChainsForBtp, chains.ChainB) {
		return fmt.Errorf("invalid Chain: %s", chains.ChainB)
	}
	return nil
}

func (chains *Chains) CheckForIbcSupportedChains() error {
	if !slices.Contains(supportedChainsForIbc, chains.ChainA) {
		return fmt.Errorf("invalid Chain: %s", chains.ChainA)
	}
	if !slices.Contains(supportedChainsForIbc, chains.ChainB) {
		return fmt.Errorf("invalid Chain: %s", chains.ChainB)
	}
	return nil
}

func (chains *Chains) CheckChainServiceNamesEmpty() bool {
	return (chains.ChainAServiceName != "" && chains.ChainBServiceName != "")
}
