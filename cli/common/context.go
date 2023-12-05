package common

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

const (
	removeServiceStarlarkScript = `
def run(plan,args):
		plan.remove_service(name=args["service_name"])
`

	stopServiceStarlarkScript = `
def run(plan, args):
	plan.stop_service(name=args["service_name"])
`
)

var (
	kurtosisContextErr error
)

type diveContext struct {
	mu               sync.Mutex
	ctx              context.Context
	kurtosisContext  *kurtosis_context.KurtosisContext
	kurtosisInitDone bool
}

func NewDiveContext1() *diveContext {

	return &diveContext{ctx: context.Background()}
}
func (dc *diveContext) GetContext() context.Context {
	return dc.ctx
}

func (dc *diveContext) GetKurtosisContext() (*kurtosis_context.KurtosisContext, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if !dc.kurtosisInitDone {
		kurtosisContext, err := dc.initKurtosisContext()
		if err != nil {
			return nil, err
		}
		dc.kurtosisContext = kurtosisContext
		dc.kurtosisInitDone = true
	}
	return dc.kurtosisContext, nil
}

func (dc *diveContext) GetEnclaves() ([]EnclaveInfo, error) {
	enclavesInfo, err := dc.kurtosisContext.GetEnclaves(dc.ctx)

	if err != nil {
		return nil, Errorc(KurtosisContextError, err.Error())
	}
	var enclaves []EnclaveInfo
	enclaveMap := enclavesInfo.GetEnclavesByName()
	for _, enclaveInfoList := range enclaveMap {
		for _, enclaveInfo := range enclaveInfoList {

			enclaves = append(enclaves, EnclaveInfo{Name: enclaveInfo.Name, Uuid: enclaveInfo.EnclaveUuid, ShortUuid: enclaveInfo.ShortenedUuid})
		}
	}
	return enclaves, nil
}

func (dc *diveContext) GetEnclaveContext(enclaveName string) (*enclaves.EnclaveContext, error) {
	// check enclave exist
	enclaveInfo, err := dc.checkEnclaveExist(enclaveName)
	if err != nil {
		return dc.CreateEnclave(enclaveName)
	}
	enclaveCtx, err := dc.kurtosisContext.GetEnclaveContext(dc.ctx, enclaveInfo.Name)
	if err != nil {
		return nil, Errorc(InvalidEnclaveNameError, err.Error())
	}

	return enclaveCtx, nil

}

func (dc *diveContext) CleanEnclaves() ([]*EnclaveInfo, error) {
	enclaves, err := dc.kurtosisContext.Clean(dc.ctx, true)

	if err != nil {
		return nil, Errorc(InvalidEnclaveContextError, err.Error())
	}

	var enclaveInfo []*EnclaveInfo

	for _, info := range enclaves {
		enclaveInfo = append(enclaveInfo, &EnclaveInfo{Name: info.Name, Uuid: info.Uuid})
	}

	return enclaveInfo, nil
}

func (dc *diveContext) CleanEnclaveByName(enclaveName string) error {

	enclaveInfo, err := dc.checkEnclaveExist(enclaveName)

	if err != nil {
		return Errorc(KurtosisContextError, err.Error())
	}
	err = dc.kurtosisContext.DestroyEnclave(dc.ctx, enclaveInfo.Uuid)
	if err != nil {
		return Errorc(KurtosisContextError, err.Error())
	}
	return nil
}

func (dc *diveContext) CheckSkippedInstructions(instructions map[string]bool) bool {

	return len(instructions) != 0
}

func (dc *diveContext) StopService(serviceName string, enclaveName string) error {

	enclaveContext, err := dc.GetEnclaveContext(enclaveName)
	if err != nil {
		return err
	}

	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
	starlarkConfig := GetStarlarkRunConfig(params, "", "")

	_, err = enclaveContext.RunStarlarkScriptBlocking(dc.ctx, stopServiceStarlarkScript, starlarkConfig)
	if err != nil {
		return err
	}
	return nil
}

func (dc *diveContext) StopServices(enclaveName string) error {

	enclaveCtx, err := dc.GetEnclaveContext(enclaveName)

	if err != nil {
		return WrapMessageToError(err, "Failed To Stop Services")
	}

	services, err := enclaveCtx.GetServices()

	if err != nil {
		return WrapMessageToError(err, "Failed To Stop Services")
	}

	for serviceName := range services {
		params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
		starlarkConfig := GetStarlarkRunConfig(params, "", "")

		_, err = enclaveCtx.RunStarlarkScriptBlocking(dc.ctx, stopServiceStarlarkScript, starlarkConfig)
		if err != nil {
			return err
		}

	}

	return nil
}

func (dc *diveContext) RemoveServices(enclaveName string) error {
	enclaveCtx, err := dc.GetEnclaveContext(enclaveName)

	if err != nil {
		return WrapMessageToError(err, "Failed To Remove Services")
	}

	services, err := enclaveCtx.GetServices()

	if err != nil {
		return WrapMessageToError(err, "Failed To Remove Services")
	}

	for serviceName := range services {
		params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
		starlarkConfig := GetStarlarkRunConfig(params, "", "")

		_, err = enclaveCtx.RunStarlarkScriptBlocking(dc.ctx, removeServiceStarlarkScript, starlarkConfig)
		if err != nil {
			return err
		}

	}

	return nil
}

func (dc *diveContext) RemoveService(serviceName string, enclaveName string) error {
	enclaveContext, err := dc.GetEnclaveContext(enclaveName)
	if err != nil {
		return err
	}

	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
	starlarkConfig := GetStarlarkRunConfig(params, "", "")

	_, err = enclaveContext.RunStarlarkScriptBlocking(dc.ctx, removeServiceStarlarkScript, starlarkConfig)
	if err != nil {
		return err
	}
	return nil
}

func (dc *diveContext) CreateEnclave(enclaveName string) (*enclaves.EnclaveContext, error) {

	enclaveContext, err := dc.kurtosisContext.CreateEnclave(dc.ctx, enclaveName)

	if err != nil {
		return nil, Errorc(InvalidEnclaveContextError, err.Error())
	}

	return enclaveContext, nil
}

func (dc *diveContext) initKurtosisContext() (*kurtosis_context.KurtosisContext, error) {

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()

	if err != nil {
		return nil, WrapMessageToError(ErrKurtosisContext, err.Error())

	}

	return kurtosisContext, nil
}
func (dc *diveContext) checkEnclaveExist(enclaveName string) (*EnclaveInfo, error) {

	kurtosisContext, err := dc.GetKurtosisContext()
	if err != nil {
		return nil, Errorc(KurtosisContextError, err.Error())
	}
	enclaveInfo, err := kurtosisContext.GetEnclave(dc.ctx, enclaveName)
	if err != nil {
		return nil, Errorc(KurtosisContextError, err.Error())
	}

	return &EnclaveInfo{
		Name:      enclaveInfo.Name,
		Uuid:      enclaveInfo.EnclaveUuid,
		ShortUuid: enclaveInfo.ShortenedUuid,
	}, nil
}

func GetStarlarkRunConfig(params string, relativePathToMainFile string, mainFunctionName string) *starlark_run_config.StarlarkRunConfig {

	starlarkConfig := &starlark_run_config.StarlarkRunConfig{
		RelativePathToMainFile:   relativePathToMainFile,
		MainFunctionName:         mainFunctionName,
		DryRun:                   DiveDryRun,
		SerializedParams:         params,
		Parallelism:              DiveDefaultParallelism,
		ExperimentalFeatureFlags: []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{},
	}
	return starlarkConfig
}

func GetSerializedData(cliContext *Cli, response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error) {

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

		cliContext.log.Info(executionResponse.String())

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
			cliContext.spinner.SetColor("blue")
			if executionResponse.GetProgressInfo() != nil {

				cliContext.spinner.SetSuffixMessage(strings.ReplaceAll(executionResponse.GetProgressInfo().String(), "current_step_info:", " "), "fgGreen")

			}
		}

	}

	return serializedOutputObj, services, skippedInstruction, nil
}
