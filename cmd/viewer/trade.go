// +build trade

package main

import (
	"encoding/json"
	"fmt"

	"dev.33.cn/33/btrade/msq"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	LastBlockKey = "lastblock"
)

// for trade.db
func LoadAppInfo(ldb dbm.DB) {
	bytez := ldb.Get([]byte(LastBlockKey))

	lastBlock := &msq.LastBlock{}
	err := msq.UnmarshalMessage(bytez, lastBlock)
	if err != nil {
		panic(err)
	}

	res, _ := json.Marshal(lastBlock)
	fmt.Println(string(res))
}
