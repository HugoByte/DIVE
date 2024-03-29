package neutron

import (
	"github.com/hugobyte/dive/cli/cmd/chains/utils"
	"github.com/hugobyte/dive/cli/common"
)

func RunNeutron(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Neutron Run Failed")
	}

	var serviceConfig = &utils.CosmosServiceConfig{}
	chainName := "neutron"
	serviceConfig.ChainName = &chainName

	err = common.LoadConfig(cli, serviceConfig, configFilePath)

	if err != nil {
		return nil, err
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	runConfig := common.GetStarlarkRunConfig(encodedServiceConfigDataString, common.DiveCosmosDefaultNodeScript, runNeutronNodeWithDefaultConfigFunctionName)

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Neutron Run Failed")
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Neutron Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Neutron Run Failed ")
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	neutronResponseData := &common.DiveServiceResponse{}
	result, err := neutronResponseData.Decode([]byte(responseData))

	if err != nil {

		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToErrorf(errRemove, "%s.%s", errRemove, "Neutron Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Neutron Run Failed ")
	}

	return result, nil
}
