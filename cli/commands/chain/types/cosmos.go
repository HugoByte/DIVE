package types

import (
	"github.com/hugobyte/dive/common"
	"github.com/spf13/cobra"
)

func NewCosmosCmd(diveContext *common.DiveContext) *cobra.Command {

	cosmosCmd := &cobra.Command{
		Use:   "cosmos",
		Short: "Build, initialize and start a cosmos node.",
		Long: `The command starts an Cosmos node, initiating the process of setting up and launching a local cosmos network. 
It establishes a connection to the Cosmos network and allows the node in executing smart contracts and maintaining the decentralized ledger.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cosmosCmd
}
