package ovms_capi

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"
import (
	"fmt"
)

type OVMS_Server struct {
	OVMS_Server         *C.OVMS_Server
	OVMS_ServerSettings *OVMS_ServerSettings
	OVMS_ModelsSettings *OVMS_ModelsSettings
}

func OVMS_ServerNew() (*OVMS_Server, error) {
	server := OVMS_Server{}
	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_ServerNew(&server.OVMS_Server)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("failed to create OVMS server, status code=%v", code)
	}
	return &server, nil
}

func (server *OVMS_Server) OVMS_ServerDelete() {
	C.OVMS_ServerDelete(server.OVMS_Server)
	server.OVMS_Server = nil
}

func (server *OVMS_Server) OVMS_ServerStartFromConfigurationFile() error {
	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_ServerStartFromConfigurationFile(
		server.OVMS_Server,
		server.OVMS_ServerSettings.OVMS_ServerSettings,
		server.OVMS_ModelsSettings.OVMS_ModelsSettings,
	)

	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to start OVMS server from configuration file, status code=%v", code)
	}
	return nil
}