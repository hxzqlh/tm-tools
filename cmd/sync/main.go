package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hxzqlh/tm-tools/util"
	wire "github.com/tendermint/go-wire"
	dbm "github.com/tendermint/tmlibs/db"
)

// TODO: not needed
// sync tendemrint State(LastBlockHeight & AppHash) with app's LastBlockInfo

var stateDir = flag.String("tm", os.ExpandEnv("$HOME/.tendermint")+"/data/state.db", "tendermint state db")
var appDir = flag.String("app", os.ExpandEnv("$HOME/.tendermint")+"/data/trade.db", "temdermint app db")
var stateDb, appDb dbm.DB

var (
	lastBlockKey = []byte("lastblock")
)

type LastBlockInfo struct {
	Height  uint64
	AppHash []byte
}

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -tm db -app db\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	stateDb = dbm.NewDB(FileNameNoExt(*stateDir), "leveldb", filepath.Dir(*stateDir))
	appDb = dbm.NewDB(FileNameNoExt(*appDir), "leveldb", filepath.Dir(*appDir))
	defer stateDb.Close()
	defer appDb.Close()

	lastBlock := LoadLastBlock(appDb)
	fmt.Printf("app height=%v hash=%X\n", lastBlock.Height, lastBlock.AppHash)

	s := util.LoadNewState(stateDb)
	fmt.Printf("state height=%v\n", s.LastBlockHeight)

	// app ---> state
	s.LastBlockHeight = int(lastBlock.Height)
	s.AppHash = lastBlock.AppHash

	util.SaveNewState(stateDb, s)
}

// Get the last block from the db
func LoadLastBlock(db dbm.DB) (lastBlock LastBlockInfo) {
	buf := db.Get(lastBlockKey)
	if len(buf) != 0 {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		wire.ReadBinaryPtr(&lastBlock, r, 0, n, err)
		if *err != nil {
			panic("cannot load last block (data has been corrupted or its spec has changed)")
		}
		// TODO: ensure that buf is completely read.
	}

	return lastBlock
}

func SaveLastBlock(db dbm.DB, lastBlock LastBlockInfo) {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		panic("cannot save last block")
	}
	db.Set(lastBlockKey, buf.Bytes())
}

func FileNameNoExt(fpath string) string {
	base := filepath.Base(fpath)
	return strings.TrimSuffix(base, filepath.Ext(fpath))
}
