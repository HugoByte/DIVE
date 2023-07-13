package types

import (
	"context"
	"fmt"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Id      string
	Genesis string
)

func NewIconCmd(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *cobra.Command {
	var iconCmd = &cobra.Command{
		Use:   "icon",
		Short: "Runs Icon Chain Node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runIconNode(ctx, kurtosisEnclaveContext)
		},
	}

	iconCmd.Flags().StringVarP(&Id, "id", "i", "", "chain id")
	iconCmd.Flags().StringVarP(&Genesis, "genesis", "g", "", "gen file")

	return iconCmd
}

func runIconNode(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) {

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "get_service_config", `{"id":"0","private_port":9080,"public_port":8090,"p2p_listen_address":"7080","p2p_address":"8080","cid":"0xacbc4e"}`, false, 4)

	if err != nil {
		fmt.Println(err)
	}

	var serializedOutputObj string
	for executionResponseLine := range data {
		runFinishedEvent := executionResponseLine.GetRunFinishedEvent()
		if runFinishedEvent == nil {
			logrus.Info("Execution in progress...")
		} else {
			logrus.Info("Execution finished successfully")
			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()
			} else {
				panic("Starlark run failed")
			}
		}
	}

	params := fmt.Sprintf(`{"service_config":%s,"id":"0","start_file_name":"start-icon.sh"}`, serializedOutputObj)

	kurtosisEnclaveContext.UploadFiles("/Users/soul/Garage/HugoByte/DIVE/services/jvm/icon/static-files/config/genesis-icon-0.zip", "genisis")

	icon_data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "start_icon_node", params, false, 4)

	if err != nil {
		fmt.Println(err)
	}

	for executionResponseLine := range icon_data {
		runFinishedEvent := executionResponseLine.GetRunFinishedEvent()
		if runFinishedEvent == nil {
			logrus.Info("Execution in progress...")
		} else {
			logrus.Info("Execution finished successfully")
			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()
			} else {
				panic("Starlark run failed")
			}
		}
	}

	fmt.Println(serializedOutputObj)
}
