package hardhat

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var HardhatCmd = common.NewDiveCommandBuilder().
	SetUse("hardhat").
	SetShort("Build, initialize and start a hardhat node.").
	SetLong(`The command starts an hardhat node, initiating the process of setting up and launching a local hardhat network. 
It establishes a connection to the hardhat network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	SetRun(hardhat).
	Build()

func hardhat(cmd *cobra.Command, args []string) {}
