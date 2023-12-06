package eth

import (
	"github.com/hugobyte/dive-core/cli/common"
)

func RunEth(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.DiveEnclave)
	if err != nil {
		return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
	}
	runConfig := common.GetStarlarkRunConfig(`{}`, common.DiveEthHardhatNodeScript, "start_eth_node")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapCodeToError(err, common.KurtosisContextError, "Starlark Run Failed")
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

	ethResponseData := &common.DiveServiceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		err = cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.Errorc(common.InvalidEnclaveContextError, err.Error())
		}

		return nil, common.Errorc(common.KurtosisContextError, err.Error())
	}

	return result, nil
}
