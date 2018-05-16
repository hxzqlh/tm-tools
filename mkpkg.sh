#!/bin/bash
set -x

VERSION=1.1.0.0
ROOT_PATH=$(cd $(dirname $0) && pwd)
PKG=tm_tools_v${VERSION}

cd "$ROOT_PATH"
bash build.sh

mkdir -p $PKG
cp run.sh build/* $PKG

tar -zvcf ${PKG}.tar.gz $PKG
rm -rf $PKG
