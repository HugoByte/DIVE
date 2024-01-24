/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package version

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hugobyte/dive-alpha/cli/common"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func NewVersionCmd(diveContext *common.DiveContext) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the CLI version",
		Long:  `Prints the current DIVE CLI version and warns if you are using an old version.`,
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())
			diveContext.Log.SetOutput(os.Stdout)
			// Checks for latest Version
			latestVersion := common.GetLatestVersion()
			if common.DiveVersion != latestVersion {
				diveContext.Log.Warnf("Update available '%s'. Get the latest version of our DIVE CLI for bug fixes, performance improvements, and new features.", latestVersion)
			}
			version := color.New(color.Bold).Sprintf("CLI version - %s", common.DiveVersion)
			fmt.Println(version)

		},
	}

}
