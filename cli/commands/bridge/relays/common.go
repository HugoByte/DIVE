package relays

import (
	"fmt"
	"strconv"
	"strings"

	"slices"

	"github.com/hugobyte/dive/cli/common"
)

var supportedChainsForBtp = []string{"icon", "eth", "hardhat"}
var supportedChainsForIbc = []string{"archway", "neutron", "icon"}

type Chains struct {
	chainA            string
	chainB            string
	chainAServiceName string
	chainBServiceName string
	bridge            string
}

func initChains(chainA, chainB, serviceA, serviceB string, bridge bool) *Chains {
	return &Chains{
		chainA:            strings.ToLower(chainA),
		chainB:            strings.ToLower(chainB),
		chainAServiceName: serviceA,
		chainBServiceName: serviceB,
		bridge:            strconv.FormatBool(bridge),
	}
}

func (c *Chains) areChainsIcon() bool {
	return (c.chainA == "icon" && c.chainB == "icon")
}

func (chains *Chains) getParams() string {
	return fmt.Sprintf(`{"args":{"links": {"src": "%s", "dst": "%s"},"bridge":"%s"}}`, chains.chainA, chains.chainB, chains.bridge)
}
func (chains *Chains) getIbcRelayParams() string {

	return fmt.Sprintf(`{"args":{"links": {"src": "%s", "dst": "%s"}, "src_config":{"data":{}}, "dst_config":{"data":{}}}}`, chains.chainA, chains.chainB)
}

func (chains *Chains) getServicesResponse() (string, string, error) {

	serviceConfig, err := common.ReadServiceJsonFile()

	if err != nil {
		return "", "", err
	}

	chainAServiceResponse, OK := serviceConfig[chains.chainAServiceName]
	if !OK {
		return "", "", fmt.Errorf("service name not found")
	}
	chainBServiceResponse, OK := serviceConfig[chains.chainBServiceName]
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

func (chains *Chains) checkForBtpSupportedChains() error {
	if !slices.Contains(supportedChainsForBtp, chains.chainA) {
		return fmt.Errorf("invalid Chain: %s", chains.chainA)
	}
	if !slices.Contains(supportedChainsForBtp, chains.chainB) {
		return fmt.Errorf("invalid Chain: %s", chains.chainB)
	}
	return nil
}

func (chains *Chains) checkForIbcSupportedChains() error {
	if !slices.Contains(supportedChainsForIbc, chains.chainA) {
		return fmt.Errorf("invalid Chain: %s", chains.chainA)
	}
	if !slices.Contains(supportedChainsForIbc, chains.chainB) {
		return fmt.Errorf("invalid Chain: %s", chains.chainB)
	}
	return nil
}
