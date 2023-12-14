package polkadot

import (
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"

	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
)

const (
	localChain       = "local"
	configsDirectory = "/Users/abhishekharde/Desktop/hugobyte/dive-packages/services/polkadot/parachain/static_files/configs"
)

func RunPolkadot(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {
	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Polkadot Run Failed")
	}
	var serviceConfig = &utils.PolkadotServiceConfig{}

	err = common.LoadConfig(cli, serviceConfig, configFilePath)

	if err != nil {
		return nil, err
	}

	configureService(serviceConfig)

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	para := fmt.Sprintf(`{"args": %s}`, encodedServiceConfigDataString)

	fmt.Print(para)

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	runConfig := getPolkadotRunConfig(serviceConfig, enclaveContext, para)

	response, _, err := enclaveContext.RunStarlarkPackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Polkadot Run Failed ")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Polkadot Running")
	}

	polkadotResponseData := &common.DiveMultipleServiceResponse{}
	result, err := polkadotResponseData.Decode([]byte(responseData))


	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Polkadot Run Failed ")

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
			serviceConfig.RelayChain.Name = "polkadot"
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
	for _, node := range append(serviceConfig.RelayChain.Nodes, serviceConfig.Para[0].Nodes...) {
		node.Prometheus = true
	}
}

func getPolkadotRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if serviceConfig.Para[0].Name != "" {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotDefaultNodeSetupScript, runPolkadotFunctionName)
	} else {
		if serviceConfig.ChainType == localChain {
			enclaveContext.UploadFiles(configsDirectory, "configs")
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runPolkadotRelayLocal)
		} else {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runPolkadotRelayTestnetMainet)
		}

	}
}
