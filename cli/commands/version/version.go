/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package version

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hugobyte/dive/common"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the CLI version",
	Long:  `Prints the current DIVE CLI version and warns if you are using an old version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := color.New(color.Bold).Sprint("CLI version - ") + common.DiveVersion
		fmt.Println(version)

	},
}
