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

type OVMS_InferenceResponse struct {
	OVMS_InferenceResponse *C.OVMS_InferenceResponse
}

func (response *OVMS_InferenceResponse) OVMS_InferenceResponseGetOutputCount() (int, error) {
	cCount := C.uint32_t(0)
	// defer C.free(unsafe.Pointer(&cCount))

	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_InferenceResponseGetOutputCount(response.OVMS_InferenceResponse, &cCount)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		details, _ := status.OVMS_StatusGetDetails()
		return 0, fmt.Errorf("failed to get output count from inference response, status code=%v, details=%v", code, details)
	}
	return (int)(cCount), nil

}

func (response *OVMS_InferenceResponse) OVMS_InferenceResponseGetOutput(id int32, deviceId int) (*OVMSOutput, error) {
	var cName = C.CString("")
	defer C.free(unsafe.Pointer(cName))

	var dataType C.OVMS_DataType

	var cShape *C.size_t
	defer C.free(unsafe.Pointer(cShape))

	var dimCount C.uint32_t

	var data unsafe.Pointer

	var bufferType C.OVMS_BufferType

	var byteSize C.size_t

	cId := C.uint32_t(id)

	cDeviceId := C.uint32_t(deviceId)

	status := OVMS_Status{}

	status.OVMS_Status = C.OVMS_InferenceResponseGetOutput(response.OVMS_InferenceResponse,
		cId, &cName, &dataType, &cShape, &dimCount, &data, &byteSize, &bufferType, &cDeviceId)

	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("failed to get output for id=%v, status code=%v", id, code)
	}

	goDimCount := int(dimCount)
	outShape := []int64{}
	shapeSlice := unsafe.Slice(cShape, goDimCount)
	for _, val := range shapeSlice {
		outShape = append(outShape, int64(val))
	}
	tensorOutput := OVMSOutput{
		DeviceId:   int(cDeviceId),
		Name:       C.GoString(cName),
		DataType:   (OVMS_DataType)(dataType),
		Shape:      outShape,
		DimCount:   goDimCount,
		BufferType: (OVMS_BufferType)(bufferType),
		Data:       C.GoBytes(data, (C.int)(byteSize)),
		ByteSize:   int(byteSize),
	}

	return &tensorOutput, nil
}