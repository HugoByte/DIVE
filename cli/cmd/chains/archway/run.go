package archway

import (
	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
)

func RunArchway(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Archway Run Failed While Getting Enclave Context")
	}

	var serviceConfig = &utils.CosmosServiceConfig{}
	chainName := "archway"
	serviceConfig.ChainName = &chainName

	err = common.LoadConfig(cli, serviceConfig, configFilePath)
	if err != nil {
		return nil, err
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	runConfig := common.GetStarlarkRunConfig(encodedServiceConfigDataString, common.DiveCosmosDefaultNodeScript, runArchwayNodeWithDefaultConfigFunctionName)

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Archway Run Failed")
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Archway Run Failed.")
		}

		return nil, common.WrapMessageToErrorf(err, "%s. %s", err, "Archway Run Failed")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	archwayResponseData := &common.DiveServiceResponse{}
	result, err := archwayResponseData.Decode([]byte(responseData))

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Archway Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Archway Run Failed. Services Removed")

	}

	return result, nil
}
