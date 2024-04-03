package utility

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/hugobyte/dive/cli/common"
	"github.com/spf13/cobra"
)

var (
	latestVersion = ""
)

const versionFile = "version"

var VersionCmd = common.NewDiveCommandBuilder().
	SetUse("version").
	SetShort("Checks The DIVE CLI Version").
	SetLong("Checks the current DIVE CLI version and warns if you are using an old version.").
	SetRun(version).
	Build()

func version(cmd *cobra.Command, args []string) {

	cli := common.GetCli(common.EnclaveName)

	err := common.ValidateArgs(args)

	if err != nil {
		cli.Logger().SetErrorToStderr()
		cli.Logger().Fatal(common.CodeOf(err), err.Error())
		cli.Context().Exit(1)
	}
	latestVersion = GetLatestVersion(cli)
	currentVer, err := extractVersion(common.DiveVersion)
	if err != nil {
		cli.Error(common.WrapMessageToError(common.ErrInitializingCLI, err.Error()))
	}
	latestVer, err := extractVersion(latestVersion)
	if err != nil {
		cli.Error(common.WrapMessageToError(common.ErrInitializingCLI, err.Error()))
	}
	if currentVer < latestVer {
		cli.Logger().SetOutputToStdout()
		cli.Logger().Warnf("Update available '%s'. Get the latest version of our DIVE CLI for bug fixes, performance improvements, and new features.", latestVersion)
		version := color.New(color.Bold).Sprintf("CLI version - %s", common.DiveVersion)
		fmt.Println(version)
		cli.Context().Exit(0)
	}
	version := color.New(color.Bold).Sprintf("CLI version - %s", common.DiveVersion)
	fmt.Println(version)

}

func extractVersion(versionString string) (int, error) {
	// Remove the leading 'v' if present
	versionString = strings.TrimPrefix(versionString, "v")

	// Split the version string by the '-' delimiter (if it exists) and take the first part
	parts := strings.Split(versionString, "-")
	versionComponents := strings.Split(parts[0], ".")

	// Parse each component as an integer
	var versionInt int
	for _, component := range versionComponents {
		num, err := strconv.Atoi(component)
		if err != nil {
			return 0, err
		}
		versionInt = versionInt*100 + num
	}

	return versionInt, nil
}

// This function will fetch the latest version from HugoByte/Dive repo
func GetLatestVersion(cli *common.Cli) string {

	// Repo Name
	repo := "DIVE"
	owner := "HugoByte"

	versionFilePath, err := cli.FileHandler().GetAppDirPathOrAppFilePath(versionFile)
	if err != nil {
		cli.Logger().SetErrorToStderr()
		cli.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	versionFileInfo, err := os.Stat(versionFilePath)
	if os.IsNotExist(err) {

	} else if err != nil {
		cli.Logger().SetErrorToStderr()
		cli.Logger().Fatal(common.CodeOf(err), err.Error())

	}

	if versionFileInfo == nil || time.Since(versionFileInfo.ModTime()).Hours() > 1 {

		client := github.NewClient(nil)
		release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
		if err != nil {
			cli.Logger().SetErrorToStderr()
			cli.Logger().Error(common.CodeOf(err), err.Error())
			return ""
		}

		latestVersion = release.GetName()
		cli.FileHandler().WriteAppFile(versionFile, []byte(latestVersion))
		os.Chtimes(versionFilePath, time.Now(), time.Now())

	} else {

		cachedVersion, err := cli.FileHandler().ReadAppFile(versionFile)
		if err == nil && string(cachedVersion) != "" {
			latestVersion = string(cachedVersion)
		}
	}

	return latestVersion
}
