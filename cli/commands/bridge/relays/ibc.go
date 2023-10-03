package relays

import (
	"fmt"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

func IbcRelayCmd(diveContext *common.DiveContext) *cobra.Command {

	ibcRelayCommand := &cobra.Command{
		Use:   "ibc",
		Short: "Start connection between Cosmos based chainA and ChainB and initiate communication between them",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {

			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())
			diveContext.InitKurtosisContext()

			enclaveCtx, err := diveContext.GetEnclaveContext()

			if err != nil {
				diveContext.Error(err.Error())
			}

			result := startIbcRelay(diveContext, enclaveCtx)

			err = common.WriteToFile(result)

			if err != nil {
				diveContext.Error(err.Error())
			}

			diveContext.StopSpinner(fmt.Sprintf("IBC Setup Completed between %s and %s. Please find service details in current working directory(dive.json)", chainA, chainB))
		},
	}

	ibcRelayCommand.Flags().StringVar(&chainA, "chainA", "", "Mention Name of Supported Chain")
	ibcRelayCommand.Flags().StringVar(&chainB, "chainB", "", "Mention Name of Supported Chain")

	ibcRelayCommand.Flags().StringVar(&serviceA, "chainAServiceName", "", "Service Name of Chain A from services.json")
	ibcRelayCommand.Flags().StringVar(&serviceB, "chainBServiceName", "", "Service Name of Chain B from services.json")
	ibcRelayCommand.MarkFlagRequired("chainA")
	ibcRelayCommand.MarkFlagRequired("chainB")

	return ibcRelayCommand
}

func startIbcRelay(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext) string {
	diveContext.StartSpinner(" Starting IBC Setup")
	chains := initChains(chainA, chainB, serviceA, serviceB, false)
	var starlarkExecutionResponse string
	var err error

	err = chains.checkForIbcSupportedChains()
	if err != nil {
		diveContext.FatalError(err.Error(), fmt.Sprintf("Supported chains are %v", supportedChainsForIbc))
	}

	if chains.chainAServiceName != "" && chains.chainBServiceName != "" {

		srcChainServiceResponse, dstChainServiceResponse, err := chains.getServicesResponse()

		if err != nil {
			diveContext.FatalError("Failed To read ServiceFile", err.Error())
		}
		starlarkExecutionResponse, err = setupIbcRelayforAlreadyRunningCosmosChain(diveContext, enclaveContext, chains.chainA, chains.chainB, srcChainServiceResponse, dstChainServiceResponse)

		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}

	} else {
		starlarkExecutionResponse, err = startCosmosChainsAndSetupIbcRelay(diveContext, enclaveContext, chains)

		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}

	}

	if chainA == "icon" {
		startIbcRelayIconToCosmos(diveContext, enclaveContext, common.RelayServiceNameIconToCosmos)
	}

	return starlarkExecutionResponse
}

func startIbcRelayIconToCosmos(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext, serviceName string) (string, error) {
	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)

	executionData, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, "services/bridges/ibc/src/bridge.star", "start_relay", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	executionSerializedData, services, _, err := diveContext.GetSerializedData(executionData)

	if err != nil {
		diveContext.StopServices(services)
		return "", err

	}

	return executionSerializedData, nil
}

func startCosmosChainsAndSetupIbcRelay(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, chains *Chains) (string, error) {

	params := chains.getIbcRelayParams()

	executionResult, err := runStarlarkPackage(diveContext, enclaveCtx, params, "run_cosmos_ibc_setup")

	if err != nil {
		return "", err
	}

	return executionResult, nil
}

func setupIbcRelayforAlreadyRunningCosmosChain(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, chainA, chainB, chainAServiceResponse, chainBServiceResponse string) (string, error) {

	params := fmt.Sprintf(`{"src_chain_config":%s,"dst_chain_config":%s, "args":{"links": {"src": "%s", "dst": "%s"}, "src_config":{"data":{}}, "dst_config":{"data":{}}}}`, chainAServiceResponse, chainBServiceResponse, chainA, chainB)

	executionResult, err := runStarlarkPackage(diveContext, enclaveCtx, params, "run_cosmos_ibc_relay_for_already_running_chains")

	if err != nil {
		return "", err
	}

	return executionResult, nil
}

func runStarlarkPackage(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext, params, functionName string) (string, error) {
	executionData, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, functionName, params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		return "", err
	}

	executionSerializedData, services, skippedInstructions, err := diveContext.GetSerializedData(executionData)

	if err != nil {
		diveContext.StopServices(services)
		return "", err

	}

	diveContext.CheckInstructionSkipped(skippedInstructions, "Instruction Executed Already")

	return executionSerializedData, nil
}
