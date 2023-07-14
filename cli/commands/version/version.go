/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package version

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the CLI version",
	Long:  `Prints the current DIVE CLI version and warns if you are using an old version.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := color.New(color.Bold).Sprint("CLI version - ") + DiveVersion
		latestVersion := getLatestVersion()
		if DiveVersion != latestVersion {
			logrus.Warnf("Update available '%s'. Get the latest version of our DIVE CLI for bug fixes, performance improvements, and new features.", latestVersion)
		}
		fmt.Println(version)

	},
}

// This function will fetch the latest version from HugoByte/Dive repo
func getLatestVersion() string {

	// Repo Name
	repo := "DIVE"
	owner := "HugoByte"

	// Create a new github client
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Print the release version.
	return release.GetName()
}
