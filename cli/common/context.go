package common

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
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
			return nil, WrapMessageToError(ErrInitializingKurtosis, err.Error())
		}
		dc.kurtosisContext = kurtosisContext
		dc.kurtosisInitDone = true
	}
	return dc.kurtosisContext, nil
}

func (dc *diveContext) GetEnclaves() ([]EnclaveInfo, error) {
	enclavesInfo, err := dc.kurtosisContext.GetEnclaves(dc.ctx)

	if err != nil {
		return nil, WrapMessageToError(ErrInvalidEnclaveContext, err.Error())
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
		return nil, WrapMessageToError(ErrEnclaveNameInvalid, err.Error())
	}

	return enclaveCtx, nil

}

func (dc *diveContext) CleanEnclaves() ([]*EnclaveInfo, error) {
	enclaves, err := dc.kurtosisContext.Clean(dc.ctx, true)

	if err != nil {
		return nil, WrapMessageToError(ErrInvalidKurtosisContext, err.Error())
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
		return WrapMessageToError(ErrInvalidKurtosisContext, err.Error())
	}
	err = dc.kurtosisContext.DestroyEnclave(dc.ctx, enclaveInfo.Uuid)
	if err != nil {
		return WrapMessageToError(ErrInvalidKurtosisContext, err.Error())
	}
	return nil
}

func (dc *diveContext) CheckSkippedInstructions(instructions map[string]bool) bool {

	return len(instructions) != 0
}

func (dc *diveContext) StopService(serviceName string, enclaveName string) error {

	enclaveContext, err := dc.GetEnclaveContext(enclaveName)
	if err != nil {
		return WrapMessageToError(err, "Failed To Stop Services")
	}

	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
	starlarkConfig := GetStarlarkRunConfig(params, "", "")

	_, err = enclaveContext.RunStarlarkScriptBlocking(dc.ctx, stopServiceStarlarkScript, starlarkConfig)
	if err != nil {
		return ErrStarlarkRunFailed
	}
	return nil
}

func (dc *diveContext) StopServices(enclaveName string) error {

	enclaveCtx, err := dc.GetEnclaveContext(enclaveName)

	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Stop Services on enclave %s", enclaveName)
	}

	services, err := enclaveCtx.GetServices()

	if err != nil {
		return WrapMessageToErrorf(ErrKurtosisService, "%s. Failed To Stop Services", err)
	}

	for serviceName := range services {
		params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
		starlarkConfig := GetStarlarkRunConfig(params, "", "")

		_, err = enclaveCtx.RunStarlarkScriptBlocking(dc.ctx, stopServiceStarlarkScript, starlarkConfig)
		if err != nil {
			return ErrStarlarkRunFailed
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
		return WrapMessageToError(ErrKurtosisService, "Failed To Remove Services")
	}

	for serviceName := range services {
		params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
		starlarkConfig := GetStarlarkRunConfig(params, "", "")

		_, err = enclaveCtx.RunStarlarkScriptBlocking(dc.ctx, removeServiceStarlarkScript, starlarkConfig)
		if err != nil {
			return ErrStarlarkRunFailed
		}

	}

	return nil
}

func (dc *diveContext) RemoveService(serviceName string, enclaveName string) error {
	enclaveContext, err := dc.GetEnclaveContext(enclaveName)
	if err != nil {
		return WrapMessageToErrorf(err, "Failed To Remove Service %s in enclave %s", serviceName, enclaveName)
	}

	params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
	starlarkConfig := GetStarlarkRunConfig(params, "", "")

	_, err = enclaveContext.RunStarlarkScriptBlocking(dc.ctx, removeServiceStarlarkScript, starlarkConfig)
	if err != nil {
		return ErrStarlarkRunFailed
	}
	return nil
}

func (dc *diveContext) CreateEnclave(enclaveName string) (*enclaves.EnclaveContext, error) {

	enclaveContext, err := dc.kurtosisContext.CreateEnclave(dc.ctx, enclaveName)

	if err != nil {
		return nil, WrapMessageToError(ErrInvalidKurtosisContext, err.Error())
	}

	return enclaveContext, nil
}

func (dc *diveContext) initKurtosisContext() (*kurtosis_context.KurtosisContext, error) {

	kurtosisContext, err := kurtosis_context.NewKurtosisContextFromLocalEngine()

	if err != nil {
		return nil, WrapMessageToError(ErrInitializingKurtosis, err.Error())

	}

	return kurtosisContext, nil
}
func (dc *diveContext) checkEnclaveExist(enclaveName string) (*EnclaveInfo, error) {

	kurtosisContext, err := dc.GetKurtosisContext()
	if err != nil {
		return nil, WrapMessageToErrorf(err, "Failed To Check Enclave %s", enclaveName)
	}
	enclaveInfo, err := kurtosisContext.GetEnclave(dc.ctx, enclaveName)
	if err != nil {
		return nil, WrapMessageToError(ErrInvalidEnclaveContext, err.Error())
	}

	return &EnclaveInfo{
		Name:      enclaveInfo.Name,
		Uuid:      enclaveInfo.EnclaveUuid,
		ShortUuid: enclaveInfo.ShortenedUuid,
	}, nil
}

func (dc *diveContext) RemoveServicesByServiceNames(services map[string]string, enclaveName string) error {
	enclaveCtx, err := dc.GetEnclaveContext(enclaveName)

	if err != nil {
		return WrapMessageToError(err, "Failed To Remove Services")
	}

	for serviceName := range services {
		params := fmt.Sprintf(`{"service_name": "%s"}`, serviceName)
		starlarkConfig := GetStarlarkRunConfig(params, "", "")

		_, err = enclaveCtx.RunStarlarkScriptBlocking(dc.ctx, removeServiceStarlarkScript, starlarkConfig)
		if err != nil {
			return ErrStarlarkRunFailed
		}

	}

	return nil
}

func (dc *diveContext) Exit(statusCode int) {
	os.Exit(statusCode)
}
