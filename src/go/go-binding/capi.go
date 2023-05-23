// go:build (darwin && cgo) || linux
//go:build (darwin && cgo) || linux
// +build darwin,cgo linux

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hybridgroup/mjpeg"
	"github.com/openvinotoolkit/model_sever/src/go/go-binding/ovms_capi"
	tensoroutput "github.com/openvinotoolkit/model_sever/src/go/go-binding/tensorOutput"
	"github.com/openvinotoolkit/model_sever/src/go/go-binding/yolov5"
	"gocv.io/x/gocv"
)

const (
	modelName        = "yolov5s"
	modelConfigPath  = "/app/config/config_detection.json"
	inputVideoSource = "/ovms/demos/coca-cola-4465029.mp4"
)

func main() {
	server_grpc_port := 9178
	server_http_port := 11338
	// skip_server_create := 0;

	serverSettings, err := ovms_capi.OVMS_ServerSettingsNew()
	if err != nil {
		fmt.Printf("Error:%v", err)
	}

	modelSettings, err := ovms_capi.OVMS_ModelsSettingsNew()
	if err != nil {
		fmt.Printf("Error:%v", err)
	}

	server, err := ovms_capi.OVMS_ServerNew()
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

	logLevel := ovms_capi.OVMS_LogLevelNew(ovms_capi.OVMS_LOG_DEBUG)
	serverSettings.OVMS_ServerSettingsSetLogLevel(logLevel)

	// config_yolov5.json
	modelSettings.OVMS_ModelsSettingsSetConfigPath(modelConfigPath)

	server.OVMS_ServerSettings = serverSettings
	server.OVMS_ModelsSettings = modelSettings

	err = server.OVMS_ServerStartFromConfigurationFile()
	if err != nil {
		fmt.Printf("Error:%v", err)
		server.OVMS_ServerDelete()

	}
	inputSrc := inputVideoSource
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
	camHeight := float32(webcam.Get(gocv.VideoCaptureFrameHeight))
	camWidth := float32(webcam.Get(gocv.VideoCaptureFrameWidth))
	go run(webcam, &img, stream, server, camWidth, camHeight)

	// start http server
	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))

}

func run(webcam *gocv.VideoCapture, img *gocv.Mat, stream *mjpeg.Stream,
	server *ovms_capi.OVMS_Server, camWidth float32, camHeight float32) {
	var aggregateLatencyAfterInfer float64
	var aggregateLatencyAfterFinalProcess float64
	var frameNum float64

	// output latency metic to txt
	latencyMetricFile := "./results/latency.txt"
	file, err := os.Create(latencyMetricFile)
	if err != nil {
		fmt.Printf("failed to write to file:%v", err)
	}
	defer file.Close()

	for webcam.IsOpened() {
		if ok := webcam.Read(img); !ok {
			// retry once after 1 millisecond
			time.Sleep(time.Millisecond)
			continue
		}
		if img.Empty() {
			continue
		}
		frameNum++
		gocv.Resize(*img, img, image.Point{416, 416}, 0, 0, 3)

		// convert to image matrix to use float32
		img.ConvertTo(img, gocv.MatTypeCV32F)
		imgToBytes := img.ToBytes()

		start := float64(time.Now().UnixMilli())

		request, err := ovms_capi.OVMS_InferenceRequestNew(server, modelName, 0)
		if err != nil {
			fmt.Printf("failed to create new inference request: %v\n", err)
		}

		err = request.OVMS_InferenceRequestAddInput("images", ovms_capi.OVMS_DATATYPE_FP32, []int64{1, 416, 416, 3}, 4)
		if err != nil {
			fmt.Printf("failed to add input to inference request: %v\n", err)
		}

		dataSize := img.Step() * img.Rows()

		err = request.OVMS_InferenceRequestInputSetData("images", imgToBytes, dataSize, ovms_capi.OVMS_BUFFERTYPE_CPU, 0)
		if err != nil {
			fmt.Printf("failed to set input data to inference request: %v\n", err)
		}

		response, err := ovms_capi.OVMS_Inference(server, request)
		if err != nil {
			fmt.Printf("failed to get response for inference request: %v\n", err)
		}

		afterInfer := float64(time.Now().UnixMilli())
		aggregateLatencyAfterInfer += (afterInfer - start)
		latencyStr := fmt.Sprintf("latency after infer: %v\n", aggregateLatencyAfterInfer/frameNum)
		_, err = file.WriteString(latencyStr)
		if err != nil {
			fmt.Printf("failed to write to file:%v", err)
		}

		outputCount, err := response.OVMS_InferenceResponseGetOutputCount()
		if err != nil {
			fmt.Printf("failed to get output count from inference response: %v\n", err)
		}

		// collect all tensor outputs
		tensorOutputs := tensoroutput.TensorOutputs{}
		for id := (int32)(outputCount - 1); id >= 0; id-- {
			ovmsOutput, err := response.OVMS_InferenceResponseGetOutput(id, 42)
			if err != nil {
				fmt.Printf("failed to get output from inference response: %v\n", err)
			}
			tensorOutputs.RawData = append(tensorOutputs.RawData, ovmsOutput.Data)
			tensorOutputs.DataShapes = append(tensorOutputs.DataShapes, ovmsOutput.Shape)
		}

		err = tensorOutputs.ParseRawData()
		if err != nil {
			fmt.Printf("failed to parse raw byte data: %v\n", err)
		}
		detectedObjects := yolov5.DetectedObjects{}
		err = detectedObjects.Postprocess(tensorOutputs, camWidth, camHeight)
		if err != nil {
			fmt.Printf("post process failed: %v\n", err)
		}
		fmt.Printf("length of detectedObjects after  processing: %v\n", len(detectedObjects.Objects))
		detectedObjects = detectedObjects.FinalPostProcessAdvanced()
		// debug
		fmt.Printf("length of detectedObjects after final processing: %v\n", len(detectedObjects.Objects))

		// // track the latency after  processing latency
		afterFinalProcess := float64(time.Now().UnixMilli())
		aggregateLatencyAfterFinalProcess += (afterFinalProcess - start)
		latencyStr = fmt.Sprintf("average latency after final process: %v\n", aggregateLatencyAfterFinalProcess/frameNum)
		_, err = file.WriteString(latencyStr)
		if err != nil {
			fmt.Printf("failed to write to file:%v", err)
		}
		// fmt.Printf("latency after final process: %v\n", afterFinalProcess-start)

		// add bounding boxes to resixed image
		detectedObjects.AddBoxesToFrame(img, color.RGBA{0, 255, 0, 0}, camWidth, camWidth)

		buf, _ := gocv.IMEncode(".jpg", *img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}

}
