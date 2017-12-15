package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dev.33.cn/33/btrade/msq"
	"github.com/hxzqlh/tm-tools/util"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tendermint/state/txindex/kv"
	dbm "github.com/tendermint/tmlibs/db"
)

var dbDir = flag.String("db", os.ExpandEnv("$HOME/.tendermint")+"/data/tx_index.db", "tendermint tx index db")
var txHash = flag.String("h", "", "tx hash")
var db dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -db tx_index.db -h hash\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	db = dbm.NewDB(util.FileNameNoExt(*dbDir), "leveldb", filepath.Dir(*dbDir))
	defer db.Close()

	txIndex := kv.NewTxIndex(db)
	key, _ := hex.DecodeString(*txHash)
	res, err := txIndex.Get(key)
	if err != nil {
		panic(err)
	}

	bytes, _ := json.Marshal(res)
	fmt.Printf("%x ---> %s\n", res.Tx.Hash(), string(bytes))

	var obj msq.WriteRequest
	err = msq.UnmarshalMessage(res.Tx, &obj)
	if err != nil {
		fmt.Println(string(bytes))
		panic(err)
	}

	bytes, _ = json.Marshal(obj)
	fmt.Println("tx request:", string(bytes))

	if res.Result.Code == abci.CodeType_OK {
		var resp msq.Response
		err = msq.UnmarshalMessage(res.Result.Data, &resp)
		if err != nil {
			fmt.Println(string(bytes))
			panic(err)
		}

		bytes, _ = json.Marshal(resp)
		fmt.Println("tx response:", string(bytes))
	} else {
		fmt.Println("tx response:", string(res.Result.Data))
	}
}
