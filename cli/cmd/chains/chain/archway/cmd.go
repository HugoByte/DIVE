package archway

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var ArchwayCmd = common.NewDiveCommandBuilder().
	SetUse("archway").
	SetShort("Build, initialize and start a archway node").
	SetLong("The command starts the archway network and allows node in executing contracts").
	SetRun(archway).
	Build()

func archway(cmd *cobra.Command, args []string) {}
