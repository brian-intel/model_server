FROM ovms_capi_ocv_gst:latest

RUN wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN go version

RUN mkdir -p /app

WORKDIR /app

COPY /src /app