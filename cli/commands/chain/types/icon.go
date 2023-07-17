package types

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

func NewIconCmd(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *cobra.Command {
	var iconCmd = &cobra.Command{
		Use:   "icon",
		Short: "Runs Icon Node",
		Long:  ``,
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

				data := RunIconNode(ctx, kurtosisEnclaveContext, serviceConfig, genesis)

				params := GetDecentralizeParms(data.ServiceName, data.PrivateEndpoint, data.KeystorePath, data.KeyPassword, data.NetworkId)

				Decentralisation(ctx, kurtosisEnclaveContext, params)

			} else {

				data := RunIconNode(ctx, kurtosisEnclaveContext, serviceConfig, genesis)

				fmt.Println(data.EncodeToString())

			}

		},
	}

	iconCmd.Flags().StringVarP(&id, "id", "i", "", "chain id")
	iconCmd.Flags().StringVarP(&genesis, "genesis", "g", "", "gen file")
	iconCmd.Flags().StringVarP(&configFilePath, "config", "c", "", "gen file")
	iconCmd.Flags().BoolP("decentralisation", "d", false, "Decentralise Icon Node")

	decentralisationCmd := IconDecentralisationCmd(ctx, kurtosisEnclaveContext)

	iconCmd.AddCommand(decentralisationCmd)

	return iconCmd
}

func IconDecentralisationCmd(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *cobra.Command {

	var decentralisationCmd = &cobra.Command{
		Use:   "decentralize",
		Short: "Decentralise Icon Node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Decentralisation")

			params := GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

			Decentralisation(ctx, kurtosisEnclaveContext, params)

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

func RunIconNode(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext, serviceConfig *IconServiceConfig, genesisFilePath string) *common.DiveserviceResponse {

	paramData, err := serviceConfig.EncodeToString()
	if err != nil {
		logrus.Fatalln(err)
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "get_service_config", paramData, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	responseData := common.GetSerializedData(data)

	genesis_file_name := filepath.Base(genesisFilePath)
	r, d, err := kurtosisEnclaveContext.UploadFiles(genesisFilePath, genesis_file_name)

	if err != nil {
		panic(err)
	}

	logrus.Infof("File Uploaded sucessfully : UUID %s", r)
	uploadedFiles := fmt.Sprintf(`{"file_path":"%s","file_name":"%s"}`, d, genesis_file_name)

	params := fmt.Sprintf(`{"service_config":%s,"id":"%s","uploaded_genesis":%s,"genesis_file_path":"%s","genesis_file_name":"%s"}`, responseData, serviceConfig.Id, uploadedFiles, "", "")
	icon_data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "start_icon_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

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

func Decentralisation(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext, params string) {
	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/setup_icon_node.star", "configure_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	response := common.GetSerializedData(data)

	fmt.Println(response)

}

func GetDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {

	return fmt.Sprintf(`{"args":{"service_name":"%s","endpoint":"%s","keystore_path":"%s","keypassword":"%s","nid":"%s"}}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

}
