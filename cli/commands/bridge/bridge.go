/*
Copyright © 2023 Hugobyte AI Labs <hello@hugobyte.com>
*/
package bridge

import (
	"github.com/hugobyte/dive/cli/commands/bridge/relyas"
	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

func NewBridgeCmd(diveContext *common.DiveContext) *cobra.Command {

	var bridgeCmd = &cobra.Command{
		Use:   "bridge",
		Short: "Command for cross chain communication between two different chains",
		Long: `To connect two different chains using any of the supported cross chain communication protocols.
This will create an relay to connect two different chains and pass any messages between them.`,
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Help()
		},
	}

	bridgeCmd.AddCommand(relyas.BtpRelayCmd(diveContext))

	bridgeCmd.AddCommand(relyas.IbcRelayCmd(diveContext))

	return bridgeCmd
}
