package ibc

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var (
	chainA   string
	chainB   string
	serviceA string
	serviceB string
)

var IbcRelayCmd = common.NewDiveCommandBuilder().
	SetUse("ibc").
	SetShort("Start connection between Cosmos based chainA and ChainB and initiate communication between them").
	SetLong(`This Command deploy , initialize the contracts and make it ready for ibc.
Along with that setup and starts the ibc relayer to establish communication between chains specified`).
	SetRun(ibcRelay).
	AddStringFlag(&chainA, "chainA", "", "Mention Name of Supported Chain").
	AddStringFlag(&chainB, "chainB", "", "Mention Name of Supported Chain").
	AddStringFlag(&serviceA, "chainAServiceName", "", "Service Name of Chain A from services.json").
	AddStringFlag(&serviceB, "chainBServiceName", "", "Service Name of Chain B from services.json").
	MarkFlagRequired("chainA").
	MarkFlagRequired("chainB").
	Build()

func ibcRelay(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}
	cliContext.StartSpinnerIfNotVerbose("Starting IBC Setup", common.DiveLogs)
	result, err := RunIbcRelay(cliContext)
	if err != nil {
		cliContext.Fatal(err)
	}

	err = cliContext.FileHandler().WriteFile("dive.json", []byte(result))
	if err != nil {
		cliContext.Fatal(err)
	}

	cliContext.StopSpinnerIfNotVerbose(fmt.Sprintf("IBC Setup Completed between %s and %s. Please find service details in current working directory(dive.json)", chainA, chainB), common.DiveLogs)
}
