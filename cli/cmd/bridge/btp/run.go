package btp

import (
	"fmt"
	"strconv"

	"github.com/hugobyte/dive-core/cli/cmd/bridge/utils"
	"github.com/hugobyte/dive-core/cli/cmd/chains/eth"
	"github.com/hugobyte/dive-core/cli/cmd/chains/hardhat"
	"github.com/hugobyte/dive-core/cli/cmd/chains/icon"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

var runChain = map[string]func(cli *common.Cli) (*common.DiveServiceResponse, error){
	"icon": func(cli *common.Cli) (*common.DiveServiceResponse, error) {
		nodeResponse, err := icon.RunIconNode(cli)
		if err != nil {
			return nil, err
		}
		params := icon.GetDecentralizeParams(nodeResponse.ServiceName, nodeResponse.PrivateEndpoint, nodeResponse.KeystorePath, nodeResponse.KeyPassword, nodeResponse.NetworkId)

		cli.Spinner().StartWithMessage("Starting Decentralization", "green")
		err = icon.RunDecentralization(cli, params)
		if err != nil {
			return nil, err
		}

		return nodeResponse, nil

	},
	"eth": func(cli *common.Cli) (*common.DiveServiceResponse, error) {
		return eth.RunEth(cli)

	},
	"hardhat": func(cli *common.Cli) (*common.DiveServiceResponse, error) {

		return hardhat.RunHardhat(cli)
	},
}

func RunBtpSetup(cli *common.Cli, chains *utils.Chains, bridge bool) (string, error) {
	var starlarkExecutionResponse string

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return "", common.WrapMessageToError(err, "BTP Setup Run Failed")
	}
	if chains.AreChainsIcon() {
		starlarkExecutionResponse, err = runBtpSetupWhenChainsAreIcon(cli, chains, enclaveContext, bridge)
		if err != nil {
			return "", err
		}
	} else {
		starlarkExecutionResponse, err = runBtpSetupWhenChainsAreNotIcon(cli, enclaveContext, chains, bridge)
		if err != nil {
			return "", err
		}
	}

	return starlarkExecutionResponse, nil

}

func runBtpSetupWhenChainsAreIcon(cli *common.Cli, chains *utils.Chains, enclaveContext *enclaves.EnclaveContext, bridge bool) (string, error) {

	if chains.ChainAServiceName != "" && chains.ChainBServiceName == "" {
		srcChainServiceResponse, dstChainServiceResponse, err := chains.GetServicesResponse(cli)
		if err != nil {
			return "", common.WrapMessageToError(err, "BTP Setup run Failed For Icon Chains")
		}
		response, err := runBtpSetupForAlreadyRunningNodes(cli, enclaveContext, runBridgeIcon2icon, chains.ChainA, chains.ChainB, chains.ChainAServiceName, chains.ChainBServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)
		if err != nil {
			return "", common.WrapMessageToError(err, "BTP Setup run Failed For Icon Chains")
		}
		return response, nil
	} else {
		params := chains.GetParams()

		response, err := runBtpSetupByRunningNodes(cli, enclaveContext, params)
		if err != nil {
			return "", common.WrapMessageToError(err, "BTP Setup run Failed For Icon Chains")
		}
		return response, nil
	}

}

func runBtpSetupWhenChainsAreNotIcon(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, chains *utils.Chains, bridge bool) (string, error) {

	if chains.ChainAServiceName != "" && chains.ChainBServiceName != "" {
		var response string
		chainAServiceResponse, chainBServiceResponse, err := chains.GetServicesResponse(cli)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
		}
		if chains.ChainB == "icon" {
			response, err = runBtpSetupForAlreadyRunningNodes(cli, enclaveContext, runBridgeIcon2EthHardhat, chains.ChainB, chains.ChainA, chains.ChainBServiceName, chains.ChainAServiceName, bridge, chainBServiceResponse, chainAServiceResponse)
			if err != nil {
				return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
			}
		} else {
			response, err = runBtpSetupForAlreadyRunningNodes(cli, enclaveContext, runBridgeIcon2EthHardhat, chains.ChainA, chains.ChainB, chains.ChainAServiceName, chains.ChainBServiceName, bridge, chainAServiceResponse, chainBServiceResponse)
			if err != nil {
				return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
			}
		}

		return response, nil
	} else if (chains.ChainAServiceName == "" && chains.ChainBServiceName != "") || (chains.ChainAServiceName != "" && chains.ChainBServiceName == "") {
		response, err := runBtpSetupWhenSingleChainRunning(cli, enclaveContext, chains, bridge)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
		}
		return response, nil
	} else {
		params := chains.GetParams()
		response, err := runBtpSetupByRunningNodes(cli, enclaveContext, params)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
		}
		return response, nil
	}

}

func runBtpSetupWhenSingleChainRunning(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, chains *utils.Chains, bridge bool) (string, error) {
	var chainAServiceResponse, chainBServiceResponse, chainAServiceName, chainBServiceName, response string
	var services = common.Services{}
	err := cli.FileHandler().ReadJson("services.json", &services)

	if err != nil {
		return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
	}

	if chains.ChainAServiceName == "" {
		serviceResponse, OK := services[chains.ChainBServiceName]
		if !OK {
			return "", common.WrapMessageToError(common.ErrDataUnMarshall, fmt.Sprint("service name not found:", chains.ChainBServiceName))
		}
		chainBServiceName = serviceResponse.ServiceName

		chainBServiceResponse, err = serviceResponse.EncodeToString()
		if err != nil {
			return "", common.WrapMessageToError(err, "BTP Setup Failed")
		}
		responseData, err := runChain[chains.ChainA](cli)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed Due To ChainA %s Run Failed ", chains.ChainA))
		}
		chainAServiceName = responseData.ServiceName
		chainAServiceResponse, err = responseData.EncodeToString()
		if err != nil {
			return "", err
		}

	} else if chains.ChainBServiceName == "" {
		serviceResponse, OK := services[chains.ChainAServiceName]
		if !OK {
			return "", common.WrapMessageToError(common.ErrDataUnMarshall, fmt.Sprint("service name not found:", chains.ChainAServiceName))
		}
		chainAServiceName = serviceResponse.ServiceName
		chainAServiceResponse, err = serviceResponse.EncodeToString()
		if err != nil {
			return "", common.WrapMessageToError(err, "BTP Setup Failed")
		}
		responseData, err := runChain[chains.ChainB](cli)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed Due To ChainB %s Run Failed ", chains.ChainB))
		}
		chainBServiceName = responseData.ServiceName
		chainBServiceResponse, err = responseData.EncodeToString()
		if err != nil {
			return "", err
		}

	}

	if chains.ChainB == "icon" {
		response, err = runBtpSetupForAlreadyRunningNodes(cli, enclaveContext, runBridgeIcon2EthHardhat, chains.ChainB, chains.ChainA, chainBServiceName, chainAServiceName, bridge, chainBServiceResponse, chainAServiceResponse)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
		}
	} else {
		response, err = runBtpSetupForAlreadyRunningNodes(cli, enclaveContext, runBridgeIcon2EthHardhat, chains.ChainA, chains.ChainB, chainAServiceName, chainBServiceName, bridge, chainAServiceResponse, chainBServiceResponse)
		if err != nil {
			return "", common.WrapMessageToError(err, fmt.Sprintf("BTP Setup Failed For ChainA %s and ChainB %s", chains.ChainA, chains.ChainB))
		}
	}

	return response, nil
}

func runBtpSetupByRunningNodes(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, params string) (string, error) {

	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveBridgeBtpScript, bridgeMainFunction)
	executionData, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return "", common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}
	executionSerializedData, services, skippedInstructions, err := common.GetSerializedData(cli, executionData)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return "", common.WrapMessageToError(errRemove, "BTP Setup Run Failed")
		}

		return "", common.WrapMessageToError(err, "BTP Setup Run Failed")

	}
	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return "", common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}
	return executionSerializedData, nil

}

func runBtpSetupForAlreadyRunningNodes(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, mainFunctionName string, srcChain string, dstChain string, srcChainServiceName string, dstChainServiceName string, bridge bool, srcChainServiceResponse string, dstChainServiceResponse string) (string, error) {

	params := fmt.Sprintf(`{"src_chain":"%s","dst_chain":"%s", "src_chain_config":%s, "dst_chain_config":%s, "bridge":%s}`, chainA, chainB, srcChainServiceResponse, dstChainServiceResponse, strconv.FormatBool(bridge))
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveBridgeBtpScript, mainFunctionName)
	executionData, _, err := enclaveCtx.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return "", common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}
	executionSerializedData, services, skippedInstructions, err := common.GetSerializedData(cli, executionData)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return "", common.WrapMessageToError(errRemove, "BTP Setup Run Failed")
		}

		return "", common.WrapMessageToError(err, "BTP Setup Run Failed")

	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return "", common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	return executionSerializedData, nil

}
