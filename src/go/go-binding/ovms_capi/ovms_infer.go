package ovms_capi

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"
import (
	"fmt"
)

func OVMS_Inference(server *OVMS_Server, inferenceRequest *OVMS_InferenceRequest) (*OVMS_InferenceResponse, error) {
	status := OVMS_Status{}
	inferenceResponse := OVMS_InferenceResponse{}

	status.OVMS_Status = C.OVMS_Inference(server.OVMS_Server, inferenceRequest.OVMS_InferenceRequest, &inferenceResponse.OVMS_InferenceResponse)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return nil, fmt.Errorf("inference failed, status code=%v", code)
	}

	return &inferenceResponse, nil
}