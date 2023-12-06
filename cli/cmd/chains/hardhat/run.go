package hardhat

import "github.com/hugobyte/dive-core/cli/common"

func RunHardhat(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)

	if err != nil {
		return nil, common.WrapMessageToError(err, "Hardhat Run Failed")
	}

	runConfig := common.GetStarlarkRunConfig(`{}`, common.DiveEthHardhatNodeScript, "start_hardhat_node")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Hardhat Run Failed")
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)
	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if errRemove != nil {
			return nil, common.WrapMessageToError(errRemove, "Hardhat Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Hardhat Run Failed ")

	}
	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	hardhatResponseData := &common.DiveServiceResponse{}

	result, err := hardhatResponseData.Decode([]byte(responseData))

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.WrapMessageToError(errRemove, "Hardhat Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Hardhat Run Failed ")
	}

	return result, nil
}
