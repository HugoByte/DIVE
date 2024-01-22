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

func (c *Chains) AreChainsCosmos() bool {
	return (c.ChainA == "archway" || c.ChainA == "neutron") && (c.ChainB == "archway" || c.ChainB == "neutron")
}

func (chains *Chains) GetParams(src_service_config string, dst_service_config string) string {
	return fmt.Sprintf(`{"src_chain": "%s", "dst_chain": "%s", "bridge":"%s", "src_service_config": %s, "dst_service_config": %s}`, chains.ChainA, chains.ChainB, chains.Bridge, src_service_config, dst_service_config)
}

func (chains *Chains) GetIbcRelayParams(src_service_config string, dst_service_config string) string {
	return fmt.Sprintf(`{"src_chain": "%s", "dst_chain": "%s", "src_service_config": %s, "dst_service_config": %s}`, chains.ChainA, chains.ChainB, src_service_config, dst_service_config)
}

func (chains *Chains) GetServicesResponse(cli *common.Cli) (string, string, error) {

	var serviceConfig = common.Services{}
	
	shortUuid, err := cli.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		return "", "", fmt.Errorf("failed to get short uuid of enclave")
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

	err = cli.FileHandler().ReadJson(serviceFileName, &serviceConfig)

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
		return fmt.Errorf("invalid chain: %s", chains.ChainA)
	}
	if !slices.Contains(supportedChainsForBtp, chains.ChainB) {
		return fmt.Errorf("invalid chain: %s", chains.ChainB)
	}
	return nil
}

func (chains *Chains) CheckForIbcSupportedChains() error {
	if !slices.Contains(supportedChainsForIbc, chains.ChainA) {
		return fmt.Errorf("invalid chain: %s", chains.ChainA)
	}
	if !slices.Contains(supportedChainsForIbc, chains.ChainB) {
		return fmt.Errorf("invalid chain: %s", chains.ChainB)
	}
	return nil
}

func (chains *Chains) CheckChainServiceNamesEmpty() bool {
	return (chains.ChainAServiceName != "" && chains.ChainBServiceName != "")
}
