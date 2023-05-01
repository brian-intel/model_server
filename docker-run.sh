#!/bin/bash -e

docker run --rm -it --net host --entrypoint /bin/bash --privileged cgobinding:dev
