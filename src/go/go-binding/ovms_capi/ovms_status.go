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

type OVMS_Status struct {
	OVMS_Status *C.OVMS_Status
}

func (status *OVMS_Status) OVMS_StatusGetCode() (int32, error) {
	cINT32 := C.uint32_t(0)

	rtnStatus := OVMS_Status{}
	rtnStatus.OVMS_Status = C.OVMS_StatusGetCode(status.OVMS_Status, &cINT32)
	if rtnStatus.OVMS_Status != nil {
		code, _ := rtnStatus.OVMS_StatusGetCode()
		return code, fmt.Errorf("failed to get status code, Status code=%v", code)
	}

	return int32(cINT32), nil
}

func (status *OVMS_Status) OVMS_StatusGetDetails() (string, error) {
	cStr := C.CString("")
	defer C.free(unsafe.Pointer(cStr))

	status.OVMS_Status = C.OVMS_StatusGetDetails(status.OVMS_Status, &cStr)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return "", fmt.Errorf("failed to get status details with Statcode=%v", code)
	}
	details := C.GoString(cStr)

	return details, nil
}

func (status *OVMS_Status) OVMS_StatusDelete() {
	C.OVMS_StatusDelete(status.OVMS_Status)
	status.OVMS_Status = nil
}
