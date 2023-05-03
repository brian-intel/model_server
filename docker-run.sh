#!/bin/bash -e

docker run --rm -it --net host --entrypoint /bin/bash -v $(pwd)/config:/app/config -v `pwd`/results:/app/results --privileged cgobinding:dev