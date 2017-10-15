package convert

import (
	dbm "github.com/tendermint/tmlibs/db"
)

var oBlockDb, oStateDb dbm.DB
var nBlockDb, nStateDb dbm.DB
var totalHeight int

func OnStart(oTmRoot, nTmRoot string) {
	oBlockDb = dbm.NewDB("blockstore", "leveldb", oTmRoot+"/data")
	oStateDb = dbm.NewDB("state", "leveldb", oTmRoot+"/data")
	nBlockDb = dbm.NewDB("blockstore", "leveldb", nTmRoot+"/data")
	nStateDb = dbm.NewDB("state", "leveldb", nTmRoot+"/data")
}

func OnStop() {
	oBlockDb.Close()
	oStateDb.Close()
	nBlockDb.Close()
	nStateDb.Close()
}
