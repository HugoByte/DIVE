package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"

	log "github.com/sirupsen/logrus"
)

type DiveContext struct {
	Ctx             context.Context
	KurtosisContext *kurtosis_context.KurtosisContext
	Log             *log.Logger
	spinner         *spinner.Spinner
}

func NewDiveContext() *DiveContext {

	ctx := context.Background()
	log := setupLogger()

	spinner := spinner.New(spinner.CharSets[80], 100*time.Millisecond, spinner.WithWriter(os.Stderr))

	return &DiveContext{Ctx: ctx, Log: log, spinner: spinner}
}

func (diveContext *DiveContext) GetEnclaveContext() (*enclaves.EnclaveContext, error) {

	_, err := diveContext.KurtosisContext.GetEnclave(diveContext.Ctx, DiveEnclave)
	if err != nil {
		enclaveCtx, err := diveContext.KurtosisContext.CreateEnclave(diveContext.Ctx, DiveEnclave, false)
		if err != nil {
			return nil, err

		}
		return enclaveCtx, nil
	}
	enclaveCtx, err := diveContext.KurtosisContext.GetEnclaveContext(diveContext.Ctx, DiveEnclave)

	if err != nil {
		return nil, err
	}
	return enclaveCtx, nil
}

// To get names of running enclaves, returns empty string if no enclaves
func (diveContext *DiveContext) GetEnclaves() string {
	enclaves, err := diveContext.KurtosisContext.GetEnclaves(diveContext.Ctx)
	if err != nil {
		diveContext.Log.Errorf("Getting Enclaves failed with error:  %v", err)
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
func (diveContext *DiveContext) Clean() {
	diveContext.Log.SetOutput(os.Stdout)
	diveContext.Log.Info("Successfully connected to kurtosis engine...")
	diveContext.Log.Info("Initializing cleaning process...")
	// shouldCleanAll set to true as default for beta release.
	enclaves, err := diveContext.KurtosisContext.Clean(diveContext.Ctx, true)
	if err != nil {
		diveContext.Log.SetOutput(os.Stderr)
		diveContext.Log.Errorf("Failed cleaning with error: %v", err)
	}

	// Assuming only one enclave is running for beta release
	diveContext.Log.Infof("Successfully destroyed and cleaned enclave %s", enclaves[0].Name)
}

func (diveContext *DiveContext) FatalError(message, errorMessage string) {

	diveContext.Log.SetOutput(os.Stderr)
	diveContext.spinner.Stop()
	diveContext.Log.Fatalf("%s : %s", message, errorMessage)
}

func (diveContext *DiveContext) Info(message string) {

	diveContext.Log.Infoln(message)
}

func (diveContext *DiveContext) StartSpinner(message string) {

	diveContext.spinner.Suffix = message
	diveContext.spinner.Color("green")

	diveContext.spinner.Start()
}

func (diveContext *DiveContext) SetSpinnerMessage(message string) {

	diveContext.spinner.Suffix = message
}

func (diveContext *DiveContext) StopSpinner(message string) {
	c := color.New(color.FgCyan).Add(color.Underline)
	diveContext.spinner.FinalMSG = c.Sprintln(message)
	diveContext.spinner.Stop()

}

func (diveContext *DiveContext) GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error) {
	if DiveLogs {
		diveContext.spinner.Stop()
		diveContext.Log.SetOutput(os.Stdout)
	}
	var serializedOutputObj string
	services := map[string]string{}

	skippedInstruction := map[string]bool{}
	for executionResponse := range response {

		if strings.Contains(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), "added with service") {
			res1 := strings.Split(executionResponse.GetInstructionResult().GetSerializedInstructionResult(), " ")
			serviceName := res1[1][1 : len(res1[1])-1]
			serviceUUID := res1[len(res1)-1][1 : len(res1[len(res1)-1])-1]
			services[serviceName] = serviceUUID
		}

		diveContext.Log.Info(executionResponse.String())

		if executionResponse.GetInstruction().GetIsSkipped() {
			skippedInstruction[executionResponse.GetInstruction().GetExecutableInstruction()] = executionResponse.GetInstruction().GetIsSkipped()
			break
		}

		if executionResponse.GetError() != nil {

			return "", services, nil, errors.New(executionResponse.GetError().String())

		}

		runFinishedEvent := executionResponse.GetRunFinishedEvent()

		if runFinishedEvent != nil {

			if runFinishedEvent.GetIsRunSuccessful() {
				serializedOutputObj = runFinishedEvent.GetSerializedOutput()

			} else {
				return "", services, nil, errors.New(executionResponse.GetError().String())
			}

		} else {
			diveContext.spinner.Color("blue")
			if executionResponse.GetProgressInfo() != nil {
				c := color.New(color.FgGreen)
				diveContext.spinner.Suffix = c.Sprintf(strings.ReplaceAll(executionResponse.GetProgressInfo().String(), "current_step_info:", " "))

			}
		}

	}

	return serializedOutputObj, services, skippedInstruction, nil

}

func (diveContext *DiveContext) Error(err string) {
	diveContext.Log.Error(err)
}

func (diveContext *DiveContext) InitKurtosisContext() {
	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	if err != nil {
		diveContext.Log.SetOutput(os.Stderr)
		diveContext.Log.Fatal("The Kurtosis Engine Server is unavailable and is probably not running; you will need to start it using the Kurtosis CLI before you can create a connection to it")

	}
	diveContext.KurtosisContext = kurtosisContext
}

func (diveContext *DiveContext) CheckInstructionSkipped(instuctions map[string]bool, message string) {

	if len(instuctions) != 0 {

		diveContext.StopSpinner(message)
		os.Exit(0)
	}
}

func (diveContext *DiveContext) StopServices(services map[string]string) {
	if len(services) != 0 {
		enclaveContext, err := diveContext.KurtosisContext.GetEnclaveContext(diveContext.Ctx, DiveEnclave)

		if err != nil {
			diveContext.Log.Fatal("Failed To Retrieve Enclave Context", err)
		}

		for serviceName, serviceUUID := range services {
			params := fmt.Sprintf(`{"service_name": "%s", "uuid": "%s"}`, serviceName, serviceUUID)
			_, err := enclaveContext.RunStarlarkScriptBlocking(diveContext.Ctx, "", starlarkScript, params, DiveDryRun, DiveDefaultParallelism, nil)
			if err != nil {
				diveContext.Log.Fatal("Failed To Stop Services", err)
			}

		}

	}

}
