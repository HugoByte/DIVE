package types

import (
	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

func NewEthCmd(diveContext *common.DiveContext) *cobra.Command {

	var ethCmd = &cobra.Command{
		Use:   "eth",
		Short: "Build, initialize and start a eth node.",
		Long: `The command starts an Ethereum node, initiating the process of setting up and launching a local Ethereum network. 
It establishes a connection to the Ethereum network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			common.ValidateCmdArgs(args, cmd.UsageString())

			data, err := RunEthNode(diveContext)

			if err != nil {
				diveContext.FatalError("Fail to Start ETH Node", err.Error())
			}
			diveContext.SetSpinnerMessage("Execution Completed")
			err = data.WriteDiveResponse(diveContext)
			if err != nil {
				diveContext.FatalError("Failed To Write To File", err.Error())
			}
			diveContext.StopSpinner("ETH Node Started. Please find service details in current working directory(dive.json)")
		},
	}

	return ethCmd

}

func RunEthNode(diveContext *common.DiveContext) (*common.DiveserviceResponse, error) {
	diveContext.StartSpinner(" Starting ETH Node")

	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		return nil, err
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveEthHardhatNodeScript, "start_eth_node", `{"args":{}}`, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		return nil, err
	}

	responseData, skippedInstructions, err := diveContext.GetSerializedData(data)
	if err != nil {
		diveContext.Error(err.Error())
	}
	diveContext.CheckInstructionSkipped(skippedInstructions, common.DiveEthNodeAlreadyRunning)
	ethResponseData := &common.DiveserviceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		return nil, err
	}

	return result, nil

}
