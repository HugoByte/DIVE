package eth

import (
	"strings"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var EthCmd = common.NewDiveCommandBuilder().
	SetUse("eth").
	SetShort("Build, initialize and start a eth node.").
	SetLong(`The command starts an Ethereum node, initiating the process of setting up and launching a local Ethereum network. 
It establishes a connection to the Ethereum network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	SetRun(eth).
	Build()

func eth(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	cliContext.Spinner().StartWithMessage("Starting ETH Node", "green")

	responseData, err := RunEth(cliContext)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			cliContext.Spinner().StopWithMessage("ETH Node Already Running")
			cliContext.Logger().Error(common.CodeOf(err), err.Error())
			cliContext.Context().Exit(0)
		} else {
			cliContext.Logger().SetErrorToStderr()
			cliContext.Logger().Fatalf(common.CodeOf(err), err.Error())
		}
	}

	err = common.WriteServiceResponseData(responseData.ServiceName, *responseData, cliContext)

	if err != nil {
		cliContext.Spinner().Stop()
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	cliContext.Spinner().StopWithMessage("ETH Node Started. Please find service details in current working directory(services.json)")
}
