#!/bin/bash

ROOT_PATH=$(cd $(dirname $0) && pwd)

TM_VIEWER=tm_viewer
TM_MIGRATOR=tm_migrator
TM_GENESIS=tm_genesis
TM_VALIDATOR=tm_validator

TM=$HOME/workspace/golang/src/github.com/tendermint/tendermint
mv ${TM}/vendor ${TM}/src
export GOPATH=${TM}:${GOPATH} 

cd "$ROOT_PATH"
go build -v -o "build/$TM_VIEWER" cmd/viewer/*.go
go build -v -o "build/state" cmd/state/*.go
go build -v -o "build/$TM_MIGRATOR" cmd/migrator/*.go
go build -v -o "build/$TM_GENESIS" cmd/genesis/*.go
go build -v -o "build/$TM_VALIDATOR" cmd/validator/*.go

mv ${TM}/src ${TM}/vendor 
