package icon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hugobyte/dive/cli/cmd/chains/utils"
	"github.com/hugobyte/dive/cli/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
)

func RunIconNode(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Icon Run Failed")
	}

	if configFilePath != "" && genesis == "" {
		return nil, common.WrapMessageToError(common.ErrMissingFlags, "Missing genesis flag")
	}

	var serviceConfig = &utils.IconServiceConfig{}
	err = common.LoadConfig(cli, serviceConfig, configFilePath)
	if err != nil {
		return nil, err
	}

	genesisHandler, err := genesismanager(enclaveContext)
	if err != nil {
		return nil, common.WrapMessageToError(common.ErrInvalidFile, err.Error())
	}
	params := fmt.Sprintf(`{"private_port":%d, "public_port":%d, "p2p_listen_address": %s, "p2p_address":%s, "cid": "%s","uploaded_genesis":%s,"genesis_file_path":"%s","genesis_file_name":"%s"}`, serviceConfig.Port, serviceConfig.PublicPort, serviceConfig.P2PListenAddress, serviceConfig.P2PAddress, serviceConfig.Cid, genesisHandler.uploadedFiles, genesisHandler.genesisPath, genesisHandler.genesisFile)
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveIconNodeScript, "start_icon_node")

	iconData, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Icon Run Failed")
	}

	response, services, skippedInstructions, err := common.GetSerializedData(cli, iconData)

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Icon Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Icon Run Failed ")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	iconResponseData := &common.DiveServiceResponse{}

	result, err := iconResponseData.Decode([]byte(response))

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Icon Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Icon Run Failed ")

	}

	return result, nil
}

func RunDecentralization(cli *common.Cli, params string) error {

	kurtosisEnclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return common.WrapMessageToError(err, "Icon Decentralization Failed")
	}
	starlarkConfig := common.GetStarlarkRunConfig(params, common.DiveIconDecentralizeScript, "configure_node")
	data, _, err := kurtosisEnclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, starlarkConfig)

	if err != nil {
		return common.WrapMessageToError(common.ErrStarlarkRunFailed, err.Error())
	}

	_, services, skippedInstructions, err := common.GetSerializedData(cli, data)
	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return common.WrapMessageToError(errRemove, "Icon Decentralization Failed ")
		}

		return common.WrapMessageToError(err, "Icon Decentralization Failed ")
	}
	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	return nil

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
			return nil, common.WrapMessageToError(common.ErrInvalidFile, err.Error())
		}

		_, d, err := enclaveContext.UploadFiles(genesisFilePath, genesisFileName)
		if err != nil {
			return nil, common.WrapMessageToError(common.ErrInvalidEnclaveContext, err.Error())
		}

		gm.uploadedFiles = fmt.Sprintf(`{"file_path":"%s","file_name":"%s"}`, d, genesisFileName)
	} else {
		gm.genesisFile = filepath.Base(DefaultIconGenesisFile)
		gm.genesisPath = DefaultIconGenesisFile
		gm.uploadedFiles = `{}`

	}

	return &gm, nil
}

func GetDecentralizeParams(serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID string) string {
	return fmt.Sprintf(`{"service_name":"%s","uri":"%s","keystorepath":"%s","keypassword":"%s","nid":"%s"}`, serviceName, nodeEndpoint, keystorePath, keystorepassword, networkID)
}
