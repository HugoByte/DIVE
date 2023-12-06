package utility

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const tutorialURL = "https://www.youtube.com/playlist?list=PL5Xd9z-fRL1vKtRlOzIlkhROspSSDeGyG"

var TutorialCmd = common.NewDiveCommandBuilder().
	SetUse("tutorial").
	SetShort("Opens DIVE Tutorial Youtube Playlist").
	SetLong(
		`The command opens the YouTube playlist containing DIVE tutorials. It launches a web browser or the YouTube application,
	directing users to a curated collection of tutorial videos specifically designed to guide and educate users about DIVE. The playlist 
	offers step-by-step instructions, tips, and demonstrations to help users better understand and utilize the features and functionalities of DIVE.`,
	).
	SetRun(tutorial).
	Build()

func tutorial(cmd *cobra.Command, args []string) {

	cli := common.GetCli()

	cli.Logger().SetOutputToStdout()

	err := common.ValidateArgs(args)

	if err != nil {
		cli.Logger().Error(common.CodeOf(err), common.Errorcf(common.CodeOf(err), "error %s !!! \n%s", err, cmd.UsageString()).Error())
	}

	cli.Logger().Info("Redirecting to YouTube.....")

	if err := common.OpenFile(tutorialURL); err != nil {
		cli.Logger().Fatalf(common.InvalidCommandError, common.Errorcf(common.CodeOf(err), "Failed to open Dive YouTube chanel with error %v", err).Error())
	}

}
