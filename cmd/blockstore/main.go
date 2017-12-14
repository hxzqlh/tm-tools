package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/blockchain"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	blockStoreKey = "blockStore"
)

var dbDir = flag.String("db", os.ExpandEnv("$HOME/.tendermint")+"/data/blockstore.db", "tendermint blockstore db")
var height = flag.Int("h", 0, "height")
var blockStoreDb dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -db blockstore.db -h height\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	blockStoreDb = dbm.NewDB(util.FileNameNoExt(*dbDir), "leveldb", filepath.Dir(*dbDir))
	defer blockStoreDb.Close()

	store := blockchain.LoadBlockStoreStateJSON(blockStoreDb)
	log.Println("old store:", store)
	store.Height = *height
	log.Println("new store:", store)
	store.Save(blockStoreDb)
}
