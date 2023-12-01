package social

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const diveURL = "https://discord.gg/GyRQSBN3Cu"

var DiscordCmd = common.NewDiveCommandBuilder().
	SetUse("discord").
	SetShort("Opens DIVE discord channel").
	SetLong(`The command opens the Discord channel for DIVE, providing a direct link or launching the Discord application to access the dedicated DIVE community. It allows users to engage in discussions, seek support, share insights, and collaborate with other members of the DIVE community within the Discord platform.`).
	SetRun(discord).Build()

func discord(cmd *cobra.Command, args []string) {}
