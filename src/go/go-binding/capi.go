package main

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include <../../ovms.h>
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
	OVMS_InferenceResponse *C.OVMS_InferenceResponse
}

func OVMS_InferenceRequestNew() {
	return
}

func OVMS_InferenceRequestAddInput() {
	return
}

func OVMS_InferenceRequestInputSetData() {
	return
}

func C_API_OVMS_Inference(server *OVMS_Server, inferenceRequest *OVMS_InferenceRequest, inferenceResponse **OVMS_InferenceResponse) interface{} {
	// serverPointer := (C.OVMS_Server)(unsafe.Pointer(server.OVMS_Server))

	// inferenceRsponsePointer := unsafe.Pointer(&inferenceResponse.OVMS_InferenceResponse)
	// defer C.free(unsafe.Pointer(inferenceResponse.OVMS_InferenceResponse))

	return C.OVMS_Inference(&server.OVMS_Server, &inferenceRequest.OVMS_InferenceRequest, &((*inferenceResponse).OVMS_InferenceResponse))
}

func OVMS_StatusGetCode() {
	return
}

func OVMS_StatusGetDetails() {
	return
}

func main() {
	return
}
