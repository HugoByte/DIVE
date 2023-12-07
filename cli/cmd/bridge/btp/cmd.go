package btp

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const bridgeMainFunction = "run_btp_setup"
const runbridgeicon2icon = "start_btp_for_already_running_icon_nodes"
const runbridgeicon2ethhardhat = "start_btp_icon_to_eth_for_already_running_nodes"

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
	Build()

func btpRelay(cmd *cobra.Command, args []string) {}
