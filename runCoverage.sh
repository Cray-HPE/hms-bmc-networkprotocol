#!/usr/bin/env bash

# Build the build base image
docker build -t cray/hms-bmc-networkprotocol-build-base -f Dockerfile.build-base .

docker build -t cray/hms-bmc-networkprotocol-coverage -f Dockerfile.coverage .
