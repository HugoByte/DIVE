/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package chain

import (
	"github.com/hugobyte/dive/commands/chain/types"
	"github.com/hugobyte/dive/common"
	"github.com/spf13/cobra"
)

// chainCmd represents the chain command
func NewChainCmd(diveContext *common.DiveContext) *cobra.Command {
	var chainCmd = &cobra.Command{

		Use:   "chain",
		Short: "Build, initialize and start a given blockchain node.",
		Long: `The command builds, initializes, and starts a specified blockchain node, providing a seamless setup process. It encompasses compiling and configuring the
			   necessary dependencies and components required for the blockchain network. By executing this command, the node is launched, enabling network participation, transaction
			   processing, and ledger maintenance within the specified blockchain ecosystem.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()

		},
	}

	chainCmd.AddCommand(types.NewIconCmd(diveContext))
	chainCmd.AddCommand(types.NewEthCmd(diveContext))
	chainCmd.AddCommand(types.NewHardhatCmd(diveContext))

	return chainCmd

}
