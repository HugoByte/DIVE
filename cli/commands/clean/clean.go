/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package clean

import (
	"fmt"
	"os"

	"github.com/hugobyte/dive/common"
	"github.com/spf13/cobra"
)

func NewCleanCmd(diveContext *common.DiveContext) *cobra.Command {

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Cleans up Kurtosis leftover artifacts",
		Long:  `Destroys and removes any running encalves. If no enclaves running to remove it will throw an error`,
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(args, cmd.UsageString())

			diveContext.InitKurtosisContext()
			pwd, err := os.Getwd()

			if err != nil {
				diveContext.FatalError("Failed cleaning with error: %v", err.Error())
			}

			diveOutPath := fmt.Sprintf("%s/%s", pwd, common.DiveOutFile)

			_, err = os.Stat(diveOutPath)

			if err == nil {
				os.Remove(diveOutPath)
			}
			servicesOutPath := fmt.Sprintf("%s/%s", pwd, common.ServiceFilePath)
			_, err = os.Stat(servicesOutPath)

			if err == nil {
				os.Remove(servicesOutPath)
			}

			enclaveName := diveContext.GetEnclaves()
			if enclaveName == "" {
				diveContext.Log.SetOutput(os.Stderr)
				diveContext.Error("No enclaves running to clean !!")

			} else {
				diveContext.Log.SetOutput(os.Stdout)
				diveContext.Clean()
			}
		},
	}

	return cleanCmd

}
