package types

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"os"

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

type IconserviceResponse struct {
	ServiceName     string `json:"service_name"`
	PublicEndpoint  string `json:"endpoint_public"`
	PrivateEndpoint string `json:"endpoint"`
	KeyPassword     string `json:"keypassword"`
	KeystorePath    string `json:"keystore_path"`
	Network         string `json:"network"`
	NetworkName     string `json:"network_name"`
	NetworkId       string `json:"nid"`
}

type IconServiceConfig struct {
	Id               string `json:"id" default:"0"`
	Port             int    `json:"private_port"`
	PublicPort       int    `json:"public_port"`
	P2PListenAddress string `json:"p2p_listen_address"`
	P2PAddress       string `json:"p2p_address"`
	Cid              string `json:"cid"`
}

func (sc *IconServiceConfig) defaultServiceConfig() {

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

func (icon *IconserviceResponse) EncodeToString() (string, error) {
	encodedBytes, err := json.Marshal(icon)
	if err != nil {
		return "", nil
	}

	return string(encodedBytes), nil
}
func (icon *IconserviceResponse) Decode(responseData []byte) (*IconserviceResponse, error) {

	err := json.Unmarshal(responseData, &icon)
	if err != nil {
		return nil, err
	}
	return icon, nil
}

func NewIconCmd(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext) *cobra.Command {
	var iconCmd = &cobra.Command{
		Use:   "icon",
		Short: "Runs Icon Chain Node",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			decentralisation, _ := cmd.Flags().GetBool("decentralisation")

			serviceConfig := &IconServiceConfig{}

			if configFilePath == "" {
				serviceConfig.defaultServiceConfig()
			} else {
				data, err := readConfigFile(configFilePath)
				if err != nil {
					serviceConfig.defaultServiceConfig()
				}

				err = json.Unmarshal(data, serviceConfig)

				if err != nil {
					logrus.Fatalln(err)
				}

			}

			if decentralisation {

				data := runIconNode(ctx, kurtosisEnclaveContext, serviceConfig)

				params := getDecentralizeParms(data.ServiceName, data.PrivateEndpoint, data.KeystorePath, data.KeyPassword, data.NetworkId)

				Decentralisation(ctx, kurtosisEnclaveContext, params)

			} else {

				data := runIconNode(ctx, kurtosisEnclaveContext, serviceConfig)

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

			params := getDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

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

func runIconNode(ctx context.Context, kurtosisEnclaveContext *enclaves.EnclaveContext, serviceConfig *IconServiceConfig) *IconserviceResponse {

	paramData, err := serviceConfig.EncodeToString()
	if err != nil {
		logrus.Fatalln(err)
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "get_service_config", paramData, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	responseData := getSerializedData(data)

	genesis_file_name := filepath.Base(genesis)
	r, d, err := kurtosisEnclaveContext.UploadFiles(genesis, genesis_file_name)

	if err != nil {
		panic(err)
	}

	logrus.Infof("File Uploaded sucessfully : UUID %s", r)

	params := fmt.Sprintf(`{"service_config":%s,"id":"%s","start_file_name":"start-icon.sh","genesis_file_path":"%s","genesis_file_name":"%s"}`, responseData, serviceConfig.Id, d, genesis_file_name)
	icon_data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(ctx, "../", "services/jvm/icon/src/node-setup/start_icon_node.star", "start_icon_node", params, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		fmt.Println(err)
	}

	response := getSerializedData(icon_data)

	iconResponseData := &IconserviceResponse{}

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

	response := getSerializedData(data)

	fmt.Println(response)

}

func getSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) string {

	var serializedOutputObj string
	for executionResponseLine := range response {
		fmt.Println(executionResponseLine)
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

	return serializedOutputObj

}

func getDecentralizeParms(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {

	return fmt.Sprintf(`{"args":{"service_name":"%s","endpoint":"%s","keystore_path":"%s","keypassword":"%s","nid":"%s"}}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)

}

func readConfigFile(filePath string) ([]byte, error) {

	file, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return file, nil
}
