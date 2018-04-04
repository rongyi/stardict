#!/bin/bash

GOOS=windows GOARCH=amd64 go build -o stardict-windows
GOOS=darwin GOARCH=amd64 go build -o stardict-mac
# default is linux
go build -o stardict-linux
