/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package chain

import (
	"context"

	"github.com/hugobyte/dive/commands/chain/types"
	"github.com/hugobyte/dive/common"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/spf13/cobra"
)

// chainCmd represents the chain command

var ChainCmd = &cobra.Command{
	Use:   "chain",
	Short: "runs specfied chain",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()

	},
}

func init() {

	ctx := context.Background()

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		panic(err)

	}

	enclaveCtx, err := getEnclaveContext(ctx, common.DiveEnclave, kurtosisContext)
	if err != nil {
		panic(err)

	}

	ChainCmd.AddCommand(types.NewIconCmd(ctx, enclaveCtx))
	ChainCmd.AddCommand(types.NewEthCmd(ctx, enclaveCtx))
	ChainCmd.AddCommand(types.NewHardhatCmd(ctx, enclaveCtx))

}

func getEnclaveContext(ctx context.Context, identifier string, kurtosisContext *kurtosis_context.KurtosisContext) (*enclaves.EnclaveContext, error) {

	_, err := kurtosisContext.GetEnclave(ctx, common.DiveEnclave)
	if err != nil {
		enclaveCtx, err := kurtosisContext.CreateEnclave(ctx, identifier, false)
		if err != nil {
			return nil, err

		}
		return enclaveCtx, nil
	}
	enclaveCtx, err := kurtosisContext.GetEnclaveContext(ctx, identifier)

	if err != nil {
		return nil, err
	}
	return enclaveCtx, nil
}
