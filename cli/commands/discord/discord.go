/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package discord

import (
	"os/exec"
	"runtime"

	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	linuxOSName   = "linux"
	macOSName     = "darwin"
	windowsOSName = "windows"

	openFileLinuxCommandName   = "xdg-open"
	openFileMacCommandName     = "open"
	openFileWindowsCommandName = "rundll32"

	openFileWindowsCommandFirstArgumentDefault = "url.dll,FileProtocolHandler"

	diveURL = "https://discord.com/channels/1097522975630184469/1124224608250376293"
)

// discordCmd redirects users to DIVE discord channel
var DiscordCmd = &cobra.Command{
	Use:   "discord",
	Short: "Opens DIVE discord channel",
	Long:  `Redirects users to DIVE discord channel`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := OpenFile(diveURL); err != nil {
			logrus.Errorf("Failed to open Dive discord channel with error %v", err)
		}
	},
}

func OpenFile(URL string) error {
	var args []string
	switch runtime.GOOS {
	case linuxOSName:
		args = []string{openFileLinuxCommandName, URL}
	case macOSName:
		args = []string{openFileMacCommandName, URL}
	case windowsOSName:
		args = []string{openFileWindowsCommandName, openFileWindowsCommandFirstArgumentDefault, URL}
	default:
		return stacktrace.NewError("Unsupported operating system")
	}
	command := exec.Command(args[0], args[1:]...)
	if err := command.Start(); err != nil {
		return stacktrace.Propagate(err, "An error occurred while opening '%v'", URL)
	}
	return nil
}
