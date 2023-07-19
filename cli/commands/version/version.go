/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package version

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the CLI version",
	Long:  `Prints the current DIVE CLI version and warns if you are using an old version.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.ValidateCmdArgs(args, cmd.UsageString())

		// Checks for latest Version
		latestVersion := common.GetLatestVersion()
		if common.DiveVersion != latestVersion {
			logrus.Warnf("Update available '%s'. Get the latest version of our DIVE CLI for bug fixes, performance improvements, and new features.", latestVersion)
		}
		version := color.New(color.Bold).Sprintf("CLI version - %s", common.DiveVersion)
		fmt.Println(version)

	},
}
