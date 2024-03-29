package chains

import (
	"os"
	"slices"

	"github.com/hugobyte/dive/cli/cmd/chains/archway"
	"github.com/hugobyte/dive/cli/cmd/chains/eth"
	"github.com/hugobyte/dive/cli/cmd/chains/hardhat"
	"github.com/hugobyte/dive/cli/cmd/chains/icon"
	"github.com/hugobyte/dive/cli/cmd/chains/kusama"
	"github.com/hugobyte/dive/cli/cmd/chains/neutron"
	"github.com/hugobyte/dive/cli/cmd/chains/polkadot"
	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

var ChainCmd = common.NewDiveCommandBuilder().
	SetUse("chain").
	SetShort("Build, initialize and start a given blockchain node").
	SetLong(`The command builds, initializes, and starts a specified blockchain node, providing a seamless setup process.
It encompasses compiling and configuring the necessary dependencies and components required for the blockchain network. 
By executing this command, the node is launched, enabling network participation, transaction processing, and ledger 
maintenance within the specified blockchain ecosystem.`,
	).
	AddCommand(icon.IconCmd).
	AddCommand(eth.EthCmd).
	AddCommand(hardhat.HardhatCmd).
	AddCommand(archway.ArchwayCmd).
	AddCommand(neutron.NeutronCmd).
	AddCommand(polkadot.PolkadotCmd).
	AddCommand(kusama.KusamaCmd).
	SetRun(chains).
	Build()

func chains(cmd *cobra.Command, args []string) {

	validArgs := cmd.ValidArgs
	for _, c := range cmd.Commands() {
		validArgs = append(validArgs, c.Name())
	}
	cmd.ValidArgs = validArgs

	if len(args) == 0 {
		cmd.Help()

	} else if !slices.Contains(cmd.ValidArgs, args[0]) {

		cmd.Usage()
		os.Exit(1)
	}
}
