package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const genesisIcon = "github.com/hugobyte/dive/services/jvm/icon/static-files/config/genesis-icon-0.zip"

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

func (sc *IconServiceConfig) GetDefaultConfigIconNode1() {

	sc.Id = "1"
	sc.Port = 9081
	sc.PublicPort = 8091
	sc.P2PListenAddress = "7081"
	sc.P2PAddress = "8081"
	sc.Cid = "0x42f1f3"

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
		Long: `The command starts an Icon node, initiating the process of setting up and launching a local Icon network. It establishes a connection to the Icon
network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			decentralisation, _ := cmd.Flags().GetBool("decentralisation")

			serviceConfig := &IconServiceConfig{}

			if configFilePath == "" {
				serviceConfig.GetDefaultConfigIconNode0()
			} else {
				data, err := common.ReadConfigFile(configFilePath)
				if err != nil {
					serviceConfig.GetDefaultConfigIconNode0()
				}

				err = json.Unmarshal(data, serviceConfig)

				if err != nil {
					logrus.Fatalln(err)
				}

			}

			if decentralisation {

				data := RunIconNode(diveContext, serviceConfig, genesis)

				params := GetDecentralizeParms(data.ServiceName, data.PrivateEndpoint, data.KeystorePath, data.KeyPassword, data.NetworkId)

				Decentralisation(diveContext, params)

				data.WriteDiveResponse(diveContext)

			} else {

				data := RunIconNode(diveContext, serviceConfig, genesis)

				data.WriteDiveResponse(diveContext)

			}

		},
	}

	iconCmd.Flags().StringVarP(&id, "id", "i", "", "chain id")
	iconCmd.Flags().StringVarP(&genesis, "genesis", "g", "", "gen file")
	iconCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "gen file")
	iconCmd.Flags().BoolP("decentralisation", "d", false, "Decentralise Icon Node")

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
			fmt.Println("Decentralisation")

			params := GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

			Decentralisation(diveContext, params)

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

func RunIconNode(diveContext *common.DiveContext, serviceConfig *IconServiceConfig, genesisFilePath string) *common.DiveserviceResponse {

	paramData, err := serviceConfig.EncodeToString()
	if err != nil {
		logrus.Fatalln(err)
	}

	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		logrus.Errorf("Failed to fetch Enclave details :%s", err)
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "get_service_config", paramData, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	responseData := common.GetSerializedData(data)
	var genesisFile = ""
	var uploadedFiles = ""
	var genesisPath = ""

	if genesisFilePath != "" {
		genesisFileName := filepath.Base(genesisFilePath)
		r, d, err := kurtosisEnclaveContext.UploadFiles(genesisFilePath, genesisFileName)
		logrus.Infof("File Uploaded sucessfully : UUID %s", r)
		uploadedFiles = fmt.Sprintf(`{"file_path":"%s","file_name":"%s"}`, d, genesisFileName)

		if err != nil {
			panic(err)
		}
	} else {
		genesisFile = filepath.Base(genesisIcon)
		genesisPath = genesisIcon
		uploadedFiles = `{}`

	}

	params := fmt.Sprintf(`{"service_config":%s,"id":"%s","uploaded_genesis":%s,"genesis_file_path":"%s","genesis_file_name":"%s"}`, responseData, serviceConfig.Id, uploadedFiles, genesisPath, genesisFile)
	icon_data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "start_icon_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	response := common.GetSerializedData(icon_data)

	iconResponseData := &common.DiveserviceResponse{}

	result, err := iconResponseData.Decode([]byte(response))

	if err != nil {
		fmt.Println(err)
	}

	return result
}

func Decentralisation(diveContext *common.DiveContext, params string) {

	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		logrus.Errorf("Failed to fetch Enclave details :%s", err)
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", "services/jvm/icon/src/node-setup/setup_icon_node.star", "configure_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	response := common.GetSerializedData(data)

	fmt.Println(response)

}

func GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {

	return fmt.Sprintf(`{"args":{"service_name":"%s","endpoint":"%s","keystore_path":"%s","keypassword":"%s","nid":"%s"}}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

}
