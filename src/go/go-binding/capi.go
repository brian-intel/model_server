package main

// #include <stdlib.h>
// #cgo CFLAGS: -Wall -I"/ovms/lib"
// #cgo LDFLAGS: -L"/ovms/lib" -lovms_shared
// #include </ovms/include/ovms.h>
import "C"
import (
	"fmt"
	"image"
	"log"
	"net/http"
	"reflect"
	"time"
	"unsafe"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

type LogLevel string
type OVMS_DataType int
type OVMS_BufferType int

const (
	OVMS_LOG_TRACE   LogLevel = "TRACE"
	OVMS_LOG_DEBUG            = "DEBUG"
	OVMS_LOG_INFO             = "INFO"
	OVMS_LOG_WARNING          = "WARNING"
	OVMS_LOG_ERROR            = "ERROR"
)

const (
	OVMS_DATATYPE_BF16 OVMS_DataType = iota
	OVMS_DATATYPE_FP64
	OVMS_DATATYPE_FP32
	OVMS_DATATYPE_FP16
	OVMS_DATATYPE_I64
	OVMS_DATATYPE_I32
	OVMS_DATATYPE_I16
	OVMS_DATATYPE_I8
	OVMS_DATATYPE_I4
	OVMS_DATATYPE_U64
	OVMS_DATATYPE_U32
	OVMS_DATATYPE_U16
	OVMS_DATATYPE_U8
	OVMS_DATATYPE_U4
	OVMS_DATATYPE_U1
	OVMS_DATATYPE_BOOL
	OVMS_DATATYPE_CUSTOM
	OVMS_DATATYPE_UNDEFINED
	OVMS_DATATYPE_DYNAMIC
	OVMS_DATATYPE_MIXED
	OVMS_DATATYPE_Q78
	OVMS_DATATYPE_BIN
	OVMS_DATATYPE_END
)

const (
	OVMS_BUFFERTYPE_CPU OVMS_BufferType = iota
	OVMS_BUFFERTYPE_CPU_PINNED
	OVMS_BUFFERTYPE_GPU
	OVMS_BUFFERTYPE_HDDL
)

func (level LogLevel) LogLevelInt() int {
	switch level {
	case OVMS_LOG_TRACE:
		return 0
	case OVMS_LOG_DEBUG:
		return 1
	case OVMS_LOG_INFO:
		return 2
	case OVMS_LOG_WARNING:
		return 3
	case OVMS_LOG_ERROR:
		return 4
	default:
		return 224
	}
}

type TensorOutput struct {
	Name       string
	DataType   OVMS_DataType
	Shape      []int
	DimCount   int
	BufferType OVMS_BufferType
	DeviceId   int
	Data       []byte
	ByteSize   int
}

type OVMS_Status struct {
	OVMS_Status *C.OVMS_Status
}

type OVMS_ServerSettings struct {
	OVMS_ServerSettings *C.OVMS_ServerSettings
	GrpcPort            int32
	RestPort            int32
}

type OVMS_ModelsSettings struct {
	OVMS_ModelsSettings *C.OVMS_ModelsSettings
}

type OVMS_Server struct {
	OVMS_Server         *C.OVMS_Server
	OVMS_ServerSettings *OVMS_ServerSettings
	OVMS_ModelsSettings *OVMS_ModelsSettings
}

type OVMS_LogLevel struct {
	OVMS_LogLevel C.OVMS_LogLevel
}

type OVMS_InferenceRequest struct {
	OVMS_InferenceRequest *C.OVMS_InferenceRequest
}

type OVMS_InferenceResponse struct {
	OVMS_InferenceResponse *C.OVMS_InferenceResponse
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

// const char* inputName, OVMS_DataType datatype, const uint64_t* shape, uint32_t dimCount
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
	// defer C.free(unsafe.Pointer(&cBuffersize))

	cDeviceId := C.uint32_t(deviceId)
	// defer C.free(unsafe.Pointer(&cDeviceId))

	ptrData := unsafe.Pointer(&data[0])
	cBufferType := C.OVMS_BufferType(bufferType)
	status.OVMS_Status = C.OVMS_InferenceRequestInputSetData(request.OVMS_InferenceRequest, cInputName, ptrData, cBuffersize, cBufferType, cDeviceId)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		return fmt.Errorf("failed to set input data into ovms request, status code=%v", code)
	}

	return nil
}

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

// (OVMS_InferenceResponse* res, uint32_t* count) {
func (response *OVMS_InferenceResponse) OVMS_InferenceResponseGetOutputCount() (int, error) {
	cCount := C.uint32_t(0)
	// defer C.free(unsafe.Pointer(&cCount))

	status := OVMS_Status{}
	status.OVMS_Status = C.OVMS_InferenceResponseGetOutputCount(response.OVMS_InferenceResponse, &cCount)
	if status.OVMS_Status != nil {
		code, _ := status.OVMS_StatusGetCode()
		details, _ := status.OVMS_StatusGetDetails()
		return 0, fmt.Errorf("failed to get output count from inference response, status code=%v, details=%v", code,details)
	}
	return (int)(cCount), nil

}
func Carray2slice(array *C.uint64_t, len int) []C.int {
	var list []C.int
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return list
}

func (response *OVMS_InferenceResponse) OVMS_InferenceResponseGetOutput(id int32, deviceId int) (*TensorOutput, error) {
	var cName  = C.CString("")
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
	outShape := []int{}
	shapeSlice := unsafe.Slice(cShape, goDimCount)
	for _, val := range shapeSlice {
		outShape = append(outShape, int(val))
	}
	tensorOutput := TensorOutput{
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

func (status *OVMS_Status) OVMS_StatusGetCode() (int32, error) {
	cINT32 := C.uint32_t(0)
	// defer C.free(unsafe.Pointer(&cINT32))

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

func OVMS_LogLevelNew(logLevel LogLevel) OVMS_LogLevel {
	logLevelInt := logLevel.LogLevelInt()
	cLogLevel := OVMS_LogLevel{}
	cLogLevel.OVMS_LogLevel = C.OVMS_LogLevel(logLevelInt)
	return cLogLevel
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

func main() {
	server_grpc_port := 9178
	server_http_port := 11338
	// skip_server_create := 0;

	serverSettings, err := OVMS_ServerSettingsNew()
	if err != nil {
		fmt.Printf("Error:%v", err)
	}

	modelSettings, err := OVMS_ModelsSettingsNew()
	if err != nil {
		fmt.Printf("Error:%v", err)
	}

	server, err := OVMS_ServerNew()
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	defer func() {
		server.OVMS_ServerDelete()
		serverSettings.OVMS_ServerSettingsDelete()
		modelSettings.OVMS_ModelsSettingsDelete()
	}()

	serverSettings.OVMS_ServerSettingsSetGrpcPort(int32(server_grpc_port))
	serverSettings.OVMS_ServerSettingsSetRestPort(int32(server_http_port))

	logLevel := OVMS_LogLevelNew(OVMS_LOG_DEBUG)
	serverSettings.OVMS_ServerSettingsSetLogLevel(logLevel)

	// config_yolov5.json
	modelSettings.OVMS_ModelsSettingsSetConfigPath("/ovms/demos/config_yolov5.json")

	server.OVMS_ServerSettings = serverSettings
	server.OVMS_ModelsSettings = modelSettings

	err = server.OVMS_ServerStartFromConfigurationFile()
	if err != nil {
		fmt.Printf("Error:%v", err)
		server.OVMS_ServerDelete()

	}
	inputSrc := "/ovms/demos/coca-cola-4465029.mp4"
	webcam, err := gocv.OpenVideoCapture(inputSrc) //  /dev/video4
	if err != nil {
		errMsg := fmt.Errorf("faile to open device: %s", inputSrc)
		fmt.Println(errMsg)
	}
	defer webcam.Close()

	img := gocv.NewMat()
	defer img.Close()

	// create the mjpeg stream
	stream := mjpeg.NewStream()

	go run(webcam, &img, stream, server)

	// start http server
	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}

func run(webcam *gocv.VideoCapture, img *gocv.Mat, stream *mjpeg.Stream,
		server *OVMS_Server) {
	for webcam.IsOpened() {
		if ok := webcam.Read(img); !ok {
			// retry once after 1 millisecond
			time.Sleep(time.Millisecond)
			continue
		}
		if img.Empty() {
			continue
		}

		gocv.Resize(*img, img, image.Point{416, 416}, 0, 0, 3)

		// convert to image matrix to use float32
		img.ConvertTo(img, gocv.MatTypeCV32F)
		imgToBytes := img.ToBytes()

		// start := time.Now().UnixMilli()

		request, err := OVMS_InferenceRequestNew(server,"yolov5",0)
		if err != nil {
			fmt.Printf("failed to create new inference request: %v\n",err)
		}

		err = request.OVMS_InferenceRequestAddInput("images",OVMS_DATATYPE_FP32,[]int64{1,416,416,3},4)
		if err != nil {
			fmt.Printf("failed to add input to inference request: %v\n", err)
		}

		dataSize := img.Step() * img.Rows()

		err = request.OVMS_InferenceRequestInputSetData("images",imgToBytes,dataSize,OVMS_BUFFERTYPE_CPU,0)
		if err != nil {
			fmt.Printf("failed to set input data to inference request: %v\n", err)
		}

		response, err := OVMS_Inference(server,request)
		if err != nil {
			fmt.Printf("failed to get response for inference request: %v\n", err)
		}


		outputCount, err := response.OVMS_InferenceResponseGetOutputCount()
		if err != nil {
			fmt.Printf("failed to get output count from inference response: %v\n", err)
		}
		fmt.Printf(">>>>Output Count: %v\n", outputCount)


		// collect all tensor outputs
		TensorOutputs := []*TensorOutput{}
		for id:= (int32)(outputCount-1);id>=0;id--{
			tensorOutput, err := response.OVMS_InferenceResponseGetOutput(id,42)
			if err != nil {
				fmt.Printf("failed to get output from inference response: %v\n", err)
			}
			TensorOutputs = append(TensorOutputs, tensorOutput)
		}
		

		buf, _ := gocv.IMEncode(".jpg", *img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}

}
