package types

import (
	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/spf13/cobra"
)

func NewEthCmd(diveContext *common.DiveContext) *cobra.Command {

	var ethCmd = &cobra.Command{
		Use:   "eth",
		Short: "Build, initialize and start a eth node.",
		Long: `The command starts an Ethereum node, initiating the process of setting up and launching a local Ethereum network. It establishes a connection to the Ethereum
network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

			data, err := RunEthNode(diveContext)

			if err != nil {
				diveContext.FatalError("Fail to Start ETH Node", err.Error())
			}
			data.WriteDiveResponse(diveContext)
		},
	}

	return ethCmd

}

func RunEthNode(diveContext *common.DiveContext) (*common.DiveserviceResponse, error) {

	kurtosisEnclaveContext, err := diveContext.GetEnclaveContext()

	if err != nil {
		return nil, err
	}

	data, _, err := kurtosisEnclaveContext.RunStarlarkPackage(diveContext.Ctx, "../", "services/evm/eth/src/node-setup/start-eth-node.star", "start_eth_node", `{"args":{}}`, false, 4, []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{})

	if err != nil {
		return nil, err
	}

	responseData := common.GetSerializedData(data)

	ethResponseData := &common.DiveserviceResponse{}

	result, err := ethResponseData.Decode([]byte(responseData))

	if err != nil {
		return nil, err
	}

	return result, nil

}
