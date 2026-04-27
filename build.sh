#!/bin/bash
cd /Users/dsi/projects/penelope
export GOPATH=/Users/dsi/go
export PATH=$PATH:/usr/local/go/bin
go mod tidy
BUILD_DATE=$(date -u +"%Y-%m-%d_%H:%M:%S")
go build -ldflags="-X main.BuildDate=${BUILD_DATE}" -o necromancy main.go
