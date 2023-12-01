package utility

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var VersionCmd = common.NewDiveCommandBuilder().
	SetUse("version").
	SetShort("Checks The DIVE CLI Version").
	SetLong("Checks the current DIVE CLI version and warns if you are using an old version.").
	SetRun(version).
	Build()

func version(cmd *cobra.Command, args []string) {}
