package ovms_capi

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type OVMS_ModelsSettings struct {
	OVMS_ModelsSettings *C.OVMS_ModelsSettings
}

func OVMS_ModelsSettingsNew() (*OVMS_ModelsSettings, error) {
	settings := OVMS_ModelsSettings{}
	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_ModelsSettingsNew(&settings.OVMS_ModelsSettings)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("failed to create model settings, status code=%v", code)
	}
	return &settings, nil
}

func (settings *OVMS_ModelsSettings) OVMS_ModelsSettingsDelete() {
	C.OVMS_ModelsSettingsDelete(settings.OVMS_ModelsSettings)
	settings.OVMS_ModelsSettings = nil
}

func (settings *OVMS_ModelsSettings) OVMS_ModelsSettingsSetConfigPath(configPath string) error {
	status := OVMS_Status{}
	cConfigPath := C.CString(configPath)
	defer C.free(unsafe.Pointer(cConfigPath))

	status.OVMS_Status = C.OVMS_ModelsSettingsSetConfigPath(settings.OVMS_ModelsSettings, cConfigPath)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set config path in model settings, status code =%v", code)
	}
	return nil
}