/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package clean

import (
	"github.com/hugobyte/dive/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCleanCmd(diveContext *common.DiveContext) *cobra.Command {

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Cleans up Kurtosis leftover artifacts",
		Long:  `Destroys and removes any running encalves. If no enclaves running to remove it will throw an error`,
		Run: func(cmd *cobra.Command, args []string) {
			common.ValidateCmdArgs(args, cmd.UsageString())
			enclaveName := diveContext.GetEnclaves()
			if enclaveName == "" {
				logrus.Errorf("No enclaves running to clean !!")
			} else {
				diveContext.Clean()
			}
		},
	}

	return cleanCmd

}
