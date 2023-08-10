/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"

	"github.com/hugobyte/dive/commands/chain/types"
	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
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

type Chains struct {
	chainA            string
	chainB            string
	chainAServiceName string
	chainBServiceName string
	bridge            string
}

func initChains(chainA, chainB, serviceA, serviceB string, bridge bool) *Chains {
	return &Chains{
		chainA:            strings.ToLower(chainA),
		chainB:            strings.ToLower(chainB),
		chainAServiceName: serviceA,
		chainBServiceName: serviceB,
		bridge:            strconv.FormatBool(bridge),
	}
}

func (c *Chains) areChainsIcon() bool {
	return (c.chainA == "icon" && c.chainB == "icon")
}

func NewBridgeCmd(diveContext *common.DiveContext) *cobra.Command {

	var bridgeCmd = &cobra.Command{
		Use:   "bridge",
		Short: "Command for cross chain communication between two different chains",
		Long: `To connect two different chains using any of the supported cross chain communication protocols.
This will create an relay to connect two different chains and pass any messages between them.`,
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Help()
		},
	}

	bridgeCmd.AddCommand(btpBridgeCmd(diveContext))

	return bridgeCmd
}

func btpBridgeCmd(diveContext *common.DiveContext) *cobra.Command {

	var btpbridgeCmd = &cobra.Command{
		Use:   "btp",
		Short: "Starts BTP Bridge between ChainA and ChainB",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(args, cmd.UsageString())

			diveContext.InitKurtosisContext()
			enclaveCtx, err := diveContext.GetEnclaveContext()

			if err != nil {
				diveContext.Error(err.Error())
			}

			bridge, _ := cmd.Flags().GetBool("bridge")

			chains := initChains(chainA, chainB, serviceA, serviceB, bridge)
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
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}
					runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainA, chains.chainB, chains.chainAServiceName, chains.chainBServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)
				} else if (chains.chainAServiceName == "" && chains.chainBServiceName != "") || (chains.chainAServiceName != "" && chains.chainBServiceName == "") {

					var chainAServiceResponse string
					var chainBServiceResponse string
					var chainAServiceName string
					var chainBServiceName string

					serviceConfig, err := common.ReadServiceJsonFile()
					if err != nil {
						diveContext.FatalError("Failed To Get Service Data", err.Error())
					}

					if chains.chainAServiceName == "" {
						response := runChain[chains.chainA](diveContext)
						chainAServiceName = response.ServiceName
						chainAServiceResponse, err = response.EncodeToString()
						if err != nil {
							diveContext.FatalError("Failed To Get Service Data", err.Error())
						}
						chainBServiceName = serviceConfig[chains.chainBServiceName].ServiceName

						chainBServiceResponse, err = serviceConfig[chains.chainBServiceName].EncodeToString()

						if err != nil {
							diveContext.FatalError("Failed To Get Service Data", err.Error())
						}

					} else {
						response := runChain[chains.chainB](diveContext)
						chainBServiceName = response.ServiceName
						chainBServiceResponse, err = response.EncodeToString()
						if err != nil {
							diveContext.FatalError("Failed To Get Service Data", err.Error())
						}
						chainAServiceName = serviceConfig[chains.chainBServiceName].ServiceName
						chainAServiceResponse, err = serviceConfig[chains.chainBServiceName].EncodeToString()

						if err != nil {
							diveContext.FatalError("Failed To Get Service Data", err.Error())
						}
					}

					runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, chains.chainA, chains.chainB, chainAServiceName, chainBServiceName, bridge, chainAServiceResponse, chainBServiceResponse)

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
	btpbridgeCmd.Flags().Bool("bridge", false, "Mention Bridge ENV")

	btpbridgeCmd.Flags().StringVar(&serviceA, "chainAServiceName", "", "Service Name of Chain A from services.json")
	btpbridgeCmd.Flags().StringVar(&serviceB, "chainBServiceName", "", "Service Name of Chain B from services.json")
	btpbridgeCmd.MarkFlagRequired("chainA")
	btpbridgeCmd.MarkFlagRequired("chainB")

	return btpbridgeCmd
}

func runBtpSetupByRunningNodes(diveContext *common.DiveContext, enclaveCtx *enclaves.EnclaveContext, params string) {
	diveContext.SetSpinnerMessage(" Executing BTP Starlark Package")

	data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, bridgeMainFunction, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

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

	configData := fmt.Sprintf(`{"links": {"src":"%s","dst":"%s"},"chains" : { "%s" : %s,"%s" : %s},"contracts" : {"%s"  : {},"%s" : {}},"bridge" : "%s"}`, srcChain, dstChain, srcChainServiceName, srcChainServiceResponse, dstChainServiceName, dstChainServiceResponse, srcChainServiceName, dstChainServiceName, strconv.FormatBool(bridge))

	params := fmt.Sprintf(`{"src_chain":"%s", "dst_chain":"%s", "config_data":%s, "src_service_name":"%s", "dst_service_name":"%s"}`, srcChain, dstChain, configData, srcChainServiceName, dstChainServiceName)

	data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, mainFunctionName, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

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

func (chains *Chains) getParams() string {
	return fmt.Sprintf(`{"args":{"links": {"src": "%s", "dst": "%s"},"bridge":"%s"}}`, chains.chainA, chains.chainB, chains.bridge)
}

func (chains *Chains) getServicesResponse() (string, string, error) {

	serviceConfig, err := common.ReadServiceJsonFile()

	if err != nil {
		return "", "", err
	}

	srcChainServiceResponse, err := serviceConfig[chains.chainAServiceName].EncodeToString()
	if err != nil {
		return "", "", err
	}
	dstChainServiceResponse, err := serviceConfig[chains.chainBServiceName].EncodeToString()

	if err != nil {
		return "", "", err
	}

	return srcChainServiceResponse, dstChainServiceResponse, nil
}
