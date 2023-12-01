package utility

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var CleanCmd = common.NewDiveCommandBuilder().
	SetUse("clean").
	SetShort("Cleans up Kurtosis leftover artifacts").
	SetLong("Destroys and removes any running encalves. If no enclaves running to remove it will throw an error").
	SetRun(clean).Build()

func clean(cmd *cobra.Command, args []string) {
}
