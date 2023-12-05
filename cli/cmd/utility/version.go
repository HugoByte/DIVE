package utility

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/hugobyte/dive-core/cli/common"
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

	cli := common.GetCli()

	err := common.ValidateArgs(args)

	if err != nil {
		cli.Logger().SetErrorToStderr()
		cli.Logger().Fatal(common.CodeOf(err), err.Error())
		cli.Context().Exit(1)
	}

	fmt.Println(GetLatestVersion(cli))

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
