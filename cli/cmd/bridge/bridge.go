package bridge

import (
	"os"
	"slices"

	"github.com/hugobyte/dive-core/cli/cmd/bridge/btp"
	"github.com/hugobyte/dive-core/cli/cmd/bridge/ibc"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var BridgeCmd = common.NewDiveCommandBuilder().
	SetUse("bridge").
	SetShort("Command for cross chain communication between two different chains").
	SetLong(`To connect two different chains using any of the supported cross chain communication protocols.
This will create an relay to connect two different chains and pass any messages between them.`).
	SetRun(bridge).
	AddCommand(btp.BtpRelayCmd).
	AddCommand(ibc.IbcRelayCmd).
	Build()

func bridge(cmd *cobra.Command, args []string) {
	cli := common.GetCli(common.EnclaveName)
	validArgs := cmd.ValidArgs
	for _, c := range cmd.Commands() {
		validArgs = append(validArgs, c.Name())
	}
	cmd.ValidArgs = validArgs

	if len(args) == 0 {
		cmd.Help()

	} else if !slices.Contains(cmd.ValidArgs, args[0]) {
		cli.Error(common.WrapMessageToErrorf(common.ErrInvalidCommand, "%s", cmd.UsageString()))
		os.Exit(1)
	}
}
