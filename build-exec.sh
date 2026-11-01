#!/usr/bin/env bash

mkdir -p build
CGO_ENABLED=0 go build -o build/agent . 