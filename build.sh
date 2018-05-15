#!/bin/bash

ROOT_PATH=$(cd $(dirname $0) && pwd)

TM=$GOPATH/src/github.com/tendermint/tendermint
mv ${TM}/vendor ${TM}/src
export GOPATH=${TM}:${GOPATH} 

cd "$ROOT_PATH"

go build -v -o "build/tm_sync" cmd/sync/*.go
go build -v -o "build/tm_viewer" cmd/viewer/*.go
go build -v -o "build/tm_migrator" cmd/migrator/*.go
go build -v -o "build/tm_genesis" cmd/genesis/*.go
go build -v -o "build/tm_validator" cmd/validator/*.go

mv ${TM}/src ${TM}/vendor 
