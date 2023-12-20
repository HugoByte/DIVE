package kusama

import (
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"

	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
)

const (
	localChain = "local"
)

func RunKusama(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return nil, common.WrapMessageToError(err, "Kusama Run Failed")
	}

	var serviceConfig = &utils.PolkadotServiceConfig{}

	err = flagCheck()
	if err != nil {
		return nil, err
	}

	err = common.LoadConfig(cli, serviceConfig, configFilePath)
	if err != nil {
		return nil, err
	}

	err = configureService(serviceConfig)
	if err != nil {
		return nil, err
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	err = uploadFiles(cli, enclaveContext)
	if err != nil {
		return nil, err
	}

	para := fmt.Sprintf(`{"args": %s}`, encodedServiceConfigDataString)
	runConfig := getKusamaRunConfig(serviceConfig, enclaveContext, para)

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)
	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Kusama Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Kusama Run Failed ")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Kusama already Running")
	}

	KusamaResponseData := &common.DiveMultipleServiceResponse{}

	result, err := KusamaResponseData.Decode([]byte(responseData))
	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Kusama Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Kusama Run Failed ")

	}

	return result, nil
}

func configureService(serviceConfig *utils.PolkadotServiceConfig) error {
	if paraChain != "" {
		serviceConfig.Para = []utils.ParaNodeConfig{}
		serviceConfig.Para = append(serviceConfig.Para, utils.ParaNodeConfig{
			Name: paraChain,
			Nodes: []utils.NodeConfig{
				{Name: "alice", NodeType: "full", Prometheus: false},
			},
		})
	}

	if network != "" {
		serviceConfig.ChainType = network
		if network == "testnet" || network == "mainnet" {
			configureFullNodes(serviceConfig)
		}
	}

	if explorer {
		serviceConfig.Explorer = true
	}

	if metrics {
		configureMetrics(serviceConfig)
	}

	if noRelay && serviceConfig.ChainType == "local" {
		return common.WrapMessageToError(common.ErrInvalidFlag, "Cannot pass --no-relay flag with local network")
	} else if noRelay && serviceConfig.ChainType != "local" {
		serviceConfig.RelayChain = utils.RelayChainConfig{}
	}

	return nil
}

func configureFullNodes(serviceConfig *utils.PolkadotServiceConfig) {
	if network == "testnet" {
		serviceConfig.RelayChain.Name = "rococo"
	} else if network == "mainnet" {
		serviceConfig.RelayChain.Name = "kusama"
	}

	serviceConfig.RelayChain.Nodes = []utils.NodeConfig{}
	serviceConfig.RelayChain.Nodes = append(serviceConfig.RelayChain.Nodes, utils.NodeConfig{
		Name:       "alice",
		NodeType:   "full",
		Prometheus: false,
	})
}

func configureMetrics(serviceConfig *utils.PolkadotServiceConfig) {
	for i := range serviceConfig.RelayChain.Nodes {
		serviceConfig.RelayChain.Nodes[i].Prometheus = true
	}
	for i := range serviceConfig.Para[0].Nodes {
		serviceConfig.Para[0].Nodes[i].Prometheus = true
	}
}

func flagCheck() error {
	if configFilePath != "" {
		if paraChain != "" || network != "" || explorer || metrics {
			return common.WrapMessageToError(common.ErrInvalidFlag, "Additional Flags Found")
		}
	}

	if noRelay && (network == "testnet" || network == "mainnet") {
		if paraChain == "" {
			return common.WrapMessageToError(common.ErrMissingFlags, "Missing Parachain Flag")
		}
	}
	return nil
}

func getKusamaRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if len(serviceConfig.Para) != 0 && serviceConfig.Para[0].Name != "" {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotDefaultNodeSetupScript, runKusamaFunctionName)
	} else {
		if serviceConfig.ChainType == localChain {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayLocal)
		} else {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayTestnetMainet)
		}

	}
}

func uploadFiles(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext) error {
	runConfig := common.GetStarlarkRunConfig("{}", common.DivePolkaDotUtilsPath, runUploadFiles)
	_, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	return nil
}
