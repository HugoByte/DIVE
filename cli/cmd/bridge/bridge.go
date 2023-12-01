package bridge

import (
	"github.com/hugobyte/dive-core/cli/cmd/bridge/relays"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var BridgeCmd = common.NewDiveCommandBuilder().
	SetUse("bridge").
	SetShort("Command for cross chain communication between two different chains").
	SetLong(`To connect two different chains using any of the supported cross chain communication protocols.
This will create an relay to connect two different chains and pass any messages between them.`).
	SetRun(bridge).
	AddCommand(relays.BtpRelayCmd).
	AddCommand(relays.IbcRelayCmd).
	Build()

func bridge(cmd *cobra.Command, args []string) {}
