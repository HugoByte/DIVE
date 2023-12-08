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

	cliContext := common.GetCli()

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

	err = cliContext.FileHandler().WriteFile("dive.json", []byte(result))
	if err != nil {
		cliContext.Fatal(err)
	}
	cliContext.StopSpinnerIfNotVerbose(fmt.Sprintf("BTP Setup Completed between %s and %s. Please find service details in current working directory(dive.json)", chainA, chainB), common.DiveLogs)
}
