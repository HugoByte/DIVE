/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package chain

import (
	"fmt"
	"os"
	"slices"

	"github.com/hugobyte/dive/cli/commands/chain/types"
	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

// chainCmd represents the chain command
func NewChainCmd(diveContext *common.DiveContext) *cobra.Command {
	var chainCmd = &cobra.Command{

		Use:   "chain",
		Short: "Build, initialize and start a given blockchain node",
		Long: `The command builds, initializes, and starts a specified blockchain node, providing a seamless setup process.
It encompasses compiling and configuring the necessary dependencies and components required for the blockchain network. 
By executing this command, the node is launched, enabling network participation, transaction processing, and ledger 
maintenance within the specified blockchain ecosystem.`,

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

	chainCmd.AddCommand(types.NewIconCmd(diveContext))
	chainCmd.AddCommand(types.NewEthCmd(diveContext))
	chainCmd.AddCommand(types.NewHardhatCmd(diveContext))
	chainCmd.AddCommand(types.NewArchwayCmd(diveContext))
	chainCmd.AddCommand(types.NewNeutronCmd(diveContext))

	return chainCmd

}

func addSubcommandsToValidArgs(cmd *cobra.Command) {
	validArgs := cmd.ValidArgs
	for _, c := range cmd.Commands() {
		validArgs = append(validArgs, c.Name())
	}
	cmd.ValidArgs = validArgs
}