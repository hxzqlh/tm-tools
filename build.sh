#!/bin/bash

ROOT_PATH=$(cd $(dirname $0) && pwd)

TM=$GOPATH/src/github.com/tendermint/tendermint
mv ${TM}/vendor ${TM}/src
export GOPATH=${TM}:${GOPATH} 

cd "$ROOT_PATH"

go build -v -o "build/tm_viewer" cmd/viewer/*.go
go build -v -o "build/tm_migrator" cmd/migrator/*.go
go build -v -o "build/tm_genesis" cmd/genesis/*.go
go build -v -o "build/tm_validator" cmd/validator/*.go
go build -v -o "build/set_store" cmd/blockstore/*.go
go build -v -o "build/set_app" cmd/app/*.go
go build -v -o "build/set_state" cmd/state/*.go

mv ${TM}/src ${TM}/vendor 
