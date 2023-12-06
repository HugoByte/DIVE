package hardhat

import (
	"strings"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var HardhatCmd = common.NewDiveCommandBuilder().
	SetUse("hardhat").
	SetShort("Build, initialize and start a hardhat node.").
	SetLong(`The command starts an hardhat node, initiating the process of setting up and launching a local hardhat network. 
It establishes a connection to the hardhat network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	SetRun(hardhat).
	Build()

func hardhat(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.Spinner().StartWithMessage("Starting Hardhat Node", "green")

	responseData, err := RunHardhat(cliContext)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			cliContext.Error(err)
			cliContext.Context().Exit(0)
		} else {
			cliContext.Fatal(err)
		}
	}

	err = common.WriteServiceResponseData(responseData.ServiceName, *responseData, cliContext)

	if err != nil {
		cliContext.Fatal(err)
	}

	cliContext.Spinner().StopWithMessage("Hardhat Node Started. Please find service details in current working directory(services.json)")

}
