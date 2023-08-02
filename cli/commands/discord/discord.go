/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package discord

import (
	"os"

	"github.com/hugobyte/dive/common"
	"github.com/spf13/cobra"
)

const diveURL = "https://discord.gg/GyRQSBN3Cu"

// discordCmd redirects users to DIVE discord channel
func NewDiscordCmd(diveContext *common.DiveContext) *cobra.Command {
	return &cobra.Command{
		Use:   "discord",
		Short: "Opens DIVE discord channel",
		Long: `The command opens the Discord channel for DIVE, providing a direct link or launching the Discord application
to access the dedicated DIVE community. It allows users to engage in discussions, seek support, share insights, and 
collaborate with other members of the DIVE community within the Discord platform.`,
		Run: func(cmd *cobra.Command, args []string) {
			diveContext.Log.SetOutput(os.Stdout)
			common.ValidateCmdArgs(args, cmd.UsageString())
			diveContext.Log.Info("Redirecting to DIVE discord channel...")
			if err := common.OpenFile(diveURL); err != nil {
				diveContext.Log.Errorf("Failed to open Dive discord channel with error %v", err)
			}
		},
	}

}
