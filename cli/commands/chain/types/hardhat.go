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

			if len(args) != 0 {
				diveContext.FatalError("Invalid Usage of command. Find cmd", cmd.UsageString())

			}

			data, err := RunHardhatNode(diveContext)

			if err != nil {
				diveContext.FatalError("Fail to Start Hardhat Node", err.Error())
			}
			diveContext.SetSpinnerMessage("Execution Completed")

			err = data.WriteDiveResponse(diveContext)
			if err != nil {
				diveContext.FatalError("Failed To Write To File", err.Error())
			}
			diveContext.StopSpinner("Hardhat Node Started")
		},
	}

	return ethCmd

}

func RunHardhatNode(diveContext *common.DiveContext) (*common.DiveserviceResponse, error) {

	diveContext.StartSpinner(" Starting Hardhat Node")

	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		return nil, err
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveEthHardhatNodeScript, "start_hardhat_node", "{}", common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		return nil, err
	}

	responseData := diveContext.GetSerializedData(data)

	hardhatResponseData := &common.DiveserviceResponse{}

	result, err := hardhatResponseData.Decode([]byte(responseData))

	if err != nil {
		return nil, err
	}

	return result, nil

}
