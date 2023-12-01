package relays

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var BtpRelayCmd = common.NewDiveCommandBuilder().
	SetUse("btp").
	SetShort("Starts BTP Relay to bridge between ChainA and ChainB").
	SetLong("Setup and Starts BTP Relay between ChainA and ChainB a").
	SetRun(btpRelay).
	Build()

func btpRelay(cmd *cobra.Command, args []string) {}
