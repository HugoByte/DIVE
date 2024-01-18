package btp

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/cmd/bridge/utils"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const bridgeMainFunction = "run_btp_setup"
const runBridgeIcon2icon = "start_btp_for_already_running_icon_nodes"
const runBridgeIcon2EthHardhat = "start_btp_icon_to_eth_for_already_running_nodes"

var (
	chainA   string
	chainB   string
	serviceA string
	serviceB string
)

var BtpRelayCmd = common.NewDiveCommandBuilder().
	SetUse("btp").
	SetShort("Starts BTP Relay to bridge between ChainA and ChainB").
	SetLong("Setup and Starts BTP Relay between ChainA and ChainB a").
	SetRun(btpRelay).
	AddStringFlag(&chainA, "chainA", "", "Mention Name of Supported Chain").
	AddStringFlag(&chainB, "chainB", "", "Mention Name of Supported Chain").
	AddBoolFlagP("bmvbridge", "b", false, "To Specify Which Type of BMV to be used in btp setup(if true BMV bridge is used else BMV-BTP Block is used)").
	AddStringFlag(&serviceA, "chainAServiceName", "", "Service Name of Chain A from services.json").
	AddStringFlag(&serviceB, "chainBServiceName", "", "Service Name of Chain B from services.json").
	MarkFlagRequired("chainA").
	MarkFlagRequired("chainB").
	Build()

func btpRelay(cmd *cobra.Command, args []string) {

	cliContext := common.GetCli(common.EnclaveName)

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Fatal(err)
	}
	cliContext.StartSpinnerIfNotVerbose("Starting BTP Setup", common.DiveLogs)

	bridge, err := cmd.Flags().GetBool("bmvbridge")
	if err != nil {
		cliContext.Fatal(common.WrapMessageToError(common.ErrInvalidCommandArguments, err.Error()))
	}

	chains := utils.InitChains(chainA, chainB, serviceA, serviceB, bridge)

	result, err := RunBtpSetup(cliContext, chains, bridge)

	if err != nil {
		cliContext.Fatal(err)
	}

	shortUuid, err := cliContext.Context().GetShortUuid(common.EnclaveName)
	if err != nil {
		cliContext.Fatal(err)
	}

	serviceFileName := fmt.Sprintf(common.DiveOutFile, common.EnclaveName, shortUuid)
	serviceName := fmt.Sprintf("btp-bridge-%s-%s", chainA, chainB)

	err = common.WriteBridgeResponseData(serviceName, result, cliContext, serviceFileName)
	if err != nil {
		cliContext.Fatal(err)
	}

	cliContext.StopSpinnerIfNotVerbose(fmt.Sprintf("BTP Setup Completed between %s and %s. Please find the service details in the output folder present in current working directory - (output/%s/%s)\n", chainA, chainB, common.EnclaveName, serviceFileName), common.DiveLogs)
}
