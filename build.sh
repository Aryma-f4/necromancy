#!/bin/bash
cd /Users/dsi/projects/penelope
export GOPATH=/Users/dsi/go
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o necromancy main.go
