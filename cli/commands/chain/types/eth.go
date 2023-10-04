package types

import (
	"os"
	"strings"

	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

func NewEthCmd(diveContext *common.DiveContext) *cobra.Command {

	var ethCmd = &cobra.Command{
		Use:   "eth",
		Short: "Build, initialize and start a eth node.",
		Long: `The command starts an Ethereum node, initiating the process of setting up and launching a local Ethereum network. 
It establishes a connection to the Ethereum network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())

			data := RunEthNode(diveContext)

			diveContext.SetSpinnerMessage("Execution Completed")
			err := common.WriteToServiceFile(data.ServiceName, *data)
			if err != nil {
				diveContext.FatalError("Failed To Write To File", err.Error())
			}
			diveContext.StopSpinner("ETH Node Started. Please find service details in current working directory(services.json)")

		},
	}

	return ethCmd

}

func RunEthNode(diveContext *common.DiveContext) *common.DiveserviceResponse {

	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}
	diveContext.StartSpinner(" Starting ETH Node")
	starlarkConfig := diveContext.GetStarlarkRunConfig(`{"args":{}}`, common.DiveEthHardhatNodeScript, "start_eth_node")
	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {

		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	responseData, services, skippedInstructions, err := diveContext.GetSerializedData(data)

	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			diveContext.StopSpinner("Eth Node Already Running")
			os.Exit(0)
		} else {
			diveContext.StopServices(services)
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, common.DiveEthNodeAlreadyRunning)

	ethResponseData := &common.DiveserviceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Fail to Start ETH Node", err.Error())
	}

	return result

}
