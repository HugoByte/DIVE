package utility

import (
	"fmt"

	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

var CleanCmd = common.NewDiveCommandBuilder().
	SetUse("clean").
	SetShort("Cleans up Kurtosis leftover artifacts").
	SetLong("Destroys and removes any running encalves. If no enclaves running to remove it will throw an error").
	AddBoolFlagP("all", "a", false, "To Clean All the Service in Enclave").
	SetRun(clean).Build()

func clean(cmd *cobra.Command, args []string) {
	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	cleanAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Error(common.InvalidCommandError, err.Error())
	}

	if err != nil {
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	err = cliContext.FileHandler().RemoveFiles([]string{common.DiveOutFile, common.ServiceFilePath})
	if err != nil {
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	enclaves, err := cliContext.Context().GetEnclaves()
	if err != nil {
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	if len(enclaves) == 0 {
		cliContext.Logger().SetOutputToStdout()
		cliContext.Logger().Info("No Enclaves Running")
		cliContext.Context().Exit(0)

	}

	if cleanAll {
		cliContext.StartSpinnerIfNotVerbose("Cleaning All Dive Enclaves", common.DiveLogs)
		enclavesInfo, err := cliContext.Context().CleanEnclaves()
		if err != nil {
			cliContext.Logger().SetErrorToStderr()
			cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
		}

		cliContext.Logger().Info(fmt.Sprintf("Enclaves Cleaned %v", enclavesInfo))

	} else {
		cliContext.StartSpinnerIfNotVerbose(fmt.Sprintf("Cleaning Dive By Enclave %s", common.EnclaveName), common.DiveLogs)
		err = cliContext.Context().CleanEnclaveByName(common.EnclaveName)
		if err != nil {
			cliContext.Logger().SetErrorToStderr()
			cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
		}
	}

	cliContext.StopSpinnerIfNotVerbose("Clean Completed", common.DiveLogs)
}
