#!/usr/bin/env bash

DIRS=$(go list ./... | grep -v /vendor/)

go fmt $DIRS
go vet $DIRS
TRANSLATE_DEBUG=1 go test -race $DIRS