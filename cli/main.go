/*
Copyright © 2023 Hugobyte AI Labs<hello@hugobyte.com>
*/
package main

import (
	"github.com/hugobyte/dive/commands"
	"github.com/hugobyte/dive/styles"
	// "github.com/sirupsen/logrus"
)

func main() {
	// logrus.Println("hello")
	styles.RenderBanner()
	commands.Execute()

}

// ctx, cancelCtxFunc := context.WithCancel(context.Background())
// 	defer cancelCtxFunc()

// 	kurtosisCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Println(kurtosisCtx)

// 	enclaveId := fmt.Sprintf("%s-%d", enclaveIdPrefix, time.Now().Unix())

// 	enclaveCtx, err := kurtosisCtx.CreateEnclave(ctx, enclaveId, isPartitioningEnabled)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	data, _, err := enclaveCtx.RunStarlarkRemotePackage(ctx, divePackage, pathToMainFile, mainFunctionName, emptyPackageParams, noDryRun, defaultParallelism)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	var serializedOutputObj string
// 	for executionResponseLine := range data {
// 		runFinishedEvent := executionResponseLine.GetRunFinishedEvent()
// 		if runFinishedEvent == nil {
// 			logrus.Info("Execution in progress...")
// 		} else {
// 			logrus.Info("Execution finished successfully")
// 			if runFinishedEvent.GetIsRunSuccessful() {
// 				serializedOutputObj = runFinishedEvent.GetSerializedOutput()
// 			} else {
// 				panic("Starlark run failed")
// 			}
// 		}
// 	}

// 	fmt.Println(serializedOutputObj)

// "github.com/sirupsen/logrus"

// const (
// 	solidityContractsGitPath = ""
// 	javaContractGitPath      = ""
// )

// const (
// 	enclaveIdPrefix       = "quick-start-go-example"
// 	isPartitioningEnabled = false

// 	divePackage = "github.com/hugobyte/dive"

// 	defaultParallelism = 4
// 	noDryRun           = false

// 	emptyPackageParams = `{"args":{},"node_name":"icon"}`

// 	apiServiceName = "api"

// 	contentType = "application/json"

// 	pathToMainFile   = "main.star"
// 	mainFunctionName = "run_node"
// )
