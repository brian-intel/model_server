#!/bin/bash -e

docker run --rm -it --net host --entrypoint /bin/bash -v $(pwd)/config:/app/config --privileged cgobinding_capi_ms:latest
