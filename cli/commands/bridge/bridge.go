/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const bridgeMainFunction = "run_btp_setup"

var (
	chainA string
	chainB string
)

func NewBridgeCmd(diveContext *common.DiveContext) *cobra.Command {

	var bridgeCmd = &cobra.Command{
		Use:   "bridge",
		Short: "Command for cross chain communication between two different chains",
		Long:  `To connect two different chains using any of the supported cross chain communication protocols. This will create an relay to connect two different chains and pass any messages between them.`,
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
		Short: "Starts Bridge BTP between ChainA and Chain B",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				diveContext.FatalError("Invalid Usage of command. Find cmd", cmd.UsageString())

			}

			enclaveCtx, err := diveContext.GetEnclaveContext()

			if err != nil {
				logrus.Errorln(err)
			}

			bridge, _ := cmd.Flags().GetBool("bridge")

			params := fmt.Sprintf(`{"args":{"links": {"src": "%s", "dst": "%s"},"bridge":"%s"}}`, chainA, chainB, strconv.FormatBool(bridge))

			if strings.ToLower(chainA) == "icon" && strings.ToLower(chainB) == "icon" {

				data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, bridgeMainFunction, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

				if err != nil {
					fmt.Println(err)
				}
				response := diveContext.GetSerializedData(data)

				common.WriteToFile(response)
			} else {
				data, _, err := enclaveCtx.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveBridgeScript, bridgeMainFunction, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

				if err != nil {
					fmt.Println(err)
				}
				response := diveContext.GetSerializedData(data)

				common.WriteToFile(response)
			}
		},
	}

	btpbridgeCmd.Flags().StringVar(&chainA, "chainA", "", "Metion Name of Supported Chain")
	btpbridgeCmd.Flags().StringVar(&chainB, "chainB", "", "Metion Name of Supported Chain")
	btpbridgeCmd.Flags().Bool("bridge", false, "Mention Bridge ENV")

	btpbridgeCmd.MarkFlagRequired("chainA")
	btpbridgeCmd.MarkFlagRequired("chainB")

	return btpbridgeCmd
}
