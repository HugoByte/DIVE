package chain

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var EthCmd = common.NewDiveCommandBuilder().
	SetUse("eth").
	SetShort("Build, initialize and start a eth node.").
	SetLong(`The command starts an Ethereum node, initiating the process of setting up and launching a local Ethereum network. 
It establishes a connection to the Ethereum network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	SetRun(eth).
	Build()

func eth(cmd *cobra.Command, args []string) {}
