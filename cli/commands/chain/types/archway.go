package types

import (
	"encoding/json"
	"fmt"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

var (
	config string
)

type ArchwayServiceConfig struct {
	Cid             string `json:"cid"`
	Key             string `json:"key"`
	PrivateGrpcPort int    `json:"private_grpc"`
	PrivateHttpPort int    `json:"private_http"`
	PrivateTcpPort  int    `json:"private_tcp"`
	PrivateRpcPort  int    `json:"private_rpc"`
	PublicGrpcPort  int    `json:"public_grpc"`
	PublicHttpPort  int    `json:"public_http"`
	PublicTcpPort   int    `json:"public_tcp"`
	PublicRpcPort   int    `json:"public_rpc"`
	Password        string `json:"password"`
}

func (as *ArchwayServiceConfig) EncodeToString() (string, error) {

	data, err := json.Marshal(as)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

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

	return archwayCmd
}

func RunArchwayNode(diveContext *common.DiveContext) *common.DiveserviceResponse {
	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}
	diveContext.StartSpinner(" Starting Archway Node")
	var serviceConfig = &ArchwayServiceConfig{}

	if config != "" {
		data, err := common.ReadConfigFile(config)
		if err != nil {
			diveContext.FatalError("Failed to read service config", err.Error())
		}

		err = json.Unmarshal(data, serviceConfig)

		if err != nil {
			diveContext.FatalError("Failed unmarshall service config", err.Error())
		}

		encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

		if err != nil {
			diveContext.FatalError("Failed encode service config", err.Error())
		}

		response, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveCosmosNodeScript, "get_service_config", encodedServiceConfigDataString, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

		if err != nil {

			diveContext.FatalError("Starlark Run Failed", err.Error())

		}

		responseData, _, _, err := diveContext.GetSerializedData(response)
		if err != nil {

			diveContext.FatalError("Starlark Run Failed", err.Error())

		}
		params := fmt.Sprintf(`{"args":%s}`, responseData)

		response, _, err = kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveCosmosNodeScript, "start_cosmos_node", params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

		if err != nil {

			diveContext.FatalError("Starlark Run Failed", err.Error())

		}

		responseData, _, _, err = diveContext.GetSerializedData(response)
		if err != nil {

			diveContext.FatalError("Starlark Run Failed", err.Error())

		}

		fmt.Println(responseData)
	}

	return nil
}
