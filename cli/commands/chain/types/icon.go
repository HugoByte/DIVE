package types

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/spf13/cobra"
)

const DefaultIconGenesisFile = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/genesis-icon-0.zip"

var (
	id               = ""
	genesis          = ""
	serviceName      = ""
	keystorePath     = ""
	keystorepassword = ""
	networkID        = ""
	nodeEndpoint     = ""
	configFilePath   = ""
)

type IconServiceConfig struct {
	Id               string `json:"id" default:"0"`
	Port             int    `json:"private_port"`
	PublicPort       int    `json:"public_port"`
	P2PListenAddress string `json:"p2p_listen_address"`
	P2PAddress       string `json:"p2p_address"`
	Cid              string `json:"cid"`
}

func (sc *IconServiceConfig) GetDefaultConfigIconNode0() {

	sc.Id = "0"
	sc.Port = 9080
	sc.PublicPort = 8090
	sc.P2PListenAddress = "7080"
	sc.P2PAddress = "8080"
	sc.Cid = "0xacbc4e"

}

func (sc *IconServiceConfig) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(sc)
	if err != nil {
		return "", nil
	}

	return string(encodedBytes), nil
}

func NewIconCmd(diveContext *common.DiveContext) *cobra.Command {
	var iconCmd = &cobra.Command{
		Use:   "icon",
		Short: "Build, initialize and start a icon node.",
		Long: `The command starts an Icon node, initiating the process of setting up and launching a local Icon network.
It establishes a connection to the Icon network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			common.ValidateCmdArgs(args, cmd.UsageString())

			decentralisation, _ := cmd.Flags().GetBool("decentralisation")

			if decentralisation {

				nodeResponse := RunIconNode(diveContext)

				params := GetDecentralizeParms(nodeResponse.ServiceName, nodeResponse.PrivateEndpoint, nodeResponse.KeystorePath, nodeResponse.KeyPassword, nodeResponse.NetworkId)

				diveContext.SetSpinnerMessage("Starting Decentralisation")
				Decentralisation(diveContext, params)

				err := common.WriteToServiceFile(nodeResponse.NetworkName, *nodeResponse)

				if err != nil {
					diveContext.FatalError("Failed To Write To File", err.Error())
				}

				diveContext.StopSpinner("Icon Node Started. Please find service details in current working directory(services.json) in current working directory")

			} else {

				nodeResponse := RunIconNode(diveContext)

				err := common.WriteToServiceFile(nodeResponse.NetworkName, *nodeResponse)

				if err != nil {
					diveContext.FatalError("Failed To Write To File", err.Error())
				}

				diveContext.StopSpinner("Icon Node Started. Please find service details in current working directory(services.json)")
			}

		},
	}

	iconCmd.Flags().StringVarP(&genesis, "genesis", "g", "", "path to custom genesis file")
	iconCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "path to custom config json file")
	iconCmd.Flags().BoolP("decentralisation", "d", false, "decentralise Icon Node")

	decentralisationCmd := IconDecentralisationCmd(diveContext)

	iconCmd.AddCommand(decentralisationCmd)

	return iconCmd
}

func IconDecentralisationCmd(diveContext *common.DiveContext) *cobra.Command {

	var decentralisationCmd = &cobra.Command{
		Use:   "decentralize",
		Short: "Decentralise already running Icon Node",
		Long:  `Decentralise Icon Node is necessary if you want to connect your local icon node to BTP network`,
		Run: func(cmd *cobra.Command, args []string) {

			params := GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)
			diveContext.InitKurtosisContext()
			Decentralisation(diveContext, params)

			diveContext.StopSpinner(fmt.Sprintln("Decentralisation Completed.Please find service details in dive.json"))
		},
	}
	decentralisationCmd.Flags().StringVarP(&serviceName, "serviceName", "s", "", "service name")
	decentralisationCmd.Flags().StringVarP(&nodeEndpoint, "nodeEndpoint", "e", "", "endpoint address")
	decentralisationCmd.Flags().StringVarP(&keystorePath, "keystorePath", "k", "", "keystore path")
	decentralisationCmd.Flags().StringVarP(&keystorepassword, "keyPassword", "p", "", "keypassword")
	decentralisationCmd.Flags().StringVarP(&networkID, "nid", "n", "", "NetworkId of Icon Node")

	decentralisationCmd.MarkFlagRequired("serviceName")
	decentralisationCmd.MarkFlagRequired("nodeEndpoint")
	decentralisationCmd.MarkFlagRequired("keystorePath")
	decentralisationCmd.MarkFlagRequired("keyPassword")
	decentralisationCmd.MarkFlagRequired("nid")

	return decentralisationCmd

}

func RunIconNode(diveContext *common.DiveContext) *common.DiveserviceResponse {

	// Initialse Kurtosis Context

	diveContext.InitKurtosisContext()

	serviceConfig, err := getConfig()
	if err != nil {
		diveContext.FatalError("Failed To Get Node Service Config", err.Error())
	}

	paramData, err := serviceConfig.EncodeToString()
	if err != nil {
		diveContext.FatalError("Encoding Failed", err.Error())
	}

	diveContext.StartSpinner(" Starting Icon Node")
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrive Enclave Context", err.Error())
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveIconNodeScript, "get_service_config", paramData, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	responseData, services, skippedInstructions, err := diveContext.GetSerializedData(data)

	if err != nil {
		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}

	genesisHandler, err := genesismanager(kurtosisEnclaveContext)
	if err != nil {
		diveContext.FatalError("Failed To Get Genesis", err.Error())
	}

	diveContext.CheckInstructionSkipped(skippedInstructions, "Instruction Executed Already")

	params := fmt.Sprintf(`{"service_config":%s,"id":"%s","uploaded_genesis":%s,"genesis_file_path":"%s","genesis_file_name":"%s"}`, responseData, serviceConfig.Id, genesisHandler.uploadedFiles, genesisHandler.genesisPath, genesisHandler.genesisFile)
	icon_data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveIconNodeScript, "start_icon_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {

		diveContext.StopServices(services)

		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	diveContext.SetSpinnerMessage(" Finalizing Icon Node")

	response, services, skippedInstructions, err := diveContext.GetSerializedData(icon_data)

	if err != nil {

		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, common.DiveIconNodeAlreadyRunning)

	iconResponseData := &common.DiveserviceResponse{}

	result, err := iconResponseData.Decode([]byte(response))

	if err != nil {

		diveContext.StopServices(services)

		diveContext.FatalError("Failed To Unmarshall", err.Error())
	}

	return result
}

func Decentralisation(diveContext *common.DiveContext, params string) {

	diveContext.StartSpinner(" Starting Icon Node Decentralisation")
	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		diveContext.FatalError("Failed To Retrieve Enclave Context", err.Error())
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(diveContext.Ctx, common.DiveRemotePackagePath, common.DiveIconDecentraliseScript, "configure_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		diveContext.FatalError("Starlark Run Failed", err.Error())
	}

	_, services, skippedInstructions, err := diveContext.GetSerializedData(data)
	if err != nil {

		diveContext.StopServices(services)
		diveContext.FatalError("Starlark Run Failed", err.Error())

	}
	diveContext.CheckInstructionSkipped(skippedInstructions, "Decntralization Already Completed")

}

func GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {

	return fmt.Sprintf(`{"args":{"service_name":"%s","endpoint":"%s","keystore_path":"%s","keypassword":"%s","nid":"%s"}}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

}

func getConfig() (*IconServiceConfig, error) {
	// Init Icon Node Service Config

	serviceConfig := &IconServiceConfig{}

	if configFilePath == "" {
		serviceConfig.GetDefaultConfigIconNode0()
	} else {
		data, err := common.ReadConfigFile(configFilePath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, serviceConfig)

		if err != nil {
			return nil, err
		}

	}

	return serviceConfig, nil
}

type genesisHandler struct {
	genesisFile   string
	uploadedFiles string
	genesisPath   string
}

func genesismanager(enclaveContext *enclaves.EnclaveContext) (*genesisHandler, error) {

	gm := genesisHandler{}

	var genesisFilePath = genesis

	if genesisFilePath != "" {
		genesisFileName := filepath.Base(genesisFilePath)
		if _, err := os.Stat(genesisFilePath); err != nil {
			return nil, err
		}

		_, d, err := enclaveContext.UploadFiles(genesisFilePath, genesisFileName)
		if err != nil {
			return nil, err
		}

		gm.uploadedFiles = fmt.Sprintf(`{"file_path":"%s","file_name":"%s"}`, d, genesisFileName)
	} else {
		gm.genesisFile = filepath.Base(DefaultIconGenesisFile)
		gm.genesisPath = DefaultIconGenesisFile
		gm.uploadedFiles = `{}`

	}

	return &gm, nil
}
