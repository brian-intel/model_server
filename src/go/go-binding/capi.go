package main

// #cgo CPPFLAGS: -I"../.."
// #include <ovms.h>
import "C"

type OVMS_Status struct {
	OVMS_Status C.OVMS_Status
}

type OVMS_Server struct {
	OVMS_Server C.OVMS_Server
}

type OVMS_InferenceRequest struct {
	OVMS_InferenceRequest C.OVMS_InferenceRequest
}

type OVMS_InferenceResponse struct {
	OVMS_InferenceResponse C.OVMS_InferenceResponse
}

func OVMS_ServerStartFromConfigurationFile() {
	return
}

func OVMS_Inference(server *OVMS_Server, inferenceRequest *OVMS_InferenceRequest, inferenceResponse **OVMS_InferenceResponse) interface{} {
	// serverPointer := (C.OVMS_Server)(unsafe.Pointer(server.OVMS_Server))

	// inferenceRsponsePointer := unsafe.Pointer(&inferenceResponse.OVMS_InferenceResponse)
	// defer C.free(unsafe.Pointer(inferenceResponse.OVMS_InferenceResponse))

	return C.OVMS_Inference(&server.OVMS_Server, &inferenceRequest.OVMS_InferenceRequest, &((*inferenceResponse).OVMS_InferenceResponse))
}
