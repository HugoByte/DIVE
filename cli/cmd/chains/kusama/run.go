package kusama

import (
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"

	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
)

const (
	localChain       = "local"
	configsDirectory = "/home/riya/polakadot-kurtosis-package/parachain/static_files/configs"
)

func RunKusama(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return nil, common.WrapMessageToError(err, "Kusama Run Failed")
	}

	var serviceConfig = &utils.PolkadotServiceConfig{}

	err = common.LoadConfig(cli, serviceConfig, configFilePath)
	if err != nil {
		return nil, err
	}

	configureService(serviceConfig)

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
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

func configureService(serviceConfig *utils.PolkadotServiceConfig) {
	if paraChain != "" {
		serviceConfig.Para[0].Name = paraChain
	}

	if network != "" {
		serviceConfig.ChainType = network
		if network == "testnet" {
			serviceConfig.RelayChain.Name = "rococo"
		} else if network == "mainnet" {
			serviceConfig.RelayChain.Name = "kusama"
		}
	}

	if explorer {
		serviceConfig.Explorer = true
	}

	if metrics {
		configureMetrics(serviceConfig)
	}
}

func configureMetrics(serviceConfig *utils.PolkadotServiceConfig) {
	for i := range serviceConfig.RelayChain.Nodes {
		serviceConfig.RelayChain.Nodes[i].Prometheus = true
	}
	for i := range serviceConfig.Para[0].Nodes {
		serviceConfig.Para[0].Nodes[i].Prometheus = true
	}
}

func getKusamaRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if serviceConfig.Para[0].Name != "" {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotDefaultNodeSetupScript, runKusamaFunctionName)
	} else {
		if serviceConfig.ChainType == localChain {
			enclaveContext.UploadFiles(configsDirectory, "configs")
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayLocal)
		} else {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayTestnetMainet)
		}

	}
}
