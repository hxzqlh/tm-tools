#!/bin/bash
set -x

VERSION=1.0.0.0
ROOT_PATH=$(cd $(dirname $0) && pwd)

TM_VIEWER=tm_viewer
TM_MIGRATOR=tm_migrator
PKG=tm_tools

TM=$HOME/workspace/golang/src/github.com/tendermint/tendermint
mv ${TM}/vendor ${TM}/src
export GOPATH=${TM}:${GOPATH} 

cd "$ROOT_PATH"
go build -o "build/$TM_VIEWER" cmd/viewer/*.go
go build -o "build/$TM_MIGRATOR" cmd/migrator/*.go
mv ${TM}/src ${TM}/vendor 

mkdir -p $PKG
cp run.sh build/* $PKG

tar -zvcf ${PKG}.tar.gz $PKG
rm -rf $PKG
