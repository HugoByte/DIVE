package types

import (
	"fmt"
	"strings"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

var (
	config string
)

func NewArchwayCmd(diveContext *common.DiveContext) *cobra.Command {

	archwayCmd := &cobra.Command{
		Use:   "archway",
		Short: "Build, initialize and start a archway node",
		Long:  "The command starts the archway network and allows node in executing contracts",
		Run: func(cmd *cobra.Command, args []string) {

			RunArchwayNode(diveContext)
		},
	}
	archwayCmd.Flags().StringVarP(&config, "config", "c", "", "provide config to start archway node ")
	archwayCmd.MarkFlagRequired("config")

	return archwayCmd
}

func RunArchwayNode(diveContext *common.DiveContext) *common.DiveserviceResponse {
	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}
	diveContext.StartSpinner(" Starting Archway Node")

	params := `{"cid":"chain-1", "key":"chain-key-1", "private_grpc":9090, "private_http":9091, "private_tcp":26656, "private_rpc":26657, "public_grpc":9090, "public_http":9091, "public_tcp":26656, "public_rpc":4564, "password":"password"}`

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", common.DiveCosmosNodeScript, "get_service_config", params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	responseData, _, _, err := diveContext.GetSerializedData(data)
	if err != nil {

		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	fmt.Println(strings.TrimSpace(responseData))

	params = fmt.Sprintf(`{"args": %s}`, strings.TrimSpace(responseData))

	data, _, err = kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", common.DiveCosmosNodeScript, "start_cosmos_node", params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	responseData, _, _, err = diveContext.GetSerializedData(data)
	if err != nil {

		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	fmt.Println(responseData)

	return nil
}
