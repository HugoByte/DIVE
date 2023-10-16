package relays

import (
	"fmt"
	"strconv"

	"github.com/hugobyte/dive/cli/commands/chain/types"
	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

const bridgeMainFunction = "run_btp_setup"
const runbridgeicon2icon = "start_btp_for_already_running_icon_nodes"
const runbridgeicon2ethhardhat = "start_btp_icon_to_eth_for_already_running_nodes"

var (
	chainA   string
	chainB   string
	serviceA string
	serviceB string
)

var runChain = map[string]func(diveContext *common.DiveContext) *common.DiveserviceResponse{
	"icon": func(diveContext *common.DiveContext) *common.DiveserviceResponse {
		nodeResponse := types.RunIconNode(diveContext)
		params := types.GetDecentralizeParms(nodeResponse.ServiceName, nodeResponse.PrivateEndpoint, nodeResponse.KeystorePath, nodeResponse.KeyPassword, nodeResponse.NetworkId)

		diveContext.SetSpinnerMessage("Starting Decentralisation")
		types.Decentralisation(diveContext, params)

		return nodeResponse

	},
	"eth": func(diveContext *common.DiveContext) *common.DiveserviceResponse {
		return types.RunEthNode(diveContext)

	},
	"hardhat": func(diveContext *common.DiveContext) *common.DiveserviceResponse {

		return types.RunHardhatNode(diveContext)
	},
}

func BtpRelayCmd(diveContext *common.DiveContext) *cobra.Command {

	var btpbridgeCmd = &cobra.Command{
		Use:   "btp",
		Short: "Starts BTP Bridge between ChainA and ChainB",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())

			diveContext.InitKurtosisContext()
			enclaveCtx, err := diveContext.GetEnclaveContext()

			if err != nil {
				diveContext.Error(err.Error())
			}

			bridge, _ := cmd.Flags().GetBool("bmvbridge") // To Specify Which Type of BMV to be used in btp setup(if true BMV bridge is used else BMV-BTP Block is used)

			chains := initChains(chainA, chainB, serviceA, serviceB, bridge)

			if err := chains.checkForBtpSupportedChains(); err != nil {
				diveContext.FatalError(err.Error(), fmt.Sprintf("Supported Chains for BTP: %v", supportedChainsForBtp))
			}

			diveContext.StartSpinner(fmt.Sprintf(" Starting BTP Bridge for %s,%s", chains.chainA, chains.chainB))

			if chains.areChainsIcon() {

				if chains.chainAServiceName != "" && chains.chainBServiceName != "" {

					srcChainServiceResponse, dstChainServiceResponse, err := chains.getServicesResponse()

					if err != nil {
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}

					runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2icon, chains.chainA, chains.chainB, chains.chainAServiceName, chains.chainBServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)

				} else {

					params := chains.getParams()

					runBtpSetupByRunningNodes(diveContext, enclaveCtx, params)

				}

			} else {
				if chains.chainAServiceName != "" && chains.chainBServiceName != "" {

					srcChainServiceResponse, dstChainServiceResponse, err := chains.getServicesResponse()

					if err != nil {
						diveContext.FatalError("failed to get service data", err.Error())
					}

					if chains.chainB == "icon" {
						runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainB, chains.chainA, chains.chainBServiceName, chains.chainAServiceName, bridge, dstChainServiceResponse, srcChainServiceResponse)
					} else {

						runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainA, chains.chainB, chains.chainAServiceName, chains.chainBServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)

					}
				} else if (chains.chainAServiceName == "" && chains.chainBServiceName != "") || (chains.chainAServiceName != "" && chains.chainBServiceName == "") {

					var chainAServiceResponse string
					var chainBServiceResponse string
					var chainAServiceName string
					var chainBServiceName string

					serviceConfig, err := common.ReadServiceJsonFile()
					if err != nil {
						diveContext.FatalError("failed to get service data", err.Error())
					}

					if chains.chainAServiceName == "" {

						chainBserviceresponse, OK := serviceConfig[chains.chainBServiceName]
						if !OK {
							diveContext.FatalError("failed to get service data", fmt.Sprint("service name not found:", chains.chainBServiceName))
						}
						chainBServiceName = chainBserviceresponse.ServiceName

						chainBServiceResponse, err = chainBserviceresponse.EncodeToString()

						if err != nil {
							diveContext.FatalError("failed to get service data", err.Error())
						}

						response := runChain[chains.chainA](diveContext)
						chainAServiceName = response.ServiceName
						chainAServiceResponse, err = response.EncodeToString()
						if err != nil {
							diveContext.FatalError("failed to get service data", err.Error())
						}

					} else if chains.chainBServiceName == "" {

						chainAserviceresponse, OK := serviceConfig[chains.chainAServiceName]
						if !OK {
							diveContext.FatalError("failed to get service data", fmt.Sprint("service name not found:", chains.chainAServiceName))
						}

						chainAServiceName = chainAserviceresponse.ServiceName

						chainAServiceResponse, err = chainAserviceresponse.EncodeToString()

						if err != nil {
							diveContext.FatalError("failed to get service data", err.Error())
						}

						response := runChain[chains.chainB](diveContext)
						chainBServiceName = response.ServiceName
						chainBServiceResponse, err = response.EncodeToString()
						if err != nil {
							diveContext.FatalError("failed to get service data", err.Error())

						}
					}
					if chains.chainB == "icon" {
						runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainB, chains.chainA, chainBServiceName, chainAServiceName, bridge, chainBServiceResponse, chainAServiceResponse)
					} else {

						runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainA, chains.chainB, chainAServiceName, chainBServiceName, bridge, chainAServiceResponse, chainBServiceResponse)

					}

				} else {
					params := chains.getParams()

					runBtpSetupByRunningNodes(diveContext, enclaveCtx, params)
				}

			}

			diveContext.StopSpinner(fmt.Sprintf("BTP Bridge Setup Completed between %s and %s. Please find service details in current working directory(dive.json)", chainA, chainB))
		},
	}

	btpbridgeCmd.Flags().StringVar(&chainA, "chainA", "", "Mention Name of Supported Chain")
	btpbridgeCmd.Flags().StringVar(&chainB, "chainB", "", "Mention Name of Supported Chain")
	btpbridgeCmd.Flags().BoolP("bmvbridge", "b", false, "To Specify Which Type of BMV to be used in btp setup(if true BMV bridge is used else BMV-BTP Block is used)")

	btpbridgeCmd.Flags().StringVar(&serviceA, "chainAServiceName", "", "Service Name of Chain A from services.json")
	btpbridgeCmd.Flags().StringVar(&serviceB, "chainBServiceName", "", "Service Name of Chain B from services.json")
	btpbridgeCmd.MarkFlagRequired("chainA")
	btpbridgeCmd.MarkFlagRequired("chainB")

	return btpbridgeCmd
}

func runBtpSetupByRunningNodes(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, params string) {
	diveContext.SetSpinnerMessage(" Executing BTP Starlark Package")
	starlarkConfig := diveContext.GetStarlarkRunConfig(params, common.DiveBridgeScript, bridgeMainFunction)
	data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}
	response, services, skippedInstructions, err := diveContext.GetSerializedData(data)
	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, "Bridge Already Running")
	common.WriteToFile(response)

}

func runBtpSetupForAlreadyRunningNodes(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, mainFunctionName string, srcChain string, dstChain string, srcChainServiceName string, dstChainServiceName string, bridge bool, srcChainServiceResponse string, dstChainServiceResponse string) {

	params := fmt.Sprintf(`{"src_chain":"%s","dst_chain":"%s", "src_chain_config":%s, "dst_chain_config":%s, "bridge":%s}`, chainA, chainB, srcChainServiceResponse, dstChainServiceResponse, strconv.FormatBool(bridge))
	starlarkConfig := diveContext.GetStarlarkRunConfig(params, common.DiveBridgeScript, mainFunctionName)
	data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}
	response, services, skippedInstructions, err := diveContext.GetSerializedData(data)
	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	diveContext.CheckInstructionSkipped(skippedInstructions, "Bridge Already Running")

	common.WriteToFile(response)

}
