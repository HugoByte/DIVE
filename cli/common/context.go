package common

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
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

// The `GetKurtosisContext` function is a method of the `diveContext` struct.
// Used to Get the kurtosis context from the kurtosis engine
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

// The `GetEnclaves` function is a method of the `diveContext` struct.
// Used to Get List Of Enclaves Running / Stopped
func (dc *diveContext) GetEnclaves() ([]EnclaveInfo, error) {
	enclavesInfo, err := dc.kurtosisContext.GetEnclaves(dc.ctx)

	if err != nil {
		return nil, WrapMessageToError(ErrInvalidEnclaveContext, err.Error())
	}
	var enclaves []EnclaveInfo
	enclaveMap := enclavesInfo.GetEnclavesByName()
	for _, enclaveInfoList := range enclaveMap {
		for _, enclaveInfo := range enclaveInfoList {

			enclaves = append(enclaves, EnclaveInfo{
				Name:        enclaveInfo.Name,
				Uuid:        enclaveInfo.EnclaveUuid,
				ShortUuid:   enclaveInfo.ShortenedUuid,
				Status:      strings.Replace(enclaveInfo.ContainersStatus.String(), "EnclaveContainersStatus_", "", -1),
				CreatedTime: enclaveInfo.CreationTime.AsTime().String()},
			)
		}
	}
	return enclaves, nil
}

// The `GetEnclaveContext` function is a method of the `diveContext` struct.
// Used to Get the Kurtosis Enclave Context
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

// The `CleanEnclaves` function is a method of the `diveContext` struct.
// Used to Cleans given Running Enclaves
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

// The `CleanEnclaveByName` function is a method of the `diveContext` struct.
// Used to Clean given Enclave
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

// The `CheckSkippedInstructions` function is a method of the `diveContext` struct.
// Used to Check the Skipped Instructions
func (dc *diveContext) CheckSkippedInstructions(instructions map[string]bool) bool {

	return len(instructions) != 0
}

// The `StopService` function is a method of the `diveContext` struct. It is used to stop given service name from a specific enclave.
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

// The `StopServices` function is a method of the `diveContext` struct. It is used to stop all
// services from a specific enclave.
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

// The `RemoveServices` function is a method of the `diveContext` struct. It is used to remove all
// services from a specific enclave.
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

// The `RemoveService` function is a method of the `diveContext` struct. It is used to remove a
// specific service from an enclave.
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

// The `CreateEnclave` function is used to create a new enclave with the specified name. It takes in
// the name of the enclave as a parameter and returns an `EnclaveContext` object and an error.
func (dc *diveContext) CreateEnclave(enclaveName string) (*enclaves.EnclaveContext, error) {

	enclaveContext, err := dc.kurtosisContext.CreateEnclave(dc.ctx, enclaveName)

	if err != nil {
		return nil, WrapMessageToError(ErrInvalidKurtosisContext, err.Error())
	}

	return enclaveContext, nil
}

// The `initKurtosisContext` function is responsible for initializing the Kurtosis context by creating
// a new instance of `KurtosisContext` using the `NewKurtosisContextFromLocalEngine` function from the
// `kurtosis_context` package. If there is an error during initialization, it returns an error.
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

// The `RemoveServicesByServiceNames` function is used to remove multiple services from an enclave. It
// takes in a map of service names and their corresponding IDs, as well as the name of the enclave from
// which the services should be removed.
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

// The `Exit` function is used to terminate the program with a specified exit status code. It calls the
// `os.Exit` function, which immediately terminates the program with the given status code.
func (dc *diveContext) Exit(statusCode int) {
	os.Exit(statusCode)
}
