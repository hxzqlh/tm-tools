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
var stateDb dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -state db -hash appHash\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	stateDb = dbm.NewDB(util.FileNameNoExt(*stateDir), "leveldb", filepath.Dir(*stateDir))
	defer stateDb.Close()

	bytes, _ := hex.DecodeString(*hash)
	s := util.LoadNewState(stateDb)
	s.AppHash = bytes
	util.SaveNewState(stateDb, s)
}
