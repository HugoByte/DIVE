/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package clean

import (
	"context"

	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans up Kurtosis leftover artifacts",
	Long:  `Destroys and removes any running encalves. If no enclaves running to remove it will throw an error`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
		logrus.Info("Trying to connect to local Kurtosis Engine...")
		if err != nil {
			logrus.Errorf("Connecting to kurtosis engine failed as kurtosis engine is not running")
		}
		enclaveName := getEnclaves(ctx, kurtosisCtx)
		if enclaveName == "" {
			logrus.Errorf("No enclaves running to clean !!")
		} else {
			clean(ctx, kurtosisCtx)
		}
	},
}

// To get names of running enclaves, returns empty string if no enclaves
func getEnclaves(ctx context.Context, kurtosisCtx *kurtosis_context.KurtosisContext) string {
	enclaves, err := kurtosisCtx.GetEnclaves(ctx)
	if err != nil {
		logrus.Errorf("Getting Enclaves failed with error:  %v", err)
	}
	enclaveMap := enclaves.GetEnclavesByName()
	for _, enclaveInfoList := range enclaveMap {
		for _, enclaveInfo := range enclaveInfoList {
			return enclaveInfo.GetName()
		}
	}
	return ""
}

// Funstionality to clean the enclaves
func clean(ctx context.Context, kurtosisCtx *kurtosis_context.KurtosisContext) {
	logrus.Info("Successfully connected to kurtosis engine...")
	logrus.Info("Initializing cleaning process...")

	// shouldCleanAll set to true as default for beta release.
	enclaves, err := kurtosisCtx.Clean(ctx, true)
	if err != nil {
		logrus.Errorf("Failed cleaning with error: %v", err)
	}

	// Assuming only one enclave is running for beta release
	logrus.Infof("Successfully destroyed and cleaned enclave %s", enclaves[0].Name)
}
