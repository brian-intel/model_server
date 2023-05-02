#!/bin/bash -e

docker build . -f Dockerfile.go -t cgobinding:dev
