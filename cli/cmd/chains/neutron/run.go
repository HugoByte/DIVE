package neutron

import (
	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
)

func RunNeutron(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)

	if err != nil {
		return nil, err
	}

	var serviceConfig = &utils.CosmosServiceConfig{}

	err = common.LoadConfig(cli, serviceConfig, configFilePath)

	if err != nil {
		return nil, err
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	if err != nil {
		return nil, common.Errorc(common.InvalidEnclaveConfigError, err.Error())
	}

	runConfig := common.GetStarlarkRunConfig(encodedServiceConfigDataString, common.DiveNeutronDefaultNodeScript, runNeutronNodeWithDefaultConfigFunctionName)

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.Errorc(common.FileError, err.Error())
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {

		err = cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
		}

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.Errorc(common.KurtosisContextError, "Already Running")
	}

	neutronResponseData := &common.DiveServiceResponse{}
	result, err := neutronResponseData.Decode([]byte(responseData))

	if err != nil {

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	return result, nil
}
