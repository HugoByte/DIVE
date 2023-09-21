package types

import (
	"encoding/json"
	"fmt"

	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

// Constants for function names
const (
	runNeutronNodeWithDefaultConfigFunctionName = "start_node_service"
	runNeutronNodeWithCustomServiceFunctionName = "start_neutron_node"
	construcNeutrontServiceConfigFunctionName   = "get_service_config"
)

// Variable to store the Neutron node configuration file path
var (
	neutron_node_config string
)

// NeutronServiceConfig stores configuration parameters for the Neutron service.
type NeutronServiceConfig struct {
	ChainID    string `json:"chainId"`
	Key        string `json:"key"`
	Password   string `json:"password"`
	PublicGrpc int    `json:"public_grpc"`
	PublicTCP  int    `json:"public_tcp"`
	PublicHTTP int    `json:"public_http"`
	PublicRPC  int    `json:"public_rpc"`
}

// EncodeToString encodes the NeutronServiceConfig struct to a JSON string.
func (as *NeutronServiceConfig) EncodeToString() (string, error) {
	data, err := json.Marshal(as)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadServiceConfig reads the Neutron service configuration from a JSON file.
func (as *NeutronServiceConfig) ReadServiceConfig(path string) error {
	configData, err := common.ReadConfigFile(neutron_node_config)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configData, as)
	if err != nil {
		return err
	}
	return nil
}

// NewNeutronCmd creates a new Cobra command for the Neutron service.
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
	neutronCmd.Flags().StringVarP(&neutron_node_config, "config", "c", "", "path to custom config json file to start neutron node ")
	return neutronCmd
}

// RunNeutronNode starts the Neutron node.
func RunNeutronNode(diveContext *common.DiveContext) *common.DiveserviceResponse {
	diveContext.InitKurtosisContext()
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()
	if err != nil {
		diveContext.FatalError("Failed To Retrieve Enclave Context", err.Error())
	}

	diveContext.StartSpinner("Starting Neutron Node")
	var serviceConfig = &NeutronServiceConfig{}
	var neutronResponse = &common.DiveserviceResponse{}
	var starlarkExecutionData = ""

	if neutron_node_config != "" {
		err := serviceConfig.ReadServiceConfig(neutron_node_config)
		if err != nil {
			diveContext.FatalError("Failed read service config", err.Error())
		}

		encodedServiceConfigDataString, err := serviceConfig.EncodeToString()
		if err != nil {
			diveContext.FatalError("Failed to encode service config", err.Error())
		}

		// Run Neutron Node with custom service config
		starlarkExecutionData, err = RunNeutronWithServiceConfig(diveContext, kurtosisEnclaveContext, encodedServiceConfigDataString)
		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}
	} else {
		// Run Neutron Node with default service config
		starlarkExecutionData, err = RunNeutronWithServiceConfig(diveContext, kurtosisEnclaveContext, "{}")
		if err != nil {
			diveContext.FatalError("Starlark Run Failed", err.Error())
		}
	}

	err = json.Unmarshal([]byte(starlarkExecutionData), neutronResponse)
	if err != nil {
		diveContext.FatalError("Failed to Unmarshal Service Response", err.Error())
	}

	return neutronResponse
}

// RunNeutronWithServiceConfig runs the Neutron service with the provided configuration data.
func RunNeutronWithServiceConfig(diveContext *common.DiveContext, enclaveContext *enclaves.EnclaveContext, data string) (string, error) {
	params := fmt.Sprintf(`{"args":{"data":%s}}`, data)
	nodeServiceResponse, _, err := enclaveContext.RunStarlarkPackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveNeutronDefaultNodeScript, runNeutronNodeWithDefaultConfigFunctionName, params, common.DiveDryRun, common.DiveDefaultParallelism, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})
	if err != nil {
		return "", err
	}

	nodeServiceResponseData, services, skippedInstructions, err := diveContext.GetSerializedData(nodeServiceResponse)
	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	diveContext.CheckInstructionSkipped(skippedInstructions, "Neutron Node Already Running")
	return nodeServiceResponseData, nil
}
