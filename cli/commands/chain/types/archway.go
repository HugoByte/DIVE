package types

import (
	"encoding/json"
	"fmt"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

const (
	constructServiceConfigFunctionName          = "get_service_config"
	runArchwayNodeWithCustomServiceFunctionName = "start_cosmos_node"
	runArchwayNodeWithDefaultConfigFunctionName = "start_node_service"
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
func (as *ArchwayServiceConfig) ReadServiceConfig(path string) error {
	configData, err := common.ReadConfigFile(config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(configData, as)

	if err != nil {
		return err
	}
	return nil
}

func NewArchwayCmd(diveContext *common.DiveContext) *cobra.Command {

	archwayCmd := &cobra.Command{
		Use:   "archway",
		Short: "Build, initialize and start a archway node",
		Long:  "The command starts the archway network and allows node in executing contracts",
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())
			runResponse := RunArchwayNode(diveContext)

			common.WriteToServiceFile(runResponse.ServiceName, *runResponse)

			diveContext.StopSpinner("Archdiveway Node Started. Please find service details in current working directory(services.json)")
		},
	}
	archwayCmd.Flags().StringVarP(&config, "config", "c", "", "path to custom config json file to start archway node ")

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
	var archwayResponse = &common.DiveserviceResponse{}
	var starlarkExecutionData = ""

	if config != "" {

		err := serviceConfig.ReadServiceConfig(config)

		if err != nil {
			diveContext.FatalError("Failed read service config", err.Error())
		}
		encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

		if err != nil {
			diveContext.FatalError("Failed encode service config", err.Error())
		}

		starlarkExecutionData, err = runArchwayWithCustomServiceConfig(diveContext, kurtosisEnclaveContext, encodedServiceConfigDataString)
		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}

	} else {
		starlarkExecutionData, err = runArchwayWithDefaultServiceConfig(diveContext, kurtosisEnclaveContext)
		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}
	}

	err = json.Unmarshal([]byte(starlarkExecutionData), archwayResponse)
	if err != nil {
		diveContext.FatalError("Failed to Unmarshall Service Response", err.Error())
	}

	return archwayResponse
}

func runArchwayWithCustomServiceConfig(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext, data string) (string, error) {

	serviceExecutionResponse, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveArchwayNodeScript, constructServiceConfigFunctionName, data, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		return "", err

	}

	serviceExecutionResponseData, _, _, err := diveContext.GetSerializedData(serviceExecutionResponse)
	if err != nil {

		return "", err

	}
	params := fmt.Sprintf(`{"args":%s}`, serviceExecutionResponseData)

	nodeExecutionResponse, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveArchwayNodeScript, runArchwayNodeWithCustomServiceFunctionName, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		return "", err

	}

	nodeExecutionResponseData, _, _, err := diveContext.GetSerializedData(nodeExecutionResponse)
	if err != nil {

		return "", err

	}

	return nodeExecutionResponseData, nil
}

func runArchwayWithDefaultServiceConfig(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext) (string, error) {

	params := `{"args":{"data":{}}}`
	nodeServiceResponse, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveArchwayDefaultNodeScript, runArchwayNodeWithDefaultConfigFunctionName, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		return "", err

	}

	nodeServiceResponseData, services, skippedInstructions, err := diveContext.GetSerializedData(nodeServiceResponse)
	if err != nil {

		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, "Archway Node Already Running")

	return nodeServiceResponseData, nil
}
