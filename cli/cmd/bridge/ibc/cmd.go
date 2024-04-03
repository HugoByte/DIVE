package ibc

import (
	"fmt"

	"github.com/hugobyte/dive/cli/common"
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

	cliContext := common.GetCliWithKurtosisContext(common.EnclaveName)

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Fatalf("Error %s. %s", err, cmd.UsageString())
	}
	cliContext.StartSpinnerIfNotVerbose("Starting IBC Setup", common.DiveLogs)
	result, err := RunIbcRelay(cliContext)
	if err != nil {
		cliContext.Fatal(err)
	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.DiveOutFile, common.EnclaveName, shortUuid)

	serviceName := fmt.Sprintf("ibc-bridge-%s-%s", chainA, chainB)

	err = common.WriteBridgeResponseData(serviceName, result, cliContext, serviceFileName)
	if err != nil {
		cliContext.Fatal(err)
	}

	cliContext.StopSpinnerIfNotVerbose(fmt.Sprintf("IBC Setup Completed between %s and %s. Please find the service details in the output folder present in current working directory - (output/%s/%s)\n", chainA, chainB, common.EnclaveName, serviceFileName), common.DiveLogs)
}
