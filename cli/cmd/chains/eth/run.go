package eth

import (
	"github.com/hugobyte/dive-core/cli/common"
)

func RunEth(cli *common.Cli) (*common.DiveServiceResponse, error) {

	enclaveContext, err := cli.Context().GetEnclaveContext(common.EnclaveName)
	if err != nil {
		return nil, common.WrapMessageToError(err, "Eth Run Failed while getting Enclave Context")
	}
	runConfig := common.GetStarlarkRunConfig(`{}`, common.DiveEthHardhatNodeScript, "start_eth_node")

	response, _, err := enclaveContext.RunStarlarkRemotePackage(cli.Context().GetContext(), common.DiveRemotePackagePath, runConfig)

	if err != nil {
		return nil, common.WrapMessageToErrorf(common.ErrStarlarkRunFailed, "%s. %s", err, "Eth Run Failed")
	}

	responseData, services, skippedInstructions, err := common.GetSerializedData(cli, response)

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.EnclaveName)
		if err != nil {
			return nil, common.WrapMessageToError(errRemove, "Eth Run Failed ")
		}

		return nil, common.WrapMessageToError(err, "Eth Run Failed ")

	}

	if cli.Context().CheckSkippedInstructions(skippedInstructions) {
		return nil, common.WrapMessageToError(common.ErrStarlarkResponse, "Already Running")
	}

	ethResponseData := &common.DiveServiceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		errRemove := cli.Context().RemoveServicesByServiceNames(services, common.DiveEnclave)
		if err != nil {
			return nil, common.WrapMessageToError(errRemove, "Eth Run Failed ")
		}

		return nil, common.WrapMessageToErrorf(common.ErrDataUnMarshall, "%s.%s", err, "Eth Run Failed ")
	}

	return result, nil
}
