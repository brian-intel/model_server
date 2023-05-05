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

type OVMS_InferenceRequest struct {
	OVMS_InferenceRequest *C.OVMS_InferenceRequest
}

func OVMS_InferenceRequestNew(server *OVMS_Server, servableName string, servableVersion int32) (*OVMS_InferenceRequest, error) {
	request := OVMS_InferenceRequest{}
	status := OVMS_Status{}
	cServableName := C.CString(servableName)
	defer C.free(unsafe.Pointer(cServableName))

	cServableVersion := C.uint32_t(servableVersion)

	status.OVMS_Status = C.OVMS_InferenceRequestNew(&request.OVMS_InferenceRequest, server.OVMS_Server, cServableName, cServableVersion)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("failed to create new inference request, status code=%v", code)
	}

	return &request, nil
}

func (request *OVMS_InferenceRequest) OVMS_InferenceRequestAddInput(inputName string, datatype OVMS_DataType, shape []int64, dimCount int32) error {
	status := OVMS_Status{}
	cDataType := C.OVMS_DataType(datatype)
	cDimCount := C.uint32_t(dimCount)

	cInputName := C.CString(inputName)
	defer C.free(unsafe.Pointer(cInputName))

	cSizeArray := make([]C.size_t, len(shape))
	for i, val := range shape {
		cSizeArray[i] = (C.size_t)(val)
	}

	cShape := (*C.size_t)(&cSizeArray[0])

	// cShape := (C.uint64_t)(unsafe.Pointer(&shape[0]))

	// defer C.free(unsafe.Pointer(cShape))

	status.OVMS_Status = C.OVMS_InferenceRequestAddInput(request.OVMS_InferenceRequest, cInputName, cDataType, cShape, cDimCount)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to add input to ovms request, status code=%v", code)
	}
	return nil
}

func (request *OVMS_InferenceRequest) OVMS_InferenceRequestInputSetData(inputName string, data []byte, bufferSize int, bufferType OVMS_BufferType, deviceId int32) error {
	status := OVMS_Status{}

	cInputName := C.CString(inputName)
	defer C.free(unsafe.Pointer(cInputName))

	cBuffersize := C.size_t(bufferSize)

	cDeviceId := C.uint32_t(deviceId)

	ptrData := unsafe.Pointer(&data[0])
	cBufferType := C.OVMS_BufferType(bufferType)
	status.OVMS_Status = C.OVMS_InferenceRequestInputSetData(request.OVMS_InferenceRequest, cInputName, ptrData, cBuffersize, cBufferType, cDeviceId)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set input data into ovms request, status code=%v", code)
	}

	return nil
}