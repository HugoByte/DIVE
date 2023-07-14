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
	Short: "Build, initialize and start a given blockchain node.",
	Long:  `The command builds, initializes, and starts a specified blockchain node, providing a seamless setup process. It encompasses compiling and configuring the
necessary dependencies and components required for the blockchain network. By executing this command, the node is launched, enabling network participation, transaction
processing, and ledger maintenance within the specified blockchain ecosystem.`,

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
