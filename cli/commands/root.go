/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package commands

import (
	"os"

	"github.com/hugobyte/dive/commands/bridge"
	"github.com/hugobyte/dive/commands/chain"

	"github.com/hugobyte/dive/commands/clean"
	"github.com/hugobyte/dive/commands/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dive",
	Short: "Deployable Infrastructure for Virtually Effortless blockchain integration",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// type DiveContext struct {
// 	Ctx             context.Context
// 	KurtosisContext *kurtosis_context.KurtosisContext
// 	Log             *logrus.Logger
// }

func init() {

	// kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	// if err != nil {
	// 	panic(err)
	// }
	// diveContext := DiveContext{
	// 	Ctx:             context.Background(),
	// 	KurtosisContext: kurtosisContext,
	// 	Log:             logrus.New(),
	// }

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true

	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(chain.ChainCmd)
	rootCmd.AddCommand(bridge.BridgeCmd)
	rootCmd.AddCommand(clean.CleanCmd)
}