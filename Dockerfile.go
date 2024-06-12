ARG OPENCV_VERSION=4.6.0
FROM gocv/opencv:${OPENCV_VERSION} AS gocv

FROM ovms_capi_ocv_gst:latest

# required for installing dependencies
USER root

RUN apt-get update && apt-get install -y --no-install-recommends libjpeg62 libdc1394-25 gtk+2.0-dev

# needed for gocv bindings and env is custom
ENV CGO_CPPFLAGS="-I/usr/local/include/opencv4"
ENV CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_face -lopencv_videoio -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_objdetect -lopencv_features2d -lopencv_video -lopencv_dnn -lopencv_xfeatures2d -lopencv_calib3d -lopencv_photo"


COPY --from=gocv /usr/local /usr/local

RUN wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go version

RUN mkdir -p /app
COPY /src/go/go-binding /app

WORKDIR /app
RUN go mod tidy
RUN go build -tags customenv capi.go
