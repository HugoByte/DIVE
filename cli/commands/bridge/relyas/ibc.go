package relyas

import (
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
			diveContext.InitKurtosisContext()

			enclaveCtx, err := diveContext.GetEnclaveContext()

			if err != nil {
				diveContext.Error(err.Error())
			}

			stratIbcRelay(diveContext, enclaveCtx)
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

func stratIbcRelay(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext) {
	chains := initChains(chainA, chainB, serviceA, serviceB, false)
	result, err := startCosmosChainsAndSetupIbcRelay(diveContext, enclaveContext, chains)

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}
	common.WriteToFile(result)
}

func startCosmosChainsAndSetupIbcRelay(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, chains *Chains) (string, error) {

	params := chains.getIbcRelayParams()
	data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, "run_cosmos_ibc_setup", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		return "", err
	}

	responseData, services, skippedInstructions, err := diveContext.GetSerializedData(data)

	if err != nil {
		diveContext.StopServices(services)
		return "", err

	}

	diveContext.CheckInstructionSkipped(skippedInstructions, "Instruction Executed Already")

	return responseData, nil
}

func setupIbcRelayforAlreadyRunningCosmosChain(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext) (string, error) {
	return "", nil
}
