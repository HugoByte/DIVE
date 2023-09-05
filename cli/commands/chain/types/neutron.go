package types

import (
	"encoding/json"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

const (
	runNeutronNodeWithDefaultConfigFunctionName = "start_node_service"
)


// TODO: Implement custom Neutron node configuration.
// This code can be used when writing functionality for running a Neutron node with a custom configuration.
// Currently, it's commented out as the custom configuration feature is not yet implemented.
// Uncomment the code below and adapt it to support custom configuration options.

/*

type NeurtronServiceConfig struct {
	PrivateGrpcPort int    `json:"private_grpc"`
	PrivateHttpPort int    `json:"private_http"`
	PrivateTcpPort  int    `json:"private_tcp"`
	PrivateRpcPort  int    `json:"private_rpc"`
	PublicGrpcPort  int    `json:"public_grpc"`
	PublicHttpPort  int    `json:"public_http"`
	PublicTcpPort   int    `json:"public_tcp"`
	PublicRpcPort   int    `json:"public_rpc"`
}

func (as *NeurtronServiceConfig) EncodeToString() (string, error) {

	data, err := json.Marshal(as)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func (as *NeurtronServiceConfig) ReadServiceConfig(path string) error {
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
*/

func NewNeutronCmd(diveContext *common.DiveContext) *cobra.Command {

	neutronCmd := &cobra.Command{
		Use:   "neutron",
		Short: "Build, initialize and start a neutron node",
		Long:  "The command starts the neutron network and allows node in executing contracts",
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())
			runResponse := RunNeutronNode(diveContext)

			common.WriteToServiceFile(runResponse.ServiceName, *runResponse)

			diveContext.StopSpinner("Neutron Node Started. Please find service details in current working directory(services.json)")
		},
	}
	neutronCmd.Flags().StringVarP(&config, "config", "c", "", "path to custom config json file to start neutron node ")

	return neutronCmd
}


func RunNeutronNode(diveContext *common.DiveContext) *common.DiveserviceResponse {
	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}

	diveContext.StartSpinner(" Starting Neutron Node")
	var neutronResponse = &common.DiveserviceResponse{}
	var starlarkExecutionData = ""
	starlarkExecutionData, err = runNeutronWithDefaultServiceConfig(diveContext, kurtosisEnclaveContext)
	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}
	err = json.Unmarshal([]byte(starlarkExecutionData), neutronResponse)

	if err != nil {
		diveContext.FatalError("Failed to Unmarshall Service Response", err.Error())
	}

	return neutronResponse

}


func runNeutronWithDefaultServiceConfig(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext) (string, error) {

	params := `{"args":{"data":{}}}`
	nodeServiceResponse, _, err := enclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveNeutronDefaultNodeScript, runNeutronNodeWithDefaultConfigFunctionName, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		return "", err

	}

	nodeServiceResponseData, services, skippedInstructions, err := diveContext.GetSerializedData(nodeServiceResponse)
	if err != nil {

		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, "Nueutron Node Already Running")

	return nodeServiceResponseData, nil
}