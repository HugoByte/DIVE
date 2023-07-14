/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package discord

import (
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const diveURL = "https://discord.com/channels/1097522975630184469/1124224608250376293"

// discordCmd redirects users to DIVE discord channel
var DiscordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Opens DIVE discord channel",
	Long:  `Redirects users to DIVE discord channel`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Redirecting to DIVE discord channel...")
		if err := common.OpenFile(diveURL); err != nil {
			logrus.Errorf("Failed to open Dive discord channel with error %v", err)
		}
	},
}
