/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"fmt"
	"os"

	"github.com/hugobyte/dive-alpha/cli/commands/bridge/relays"
	"github.com/hugobyte/dive-alpha/cli/common"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

func NewBridgeCmd(diveContext *common.DiveContext) *cobra.Command {

	var bridgeCmd = &cobra.Command{
		Use:   "bridge",
		Short: "Command for cross chain communication between two different chains",
		Long: `To connect two different chains using any of the supported cross chain communication protocols.
This will create an relay to connect two different chains and pass any messages between them.`,
		Run: func(cmd *cobra.Command, args []string) {

			validArgs := cmd.ValidArgs
			for _, c := range cmd.Commands() {
				validArgs = append(validArgs, c.Name())
			}
			cmd.ValidArgs = validArgs

			if len(args) == 0 {
				cmd.Help()

			} else if !slices.Contains(cmd.ValidArgs, args[0]) {

				diveContext.Log.SetOutput(os.Stderr)
				diveContext.Error(fmt.Sprintf("Invalid Subcommand: %v", args))

				cmd.Usage()
				os.Exit(1)
			}
		},
	}

	bridgeCmd.AddCommand(relays.BtpRelayCmd(diveContext))

	bridgeCmd.AddCommand(relays.IbcRelayCmd(diveContext))

	return bridgeCmd
}
