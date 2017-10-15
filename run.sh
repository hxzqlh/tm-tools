#!/bin/bash
set -x

[ $# -lt 3 ] && echo "$0 old_tm new_tm priv_dir" && exit 1

ROOT_PATH=$(cd $(dirname $0) && pwd)

OLD_TM=$(cd $(dirname $1) && pwd)
NEW_TM=$(cd $(dirname $2) && pwd)
PRIV_DIR=$(cd $(dirname $3) && pwd)

TM=/home/hxz/gop/bin/tendermint_0.10.0

$ROOT_PATH/tm_migrator -old $OLD_TM -new $NEW_TM -priv $PRIV_DIR
TMROOT="$NEW_TM" $TM node
