package cmd

import (
	"os"

	"github.com/hugobyte/dive-core/cli/cmd/social"
	"github.com/hugobyte/dive-core/cli/cmd/utility"
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/hugobyte/dive-core/cli/styles"
	"github.com/spf13/cobra"
)

var rootCmd = common.NewDiveCommandBuilder().
	SetUse("dive").
	SetShort("Deployable Infrastructure for Virtually Effortless blockchain integration").
	ToggleHelpCommand(true).
	AddCommand(utility.CleanCmd).
	AddCommand(utility.TutorialCmd).
	AddCommand(utility.VersionCmd).
	AddCommand(social.DiscordCmd).
	AddCommand(social.TwitterCmd).
	AddBoolPersistentFlag(&common.DiveLogs, "verbose", false, "Prints out logs to Stdout").
	SetRunE(run).
	Build()

func run(cmd *cobra.Command, args []string) error {
	styles.RenderBanner()
	cmd.Help()

	return nil
}

func init() {
	common.GetDiveContext()
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}
