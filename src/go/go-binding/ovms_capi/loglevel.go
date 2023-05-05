package ovms_capi

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"


type OVMS_LogLevel struct {
	OVMS_LogLevel C.OVMS_LogLevel
}

func OVMS_LogLevelNew(logLevel LogLevel) OVMS_LogLevel {
	logLevelInt := logLevel.LogLevelInt()
	cLogLevel := OVMS_LogLevel{}
	cLogLevel.OVMS_LogLevel = C.OVMS_LogLevel(logLevelInt)
	return cLogLevel
}