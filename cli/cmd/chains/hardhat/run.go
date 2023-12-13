package hardhat

import (
	"github.com/hugobyte/dive-core/cli/cmd/chains/utils"
	"github.com/hugobyte/dive-core/cli/common"
)

func RunHardhat(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Hardhat Run Failed While Getting Enclave Context")
	}

	var serviceConfig = &utils.HardhatServiceConfig{}
	err = serviceConfig.LoadDefaultConfig()
	if err != nil {
		return nil, err
	}

	encodedServiceConfigDataString, err := serviceConfig.EncodeToString()

	if err != nil {
		return nil, common.WrapMessageToError(common.ErrDataMarshall, err.Error())
	}

	runConfig := common.GetStarlarkRunConfig(encodedServiceConfigDataString, common.DiveEthHardhatNodeScript, "start_hardhat_node")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Hardhat Run Failed While Executing Starlark Package.")
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Hardhat Run Failed. Services Removed")
		}

		return nil, common.WrapMessageToError(err, "Hardhat Run Failed ")

	}
	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	hardhatResponseData := &common.DiveServiceResponse{}

	result, err := hardhatResponseData.Decode([]byte(responseData))

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if err != nil {
			return nil, common.WrapMessageToError(errRemove, "Hardhat Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Hardhat Run Failed ")
	}

	return result, nil
}
