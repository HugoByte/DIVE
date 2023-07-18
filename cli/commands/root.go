/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package commands

import (
	"os"

	"github.com/hugobyte/dive/commands/bridge"
	"github.com/hugobyte/dive/commands/chain"
	"github.com/hugobyte/dive/common"

	"github.com/hugobyte/dive/commands/clean"
	"github.com/hugobyte/dive/commands/discord"
	"github.com/hugobyte/dive/commands/tutorial"
	"github.com/hugobyte/dive/commands/twitter"
	"github.com/hugobyte/dive/commands/version"
	"github.com/hugobyte/dive/styles"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dive",
	Short: "Deployable Infrastructure for Virtually Effortless blockchain integration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		styles.RenderBanner()
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	diveContext := common.NewDiveContext()

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.AddCommand(chain.NewChainCmd(diveContext))
	rootCmd.AddCommand(bridge.NewBridgeCmd(diveContext))
	rootCmd.AddCommand(clean.NewCleanCmd(diveContext))
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(discord.DiscordCmd)
	rootCmd.AddCommand(twitter.TwitterCmd)
	rootCmd.AddCommand(tutorial.TutorialCmd)
}
