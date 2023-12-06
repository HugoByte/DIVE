package hardhat

import "github.com/hugobyte/dive-core/cli/common"

func RunHardhat(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)

	if err != nil {
		return nil, err
	}

	runConfig := common.GetStarlarkRunConfig(`{}`, common.DiveEthHardhatNodeScript, "start_hardhat_node")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
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

	hardhatResponseData := &common.DiveServiceResponse{}

	result, err := hardhatResponseData.Decode([]byte(responseData))

	if err != nil {
		err = cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
		}

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	return result, nil
}
