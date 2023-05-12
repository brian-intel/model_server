#!/bin/bash -e

docker build . -f Dockerfile.go -t cgobinding_capi_ms:latest
