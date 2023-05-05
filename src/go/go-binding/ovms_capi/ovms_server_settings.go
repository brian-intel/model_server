package ovms_capi

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"
import (
	"fmt"
)

type OVMS_ServerSettings struct {
	OVMS_ServerSettings *C.OVMS_ServerSettings
	GrpcPort            int32
	RestPort            int32
}


func OVMS_ServerSettingsNew() (*OVMS_ServerSettings, error) {
	settings := OVMS_ServerSettings{}
	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_ServerSettingsNew(&settings.OVMS_ServerSettings)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("failed to create sever settings, status code=%v", code)
	}
	return &settings, nil
}

func (settings *OVMS_ServerSettings) OVMS_ServerSettingsDelete() {
	C.OVMS_ServerSettingsDelete(settings.OVMS_ServerSettings)
	settings.OVMS_ServerSettings = nil
}

func (settings *OVMS_ServerSettings) OVMS_ServerSettingsSetGrpcPort(grpcPort int32) error {
	status := OVMS_Status{}
	cgrpcPort := C.uint(grpcPort)
	status.OVMS_Status = C.OVMS_ServerSettingsSetGrpcPort(settings.OVMS_ServerSettings, cgrpcPort)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set grpc port to server settings, status code=%v", code)
	}
	settings.GrpcPort = grpcPort
	return nil
}

func (settings *OVMS_ServerSettings) OVMS_ServerSettingsSetRestPort(restPort int32) error {
	status := OVMS_Status{}
	crestPort := C.uint(restPort)
	status.OVMS_Status = C.OVMS_ServerSettingsSetRestPort(settings.OVMS_ServerSettings, crestPort)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set rest port to server settings, status code=%v", code)
	}
	settings.RestPort = restPort
	return nil
}

func (settings *OVMS_ServerSettings) OVMS_ServerSettingsSetLogLevel(logLevel OVMS_LogLevel) error {
	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_ServerSettingsSetLogLevel(settings.OVMS_ServerSettings, logLevel.OVMS_LogLevel)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set log level in server settings, status code=%v", code)
	}
	return nil
}

