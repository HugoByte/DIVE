package common

import (
	"context"
	"fmt"
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
	Once               sync.Once
	kurtosisContextErr error
)

type diveContext struct {
	ctx             context.Context
	kurtosisContext *kurtosis_context.KurtosisContext
}

func NewDiveContext1() *diveContext {

	return &diveContext{ctx: context.Background()}
}
func (dc *diveContext) GetContext() context.Context {
	return dc.ctx
}

func (dc *diveContext) GetKurtosisContext() (*kurtosis_context.KurtosisContext, error) {
	once.Do(func() {
		err := dc.initKurtosisContext()
		if err != nil {
			kurtosisContextErr = err
		}
	})
	if kurtosisContextErr != nil {
		return nil, kurtosisContextErr
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

func (dc *diveContext) CheckSkippedInstructions() {
	panic("not implemented") // TODO: Implement
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

func (dc *diveContext) GetSerializedData(response chan *kurtosis_core_rpc_api_bindings.StarlarkRunResponseLine) (string, map[string]string, map[string]bool, error) {
	panic("not implemented") // TODO: Implement
}

func (dc *diveContext) initKurtosisContext() error {

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()

	if err != nil {
		return WrapMessageToError(ErrKurtosisContext, err.Error())

	}

	dc.kurtosisContext = kurtosisContext

	return nil
}
func (dc *diveContext) checkEnclaveExist(enclaveName string) (*EnclaveInfo, error) {

	enclaveInfo, err := dc.kurtosisContext.GetEnclave(dc.ctx, enclaveName)
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
