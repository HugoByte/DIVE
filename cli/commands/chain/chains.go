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
		Short: "runs specfied chain",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()

		},
	}

	enclaveCtx, err := diveContext.GetEnclaveContext()
	if err != nil {
		panic(err)

	}

	chainCmd.AddCommand(types.NewIconCmd(diveContext.Ctx, enclaveCtx))
	chainCmd.AddCommand(types.NewEthCmd(diveContext.Ctx, enclaveCtx))
	chainCmd.AddCommand(types.NewHardhatCmd(diveContext.Ctx, enclaveCtx))

	return chainCmd

}
