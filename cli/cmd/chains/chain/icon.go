package chain

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var IconCmd = common.NewDiveCommandBuilder().
	SetUse("icon").
	SetShort("Build, initialize and start a icon node.").
	SetLong(`The command starts an Icon node, initiating the process of setting up and launching a local Icon network.
It establishes a connection to the Icon network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	SetRun(icon).
	Build()

func icon(cmd *cobra.Command, args []string) {}
