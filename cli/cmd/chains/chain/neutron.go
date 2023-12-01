package chain

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var NeutronCmd = common.NewDiveCommandBuilder().
	SetUse("neutron").
	SetShort("Build, initialize and start a neutron node").
	SetLong("The command starts the neutron network and allows node in executing contracts").
	SetRun(neutron).
	Build()

func neutron(cmd *cobra.Command, args []string) {}
