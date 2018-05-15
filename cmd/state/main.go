package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hxzqlh/tm-tools/util"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	stateKey = "stateKey"
)

var stateDir = flag.String("state", os.ExpandEnv("$HOME/.tendermint")+"/data/state.db", "tendermint state db")
var hash = flag.String("hash", "", "tendermint app hash")
var height = flag.Int("h", 0, "last block height")

var stateDb dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -state db -h height -hash appHash\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	stateDb = dbm.NewDB(util.FileNameNoExt(*stateDir), "leveldb", filepath.Dir(*stateDir))
	defer stateDb.Close()

	bytes, _ := hex.DecodeString(*hash)

	s := util.LoadNewState(stateDb)
	fmt.Printf("old height=%v hash=%X\n", s.LastBlockHeight, s.AppHash)
	s.AppHash = bytes
	s.LastBlockHeight = *height
	fmt.Printf("new height=%v hash=%X\n", s.LastBlockHeight, s.AppHash)

	util.SaveNewState(stateDb, s)
}
