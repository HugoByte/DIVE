package polkadot

import (
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"

	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
)

const (
	localChain = "local"
)

func RunPolkadot(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return nil, common.WrapMessageToError(err, "Failed to retrieve the enclave context for Polkadot.")
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
		return nil, common.WrapMessageToError(err, "Failed to upload the configuration files.")
	}

	result, err := startRelayAndParaChain(cli, enclaveContext, serviceConfig, encodedServiceConfigDataString)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func startRelayAndParaChain(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceConfig *utils.PolkadotServiceConfig, para string) (*common.DiveMultipleServiceResponse, error) {

	param := fmt.Sprintf(`{"args": %s}`, para)

	polkadotResponseData := &common.DiveMultipleServiceResponse{}
	paraResult := &common.DiveMultipleServiceResponse{}
	finalResult := &common.DiveMultipleServiceResponse{}
	explorerResult := &common.DiveMultipleServiceResponse{}
	metricsResult := &common.DiveMultipleServiceResponse{}

	runConfig := getPolkadotRunConfig(serviceConfig, enclaveContext, param)

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot relaychain run failed. Failed to clean up services.")
		}
		return nil, common.WrapMessageToError(err, "Polkadot relaychain run failed. Failed to serialize the response data.")
	}

	result, err := polkadotResponseData.Decode([]byte(responseData))
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Polkadot relaychain run failed. Failed to clean up services.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Polkadot relaychain run failed. Failed to decode reponse data.")
	}

	finalResult = result

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		if len(serviceConfig.Para) != 0 && serviceConfig.Para[0].Name != "" {
			ipAddress, err := getIPAddress(cli, serviceConfig, true, result)
			if err != nil {
				return nil, err
			}
			paraResult, err = startParaChains(cli, enclaveContext, serviceConfig, para, ipAddress)
			if err != nil {
				return nil, err
			}
			finalResult = concatenateDiveResults(result, paraResult)

		} else {
			return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Polkadot is already Running.")
		}
	} else {
		if len(serviceConfig.Para) != 0 && serviceConfig.Para[0].Name != "" {
			ipAddress, err := getIPAddress(cli, serviceConfig, false, result)
			if err != nil {
				return nil, err
			}
			paraResult, err = startParaChains(cli, enclaveContext, serviceConfig, para, ipAddress)
			if err != nil {
				return nil, err
			}
			finalResult = concatenateDiveResults(result, paraResult)
		}
	}

	if metrics {
		metricsResult, err = startMetrics(cli, enclaveContext, para, finalResult)
		finalResult = concatenateDiveResults(finalResult, metricsResult)
		if err != nil {
			return nil, err
		}
	}

	if explorer {
		explorerResult, err = startExplorer(cli, enclaveContext)
		if err != nil {
			return nil, err
		}
		finalResult = concatenateDiveResults(finalResult, explorerResult)
	}

	return finalResult, nil
}

func startParaChains(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceConfig *utils.PolkadotServiceConfig, para string, ipAddress string) (*common.DiveMultipleServiceResponse, error) {
	paraResult := &common.DiveMultipleServiceResponse{}
	var err error

	if serviceConfig.ChainType == localChain {
		param := fmt.Sprintf(`{"args": %s, "relay_chain_ip":"%s"}`, para, ipAddress)
		paraResult, err = runParaChain(cli, enclaveContext, serviceConfig, param)
		if err != nil {
			return nil, err
		}
	} else {
		for _, paraNode := range serviceConfig.Para {
			paraChainConfig, err := paraNode.EncodeToString()
			if err != nil {
				return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
			}
			param := fmt.Sprintf(`{"parachain":%s, "args":%s}`, paraChainConfig, para)
			paraResult, err = runParaChain(cli, enclaveContext, serviceConfig, param)
			if err != nil {
				return nil, err
			}
		}

	}

	return paraResult, nil
}

func runParaChain(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceConfig *utils.PolkadotServiceConfig, para string) (*common.DiveMultipleServiceResponse, error) {

	runParaConfig := getParaRunConfig(serviceConfig, enclaveContext, para)
	paraResponse, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runParaConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	paraResponseData, paraServices, skippedParaInstructions, err := common.GetSerializedData(cli, paraResponse)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(paraServices, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Failed to clean up services.")
		}
		return nil, common.WrapMessageToError(err, "Failed to serialize the response data.")
	}

	if cli.Context().CheckSkippedInstructions(skippedParaInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Parachain is already running.")
	}

	PolkadotParaResponseData := &common.DiveMultipleServiceResponse{}
	resultPara, err := PolkadotParaResponseData.Decode([]byte(paraResponseData))
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(paraServices, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Failed to clean up services.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Failed to decode reponse data.")

	}

	return resultPara, nil
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
		if serviceConfig.ChainType == "local" {
			return common.WrapMessageToError(common.ErrInvalidFlag, "The '--no-relay' flag cannot be used with a 'local' network. This flag is only applicable for 'testnet' and 'mainnet' networks.")
		} else {
			serviceConfig.RelayChain = utils.RelayChainConfig{}
		}
	}

	return nil
}

func configureFullNodes(serviceConfig *utils.PolkadotServiceConfig) {

	if network == "testnet" {
		serviceConfig.RelayChain.Name = "rococo"
	} else if network == "mainnet" {
		serviceConfig.RelayChain.Name = "polkadot"
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
	if len(serviceConfig.Para) != 0 {
		for i := range serviceConfig.Para[0].Nodes {
			serviceConfig.Para[0].Nodes[i].Prometheus = true
		}
	}
}

func flagCheck() error {

	if configFilePath != "" {
		if paraChain != "" || network != "" || explorer || metrics {
			return common.WrapMessageToError(common.ErrInvalidFlag, "The '-c' flag does not allow additional flags.")
		}
	}

	if noRelay && (network == "testnet" || network == "mainnet") {
		if paraChain == "" {
			return common.WrapMessageToError(common.ErrMissingFlags, "The '-p' flag is required when using '--no-relay' flag. Please provide the '-p' flag with the parachain name.")
		}
	}
	return nil
}

func getPolkadotRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if serviceConfig.ChainType == localChain {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runPolkadotRelayLocal)
	} else {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runPolkadotRelayTestnetMainet)
	}
}

func getParaRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if len(serviceConfig.Para) != 0 && serviceConfig.Para[0].Name != "" {
		if serviceConfig.ChainType == localChain {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotParachainNodeSetup, runPolkadotParaLocalFunctionName)

		} else {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotParachainNodeSetup, runPolkadotParaTestMainFunctionName)

		}
	}
	return nil
}

func getIPAddress(cli *common.Cli, serviceConfig *utils.PolkadotServiceConfig, relayReRun bool, result *common.DiveMultipleServiceResponse) (string, error) {
	var nodename string
	if serviceConfig.ChainType == localChain {
		if relayReRun {
			nodename = serviceConfig.RelayChain.Nodes[0].Name
			var services = common.Services{}

			shortUuid, err := cli.Context().GetShortUuid(common.EnclaveName)
			if err != nil {
				return "", fmt.Errorf("Failed to retrieve the UUID of the enclave.")
			}
			serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

			err = cli.FileHandler().ReadJson(serviceFileName, &services)
			if err != nil {
				return "", err
			}

			chainServiceName := fmt.Sprintf("rococo-local-%s", nodename)
			chainServiceResponse, OK := services[chainServiceName]
			if !OK {
				return "", fmt.Errorf("Service name '%s' not found", chainServiceName)
			}

			ipAddress := chainServiceResponse.IpAddress
			return ipAddress, nil
		} else {
			servicename := fmt.Sprintf("rococo-local-%s", serviceConfig.RelayChain.Nodes[0].Name)
			ipAddress := result.Dive[servicename].IpAddress
			return ipAddress, nil
		}
	}
	return "", nil
}

func uploadFiles(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext) error {
	runConfig := common.GetStarlarkRunConfig("{}", common.DivePolkaDotUtilsPath, runUploadFiles)
	_, err := enclaveCtx.RunStarlarkRemotePackageBlocking(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}
	return nil
}

func startExplorer(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext) (*common.DiveMultipleServiceResponse, error) {
	explorerResponseData := &common.DiveMultipleServiceResponse{}

	para := `{"ws_url":"ws://127.0.0.1:9944"}`
	runConfig := common.GetStarlarkRunConfig(para, common.DivePolkaDotExplorerPath, runPolkadotExplorer)
	explorerResponse, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, explorerResponse)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Explorer Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToError(err, "Explorer Run Failed. Failed to serilize response data.")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Explorer is already running.")
	}

	result, err := explorerResponseData.Decode([]byte(responseData))

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Explorer Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Explorer Run Failed. Failed to decode response data.")
	}

	return result, nil
}

func startMetrics(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, para string, final_result *common.DiveMultipleServiceResponse) (*common.DiveMultipleServiceResponse, error) {
	prometheus := &common.DiveMultipleServiceResponse{}
	grafana := &common.DiveMultipleServiceResponse{}

	service_details, err := final_result.EncodeToString()
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	paraPrometheus := fmt.Sprintf(`{"args":%s, "service_details":%s}`, para, service_details)
	runConfig := common.GetStarlarkRunConfig(paraPrometheus, common.DivePolkaDotPrometheusPath, runPolkadotPrometheus)
	prometheusResponse, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	prometheusResponseData, services, skippedInstructions, err := common.GetSerializedData(cli, prometheusResponse)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Prometheus Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToError(err, "Prometheus Run Failed. Failed to serilize response data.")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Prometheus is already Running.")
	}

	prometheusResult, err := prometheus.Decode([]byte(prometheusResponseData))
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Prometheus Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Prometheus Run Failed. Failed to decode reponse data.")
	}

	paraGrafana := `{"args":{}}`
	runConfigGrafana := common.GetStarlarkRunConfig(paraGrafana, common.DivePolkaDotGrafanaPath, runPolkadotGrafana)
	grafanaResponse, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfigGrafana)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	grafanaResponseData, services, skippedInstructions, err := common.GetSerializedData(cli, grafanaResponse)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Grafana Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToError(err, "Grafana Run Failed. Failed to serialize response data.")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Grafana is already running.")
	}

	grafanaResult, err := grafana.Decode([]byte(grafanaResponseData))
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Grafana Run Failed. Cleanup of services failed.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Grafana Run Failed. Failed to decode response data.")
	}

	result := concatenateDiveResults(prometheusResult, grafanaResult)
	return result, nil
}

func concatenateDiveResults(result1, result2 *common.DiveMultipleServiceResponse) *common.DiveMultipleServiceResponse {
	if result1 == nil {
		return result2
	} else if result2 == nil {
		return result1
	}

	concatenatedResult := &common.DiveMultipleServiceResponse{
		Dive: make(map[string]*common.DiveServiceResponse),
	}

	for key, value := range result1.Dive {
		concatenatedResult.Dive[key] = value
	}

	for key, value := range result2.Dive {
		concatenatedResult.Dive[key] = value
	}

	return concatenatedResult
}
