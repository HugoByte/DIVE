package eth

import (
	"fmt"
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
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.Spinner().StartWithMessage("Starting ETH Node", "green")

	responseData, err := RunEth(cliContext)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			cliContext.Error(err)
			cliContext.Context().Exit(0)
		} else {
			cliContext.Fatal(err)
		}
	}
	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName)
	err = common.WriteServiceResponseData(responseData.ServiceName, *responseData, cliContext, serviceFileName)

	if err != nil {
		cliContext.Fatal(err)
	}

	stopMessage := fmt.Sprintf("ETH Node Started. Please find service details in current working directory(%s)", serviceFileName)
	cliContext.Spinner().StopWithMessage(stopMessage)

}
