package types

import (
	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

func NewHardhatCmd(diveContext *common.DiveContext) *cobra.Command {

	var ethCmd = &cobra.Command{
		Use:   "hardhat",
		Short: "Build, initialize and start a hardhat node.",
		Long: `The command starts an hardhat node, initiating the process of setting up and launching a local hardhat network. 
It establishes a connection to the hardhat network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			common.ValidateCmdArgs(args, cmd.UsageString())

			data := RunHardhatNode(diveContext)

			diveContext.SetSpinnerMessage("Execution Completed")

			err := data.WriteDiveResponse()
			if err != nil {
				diveContext.FatalError("Failed To Write To File", err.Error())
			}
			diveContext.StopSpinner("Hardhat Node Started. Please find service details in current working directory(dive.json)")
		},
	}

	return ethCmd

}

func RunHardhatNode(diveContext *common.DiveContext) *common.DiveserviceResponse {

	diveContext.InitKurtosisContext()

	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}
	diveContext.StartSpinner(" Starting Hardhat Node")
	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveEthHardhatNodeScript, "start_hardhat_node", "{}", common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	responseData, services, skippedInstructions, err := diveContext.GetSerializedData(data)
	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	diveContext.CheckInstructionSkipped(skippedInstructions, common.DiveHardhatNodeAlreadyRuning)
	hardhatResponseData := &common.DiveserviceResponse{}

	result, err := hardhatResponseData.Decode([]byte(responseData))

	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Failed To Unmarshall Data", err.Error())
	}

	return result

}
