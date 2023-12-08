package ibc

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/cmd/bridge/utils"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

func RunIbcRelay(cli *common.Cli) (string, error) {
	var starlarkExecutionResponse string
	chains := utils.InitChains(chainA, chainB, serviceA, serviceB, false)

	err := chains.CheckForIbcSupportedChains()

	if err != nil {

		return "", common.WrapMessageToError(common.ErrInvalidChain, err.Error())
	}

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return "", common.WrapMessageToError(err, "IBC Setup Failed")
	}

	if chains.CheckChainServiceNamesEmpty() {
		srcChainServiceResponse, dstChainServiceResponse, err := chains.GetServicesResponse(cli)
		if err != nil {
			return "", err
		}
		starlarkExecutionResponse, err = setupIbcRelayforAlreadyRunningCosmosChain(cli, enclaveContext, chains.ChainA, chains.ChainB, srcChainServiceResponse, dstChainServiceResponse)

		if err != nil {
			return "", err
		}

	} else {
		starlarkExecutionResponse, err = startCosmosChainsAndSetupIbcRelay(cli, enclaveContext, chains)
		if err != nil {
			return "", err
		}
	}

	if chainA == "icon" {
		_, err := startIbcRelayIconToCosmos(cli, enclaveContext, common.RelayServiceNameIconToCosmos)
		if err != nil {
			return "", err
		}
	}

	return starlarkExecutionResponse, nil
}

func startIbcRelayIconToCosmos(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, serviceName string) (string, error) {
	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
	starlarkConfig := common.GetStarlarkRunConfig(params, "services/bridges/ibc/src/bridge.star", "start_relay")
	executionData, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return "", common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	executionSerializedData, services, _, err := common.GetSerializedData(cli, executionData)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return "", common.WrapMessageToError(errRemove, "IBC Setup Run Failed")
		}

		return "", common.WrapMessageToError(err, "IBC Setup Run Failed")

	}

	return executionSerializedData, nil
}

func startCosmosChainsAndSetupIbcRelay(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, chains *utils.Chains) (string, error) {

	params := chains.GetIbcRelayParams()

	executionResult, err := runStarlarkPackage(cli, enclaveCtx, params, "run_cosmos_ibc_setup")

	if err != nil {
		return "", err
	}

	return executionResult, nil
}

func setupIbcRelayforAlreadyRunningCosmosChain(cli *common.Cli, enclaveCtx *enclaves.EnclaveContext, chainA, chainB, chainAServiceResponse, chainBServiceResponse string) (string, error) {

	params := fmt.Sprintf(`{"src_chain":"%s","dst_chain":"%s", "src_chain_config":%s, "dst_chain_config":%s}`, chainA, chainB, chainAServiceResponse, chainBServiceResponse)

	executionResult, err := runStarlarkPackage(cli, enclaveCtx, params, "run_cosmos_ibc_relay_for_already_running_chains")

	if err != nil {
		return "", err
	}

	return executionResult, nil
}

func runStarlarkPackage(cli *common.Cli, enclaveContext *enclaves.EnclaveContext, params, functionName string) (string, error) {
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveBridgeIbcScript, functionName)
	executionData, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return "", err
	}

	executionSerializedData, services, skippedInstructions, err := common.GetSerializedData(cli, executionData)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return "", common.WrapMessageToError(errRemove, "IBC Setup Run Failed")
		}

		return "", common.WrapMessageToError(err, "IBC Setup Run Failed")

	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return "", common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	return executionSerializedData, nil
}
