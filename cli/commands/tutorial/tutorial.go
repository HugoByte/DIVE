/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package tutorial

import (
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const tutorialURL = "https://www.youtube.com/@hugobyte"

// tutorilaCmd redirects users to DIVE youtube playlist
var TutorialCmd = &cobra.Command{
	Use:   "tutorial",
	Short: "Opens DIVE tutorial youtube playlist",
	Long: `The command opens the YouTube playlist containing DIVE tutorials. It launches a web browser or the YouTube application,
directing users to a curated collection of tutorial videos specifically designed to guide and educate users about DIVE. The playlist 
offers step-by-step instructions, tips, and demonstrations to help users better understand and utilize the features and functionalities of DIVE.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Redirecting to YouTube...")
		if err := common.OpenFile(tutorialURL); err != nil {
			logrus.Errorf("Failed to open Dive YouTube chanel with error %v", err)
		}
	},
}
