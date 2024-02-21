package hardhat

import (
	"fmt"
	"strings"

	"github.com/hugobyte/dive/cli/common"
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
	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}

	cliContext.StartSpinnerIfNotVerbose("Starting Hardhat Node", common.DiveLogs)

	responseData, err := RunHardhat(cliContext)
	if err != nil {
		if strings.Contains(err.Error(), "already running") {
			cliContext.Error(err)
			cliContext.Context().Exit(0)
		} else {
			cliContext.Fatal(err)
		}
	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.ServiceFilePath, common.EnclaveName, shortUuid)
	err = common.WriteServiceResponseData(responseData.ServiceName, *responseData, cliContext, serviceFileName)

	if err != nil {
		cliContext.Fatal(err)
	}

	stopMessage := fmt.Sprintf("Hardhat Node Started. Please find the service details in the output folder present in current working directory - (output/%s/%s)\n", common.EnclaveName, serviceFileName)
	cliContext.StopSpinnerIfNotVerbose(stopMessage, common.DiveLogs)
}
