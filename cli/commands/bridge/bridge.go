/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"strconv"
	"strings"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

const bridgeMainFunction = "run_btp_setup"
const runbridgeicon2icon = "start_btp_for_already_running_icon_nodes"
const runbridgeicon2ethhardhat = "start_btp_icon_to_eth_for_already_running_nodes"

var (
	chainA string
	chainB string
)

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
			diveContext.StartSpinner(fmt.Sprintf(" Starting BTP Bridge for %s,%s", chainA, chainB))

			bridge, _ := cmd.Flags().GetBool("bridge")

			if strings.ToLower(chainA) == "icon" && strings.ToLower(chainB) == "icon" {

				serviceConfig, err := common.ReadServiceJsonFile()

				if err != nil {
					diveContext.FatalError("Failed To read ServiceFile", err.Error())
				}

				if len(serviceConfig) != 0 {
					srcChain := "icon"
					dstChain := "icon-1"

					srcChainServiceResponse, err := serviceConfig["icon-0"].EncodeToString()
					if err != nil {
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}
					dstChainServiceResponse, err := serviceConfig["icon-1"].EncodeToString()

					if err != nil {
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}

					srcChainServiceName := serviceConfig["icon-0"].ServiceName
					dstChainServiceName := serviceConfig["icon-1"].ServiceName

					runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2icon, srcChain, dstChain, srcChainServiceName, dstChainServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)

				} else {

					params := getParams(chainA, chainA, strconv.FormatBool(bridge))

					runBtpSetupByRunningNodes(diveContext, enclaveCtx, params)
				}

			} else if (strings.ToLower(chainA) == "icon") && (strings.ToLower(chainB) == "eth" || strings.ToLower(chainB) == "hardhat") {

				serviceConfig, err := common.ReadServiceJsonFile()

				if err != nil {
					diveContext.FatalError("Failed To read ServiceFile", err.Error())
				}

				if len(serviceConfig) != 0 {
					srcChain := strings.ToLower(chainA)
					dstChain := strings.ToLower(chainB)

					srcChainServiceResponse, err := serviceConfig["icon-0"].EncodeToString()
					if err != nil {
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}
					dstChainServiceResponse, err := serviceConfig[dstChain].EncodeToString()

					if err != nil {
						diveContext.FatalError("Failed To read ServiceFile", err.Error())
					}

					srcChainServiceName := serviceConfig["icon-0"].ServiceName
					dstChainServiceName := serviceConfig[dstChain].ServiceName

					runBtpSetupForAlreadyRunningNodes(diveContext, enclaveCtx, runbridgeicon2ethhardhat, srcChain, dstChain, srcChainServiceName, dstChainServiceName, bridge, srcChainServiceResponse, dstChainServiceResponse)

				} else {

					params := getParams(chainA, chainB, strconv.FormatBool(bridge))

					runBtpSetupByRunningNodes(diveContext, enclaveCtx, params)
				}
			} else {
				diveContext.FatalError("Chains Not Supported", "Supported Chains [icon,eth,hardhat]")
			}

			diveContext.StopSpinner(fmt.Sprintf("BTP Bridge Setup Completed between %s and %s. Please find service details in current working directory(dive.json)", chainA, chainB))
		},
	}

	btpbridgeCmd.Flags().StringVar(&chainA, "chainA", "", "Mention Name of Supported Chain")
	btpbridgeCmd.Flags().StringVar(&chainB, "chainB", "", "Mention Name of Supported Chain")
	btpbridgeCmd.Flags().Bool("bridge", false, "Mention Bridge ENV")

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

	configData := fmt.Sprintf(`{"links": {"src":"%s","dst":"%s"},"chains" : { "%s" : %s,"%s" : %s},"contracts" : {"%s"  : {},"%s" : {}},"bridge" : "%s"}`, srcChain, dstChain, srcChain, srcChainServiceResponse, dstChain, dstChainServiceResponse, srcChain, dstChain, strconv.FormatBool(bridge))

	params := fmt.Sprintf(`{"src_chain":"%s", "dst_chain":"%s", "config_data":%s, "src_service_name":"%s", "dst_src_name":"%s"}`, srcChain, dstChain, configData, srcChainServiceName, dstChainServiceName)

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

func getParams(chainSrc, chainDst, bridge string) string {
	return fmt.Sprintf(`{"args":{"links": {"src": "%s", "dst": "%s"},"bridge":"%s"}}`, chainSrc, chainDst, bridge)
}
