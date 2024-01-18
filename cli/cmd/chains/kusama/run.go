package kusama

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"

	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
)

const (
	localChain   = "local"
	polkadotJUrl = "http://127.0.0.1/?rpc=ws%3A%2F%2F127.0.0.1%3A91941#/explorer"
	PolkadotJsServiceName = "polkadot-js-explorer"
)

func RunKusama(cli *common.Cli) (*common.DiveMultipleServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return nil, common.WrapMessageToError(err, "Failed to retrieve the enclave context for Kusama.")
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

	err = serviceConfig.ValidateConfig()
	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrInvalidConfig, err.Error())
	}

	for _, paraChain := range serviceConfig.Para {
		if !slices.Contains(kusamaParachains, paraChain.Name) {
			return nil, common.WrapMessageToErrorf(common.ErrInvalidConfig, "Invalid Parachain - Parachain %s is not Supported for Kusama", paraChain.Name)
		}
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	err = uploadFiles(cli, enclaveContext)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Failed to upload the configuration files.")
	}

	result := &common.DiveMultipleServiceResponse{}

	if serviceConfig.RelayChain.Name == "" {
		result, err = startParaChains(cli, enclaveContext, serviceConfig, encodedServiceConfigDataString, "")
		if err != nil {
			return nil, err
		}
	} else {
		result, err = startRelayAndParaChain(cli, enclaveContext, serviceConfig, encodedServiceConfigDataString)
		if err != nil {
			return nil, err
		}
	}

	return result, nil

}

func startRelayAndParaChain(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceConfig *utils.PolkadotServiceConfig, para string) (*common.DiveMultipleServiceResponse, error) {

	kusamaResponseData := &common.DiveMultipleServiceResponse{}
	paraResult := &common.DiveMultipleServiceResponse{}
	finalResult := &common.DiveMultipleServiceResponse{}
	explorerResult := &common.DiveMultipleServiceResponse{}
	metricsResult := &common.DiveMultipleServiceResponse{}

	param, err := serviceConfig.GetParamsForRelay()
	if err != nil {
		return nil, err
	}

	runConfig := getKusamaRunConfig(serviceConfig, enclaveContext, param)
	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Kusama relaychain run failed. Failed to clean up services.")
		}
		return nil, common.WrapMessageToError(err, "Kusama relaychain run failed. Failed to serialize the response data.")
	}

	result, err := kusamaResponseData.Decode([]byte(responseData))
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Kusama relaychain run failed. Failed to clean up services.")
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Kusama relaychain run failed. Failed to decode reponse data.")
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
			finalResult = result.ConcatenateDiveResults(paraResult)

		} else {
			return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Kusama is already Running.")
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
			finalResult = result.ConcatenateDiveResults(paraResult)
		}
	}

	if serviceConfig.HasPrometheus() {
		metricsResult, err = startMetrics(cli, enclaveContext, finalResult)
		finalResult = finalResult.ConcatenateDiveResults(metricsResult)
		if err != nil {
			return nil, err
		}
	}

	if serviceConfig.Explorer {
		publicEndpoint := "ws://127.0.0.1:9944"
		
		var isExplorerRunning bool

		allEclaveServices, err := cli.Context().GetAllEnlavesServices()
		if err != nil {
			return nil, err
		}

		for _, enclave := range allEclaveServices {
			for serviceName, _ := range enclave {
				if serviceName ==  PolkadotJsServiceName {
					isExplorerRunning = true
				}
			} 
		}

		if len(finalResult.Dive) > 0 {
			for _, serviceResponse := range finalResult.Dive {
				publicEndpoint = serviceResponse.PublicEndpoint
				break
			}
		}
		
		if !isExplorerRunning {
			explorerResult, err = startExplorer(cli, enclaveContext, publicEndpoint)
			if err != nil {
				return nil, err
			}
			finalResult = finalResult.ConcatenateDiveResults(explorerResult)
		} else {
			cli.Logger().Info("Explorer service is already running.")	
		}
		
		url := updatePort(polkadotJUrl, extractPort(publicEndpoint))
		cli.Logger().Info("Redirecting to Polkadote explorer UI...")
		if err := common.OpenFile(url); err != nil {
			cli.Logger().Fatalf(common.CodeOf(err), "Failed to open HugoByte Polkadot explorer UI with error %v", err)
		}
	}


	return finalResult, nil
}

func startParaChains(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceConfig *utils.PolkadotServiceConfig, para string, ipAddress string) (*common.DiveMultipleServiceResponse, error) {
	paraResult := &common.DiveMultipleServiceResponse{}
	allParaResult := &common.DiveMultipleServiceResponse{}

	if serviceConfig.ChainType == localChain {
		paraNodeList := utils.ParaNodeConfigList(serviceConfig.Para)
		var paraChains string
		paraChains, err := paraNodeList.EncodeToString()
		if err != nil {
			return nil, err
		}
		param := fmt.Sprintf(`{"chain_type":"%s", "parachains": %s, "relay_chain_ip": "%s"}`, serviceConfig.ChainType, paraChains, ipAddress)
		runParaConfig := getParaRunConfig(serviceConfig, enclaveContext, param)
		paraResult, err = startService(cli, enclaveContext, runParaConfig, "Parachain")
		if err != nil {
			return nil, err
		}
		allParaResult = allParaResult.ConcatenateDiveResults(paraResult)
	} else {
		for _, paraNode := range serviceConfig.Para {
			paraChainConfig, err := paraNode.EncodeToString()
			if err != nil {
				return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
			}
			param := fmt.Sprintf(`{"chain_type": "%s", "relaychain_name": "%s", "parachain":%s}`, serviceConfig.ChainType, serviceConfig.RelayChain.Name, paraChainConfig)
			runParaConfig := getParaRunConfig(serviceConfig, enclaveContext, param)
			paraResult, err = startService(cli, enclaveContext, runParaConfig, "Parachain")
			if err != nil {
				return nil, err
			}
			allParaResult = allParaResult.ConcatenateDiveResults(paraResult)
		}
	}

	return allParaResult, nil
}

func configureService(serviceConfig *utils.PolkadotServiceConfig) error {

	if len(paraChain) != 0 {
		serviceConfig.Para = []utils.ParaNodeConfig{}
		for _, value := range paraChain {
			if value != "" {
				serviceConfig.Para = append(serviceConfig.Para, utils.ParaNodeConfig{
					Name: value,
					Nodes: []utils.NodeConfig{
						{Name: "alice", NodeType: "full", Prometheus: false},
					},
				})
			}
		}
	}

	if network != "" {
		serviceConfig.ChainType = network
		if network == "testnet" || network == "mainnet" {
			serviceConfig.ConfigureFullNodes(network)
		}
	}

	if explorer {
		serviceConfig.Explorer = true
	}

	if metrics {
		serviceConfig.ConfigureMetrics()
	}

	for i := range serviceConfig.RelayChain.Nodes {
		serviceConfig.RelayChain.Nodes[i].AssignPorts(serviceConfig.RelayChain.Nodes[i].Prometheus)
	}

	for _, paraChain := range serviceConfig.Para {
		for i := range paraChain.Nodes {
			paraChain.Nodes[i].AssignPorts(paraChain.Nodes[i].Prometheus)
		}
	}

	if noRelay && serviceConfig.ChainType == "local" {
		return common.WrapMessageToError(common.ErrInvalidFlag, "The '--no-relay' flag cannot be used with a 'local' network. This flag is only applicable for 'testnet' and 'mainnet' networks.")
	} else if noRelay && serviceConfig.ChainType != "local" {
		serviceConfig.RelayChain = utils.RelayChainConfig{}
	}

	if serviceConfig.ChainType == localChain && serviceConfig.RelayChain.Name == "" && len(serviceConfig.Para) != 0 {
		return common.WrapMessageToError(common.ErrEmptyFields, "Cannot start a Parachain in local without Relay Chain")
	}

	return nil
}

func flagCheck() error {

	if configFilePath != "" {
		if len(paraChain) != 0 || network != "" || explorer || metrics || noRelay {
			return common.WrapMessageToError(common.ErrInvalidFlag, "The '-c' flag does not allow additional flags.")
		}
	}

	if noRelay && (network == "testnet" || network == "mainnet") {
		if len(paraChain) == 0 {
			return common.WrapMessageToError(common.ErrMissingFlags, "The '-p' flag is required when using '--no-relay' flag. Please provide the '-p' flag with the parachain name.")
		}
	}
	return nil
}

func getKusamaRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if serviceConfig.ChainType == localChain {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayLocal)
	} else {
		return common.GetStarlarkRunConfig(para, common.DivePolkadotRelayNodeSetupScript, runKusamaRelayTestnetMainet)
	}
}

func getParaRunConfig(serviceConfig *utils.PolkadotServiceConfig, enclaveContext *enclaves.EnclaveContext, para string) *starlark_run_config.StarlarkRunConfig {
	if len(serviceConfig.Para) != 0 && serviceConfig.Para[0].Name != "" {
		if serviceConfig.ChainType == localChain {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotParachainNodeSetup, runKusamaParaLocalFunctionName)

		} else {
			return common.GetStarlarkRunConfig(para, common.DivePolkadotParachainNodeSetup, runKusamaParaTestMainFunctionName)

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
				return "", fmt.Errorf("failed to retrieve the UUID of the enclave")
			}
			serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)

			err = cli.FileHandler().ReadJson(serviceFileName, &services)
			if err != nil {
				return "", err
			}

			chainServiceName := fmt.Sprintf("rococo-local-%s", nodename)
			chainServiceResponse, OK := services[chainServiceName]
			if !OK {
				return "", fmt.Errorf("service name '%s' not found", chainServiceName)
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

func startExplorer(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, publicEndpoint string) (*common.DiveMultipleServiceResponse, error) {
	para := fmt.Sprintf(`{"ws_url":"%s"}`, publicEndpoint)

	runConfig := common.GetStarlarkRunConfig(para, common.DivePolkaDotExplorerPath, runKusamaExplorer)
	explorerResponseData, err := startService(cli, enclaveCtx, runConfig, "Explorer")
	if err != nil {
		return nil, err
	}
	return explorerResponseData, nil
}

func startMetrics(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, final_result *common.DiveMultipleServiceResponse) (*common.DiveMultipleServiceResponse, error) {
	service_details, err := final_result.EncodeToString()
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	paraPrometheus := fmt.Sprintf(`{"service_details":%s}`, service_details)

	runConfigPrometheus := common.GetStarlarkRunConfig(paraPrometheus, common.DivePolkaDotPrometheusPath, runKusamaPrometheus)
	prometheusResult, err := startService(cli, enclaveCtx, runConfigPrometheus, "Prometheus")
	if err != nil {
		return nil, err
	}

	runConfigGrafana := common.GetStarlarkRunConfig(`{}`, common.DivePolkaDotGrafanaPath, runKusamaGrafana)
	grafanaResult, err := startService(cli, enclaveCtx, runConfigGrafana, "Grafana")
	if err != nil {
		return nil, err
	}

	result := prometheusResult.ConcatenateDiveResults(grafanaResult)
	return result, nil
}

func startService(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, runConfig *starlark_run_config.StarlarkRunConfig, serviceName string) (*common.DiveMultipleServiceResponse, error) {
	starlarkResponseData := &common.DiveMultipleServiceResponse{}

	response, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.PolkadotRemotePackagePath, runConfig)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToErrorf(errRemove, "%s Run Failed. Cleanup of services failed.", serviceName)
		}
		return nil, common.WrapMessageToErrorf(err, "%s Run Failed. Failed to serilize response data.", serviceName)
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkResponse, "%s is already running.", serviceName)
	}

	result, err := starlarkResponseData.Decode([]byte(responseData))

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToErrorf(errRemove, "%s Run Failed. Cleanup of services failed.", serviceName)
		}
		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s Run Failed. Failed to decode response data.", serviceName)
	}

	return result, nil
}

func updatePort(url string, newPort string) string {
	re := regexp.MustCompile(`ws%3A%2F%2F127\.0\.0\.1%3A(\d+)`)
	newUrl := re.ReplaceAllString(url, fmt.Sprintf("ws%%3A%%2F%%2F127.0.0.1%%3A%s", newPort))
	return newUrl
}

func extractPort(url string) string {
	re := regexp.MustCompile(`ws://127\.0\.0\.1:(\d+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
