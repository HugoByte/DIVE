/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package commands

import (
	"os"

	"github.com/hugobyte/dive/cli/commands/bridge"
	"github.com/hugobyte/dive/cli/commands/chain"
	"github.com/hugobyte/dive/cli/commands/clean"
	"github.com/hugobyte/dive/cli/commands/discord"
	"github.com/hugobyte/dive/cli/commands/tutorial"
	"github.com/hugobyte/dive/cli/commands/twitter"
	"github.com/hugobyte/dive/cli/commands/version"
	"github.com/hugobyte/dive/cli/common"

	"github.com/hugobyte/dive/cli/styles"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	diveContext := common.NewDiveContext()
	var rootCmd = &cobra.Command{
		Use:   "dive",
		Short: "Deployable Infrastructure for Virtually Effortless blockchain integration",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			styles.RenderBanner()
			cmd.Help()

		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(chain.NewChainCmd(diveContext))
	rootCmd.AddCommand(bridge.NewBridgeCmd(diveContext))
	rootCmd.AddCommand(clean.NewCleanCmd(diveContext))
	rootCmd.AddCommand(version.NewVersionCmd(diveContext))
	rootCmd.AddCommand(discord.NewDiscordCmd(diveContext))
	rootCmd.AddCommand(twitter.NewtwitterCmd(diveContext))
	rootCmd.AddCommand(tutorial.NewTutorialCmd(diveContext))

	rootCmd.PersistentFlags().BoolVar(&common.DiveLogs, "verbose", false, "Prints out logs to Stdout")

	return rootCmd

}

func Execute() {
	err := RootCmd().Execute()

	if err != nil {
		os.Exit(1)
	}
}
