/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package clean

import (
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans up Kurtosis leftover artifacts",
	Long:  `Destroys and removes any running encalves. If no enclaves running to remove it will throw an error`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
