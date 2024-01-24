/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package tutorial

import (
	"os"

	"github.com/hugobyte/dive-alpha/cli/common"
	"github.com/spf13/cobra"
)

const tutorialURL = "https://www.youtube.com/playlist?list=PL5Xd9z-fRL1vKtRlOzIlkhROspSSDeGyG"

// tutorilaCmd redirects users to DIVE youtube playlist
func NewTutorialCmd(diveContext *common.DiveContext) *cobra.Command {

	return &cobra.Command{
		Use:   "tutorial",
		Short: "Opens DIVE tutorial youtube playlist",
		Long: `The command opens the YouTube playlist containing DIVE tutorials. It launches a web browser or the YouTube application,
directing users to a curated collection of tutorial videos specifically designed to guide and educate users about DIVE. The playlist 
offers step-by-step instructions, tips, and demonstrations to help users better understand and utilize the features and functionalities of DIVE.`,
		Run: func(cmd *cobra.Command, args []string) {
			diveContext.Log.SetOutput(os.Stdout)
			common.ValidateCmdArgs(diveContext, args, cmd.UsageString())
			diveContext.Log.Info("Redirecting to YouTube...")
			if err := common.OpenFile(tutorialURL); err != nil {
				diveContext.Log.Errorf("Failed to open Dive YouTube chanel with error %v", err)
			}
		},
	}
}
