package types

import (
	"fmt"
	"strings"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

func NewCosmosCmd(diveContext *common.DiveContext) *cobra.Command {

	cosmosCmd := &cobra.Command{
		Use:   "cosmos",
		Short: "Build, initialize and start a cosmos node.",
		Long: `The command starts an Cosmos node, initiating the process of setting up and launching a local cosmos network. 
It establishes a connection to the Cosmos network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {
			RunCosmosNode(diveContext)
		},
	}

	return cosmosCmd
}

func RunCosmosNode(diveContext *common.DiveContext) *common.DiveserviceResponse {
	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}
	diveContext.StartSpinner(" Starting Cosmos Node")

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
