package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	wire "github.com/tendermint/go-wire"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

var appDir = flag.String("app", os.ExpandEnv("$HOME/.tendermint")+"/data/trade.db", "temdermint app db")
var height = flag.Int("h", 0, "app height")
var hash = flag.String("hash", "", "app hash")
var appDb dbm.DB

var lastBlockKey = []byte("lastblock")

type LastBlockInfo struct {
	Height  uint64
	AppHash []byte
}

func main() {
	if len(os.Args) <= 6 {
		fmt.Printf("Usage: %s -app db -h height -hash appHash\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	appDb = dbm.NewDB(FileNameNoExt(*appDir), "leveldb", filepath.Dir(*appDir))
	defer appDb.Close()

	appHash, _ := hex.DecodeString(*hash)

	lastBlock := LoadLastBlock(appDb)
	bytes, _ := json.Marshal(lastBlock)
	fmt.Println("old:", string(bytes))

	lastBlock.Height = uint64(*height)
	lastBlock.AppHash = appHash
	SaveLastBlock(appDb, lastBlock)

	bytes, _ = json.Marshal(lastBlock)
	fmt.Println("new:", string(bytes))
}

// Get the last block from the db
func LoadLastBlock(db dbm.DB) (lastBlock LastBlockInfo) {
	buf := db.Get(lastBlockKey)
	if len(buf) != 0 {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		wire.ReadBinaryPtr(&lastBlock, r, 0, n, err)
		if *err != nil {
			cmn.PanicCrisis(errors.Wrap(*err, "cannot load last block (data has been corrupted or its spec has changed)"))
		}
		// TODO: ensure that buf is completely read.
	}

	return lastBlock
}

func SaveLastBlock(db dbm.DB, lastBlock LastBlockInfo) {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(lastBlock, buf, n, err)
	if *err != nil {
		// TODO
		cmn.PanicCrisis(errors.Wrap(*err, "cannot save last block"))
	}
	db.Set(lastBlockKey, buf.Bytes())
}

func FileNameNoExt(fpath string) string {
	base := filepath.Base(fpath)
	return strings.TrimSuffix(base, filepath.Ext(fpath))
}
