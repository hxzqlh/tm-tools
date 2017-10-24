#!/bin/bash
set -x

VERSION=1.0.0.0
ROOT_PATH=$(cd $(dirname $0) && pwd)
PKG=tm_tools

cd "$ROOT_PATH"
bash build.sh

mkdir -p $PKG
cp run.sh build/* $PKG

tar -zvcf ${PKG}.tar.gz $PKG
rm -rf $PKG
