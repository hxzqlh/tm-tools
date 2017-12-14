package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dev.33.cn/33/btrade/msq"
	"github.com/hxzqlh/tm-tools/util"
	abci "github.com/tendermint/abci/types"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var dbDir = flag.String("db", os.ExpandEnv("$HOME/.tendermint")+"/data/tx_index.db", "tendermint tx index db")
var txHash = flag.String("k", "", "tx hash")
var hexStr = flag.String("v", "", "value hex string")
var db dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -db tx_index.db -k hash [-v hexStr]\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	var action = "get"
	if len(os.Args) > 5 {
		if *txHash == "" || *hexStr == "" {
			fmt.Println("set but no tx hash/val specified")
			os.Exit(0)
		} else {
			action = "set"
		}
	}

	db = dbm.NewDB(util.FileNameNoExt(*dbDir), "leveldb", filepath.Dir(*dbDir))
	defer db.Close()

	if action == "get" {
		get()
	} else {
		set()
	}
}

func set() {
	key, _ := hex.DecodeString(*txHash)
	res, err := LoadTxResult(key)
	if err != nil {
		panic(err)
	}
	if res == nil {
		panic("TxResult is nil")
	}

	bytes, _ := json.Marshal(res.Tx)
	fmt.Printf("old result, Height: %v, index: %v, Tx: %s, Result: [Code: %v, Data: %X, Log: %v]\n",
		res.Height, res.Index, string(bytes), res.Result.Code, res.Result.Data, res.Result.Log)

	val, _ := hex.DecodeString(*hexStr)
	res.Result = abci.ResponseDeliverTx{abci.CodeType_OK, val, ""}
	SaveTxResult(res)

	fmt.Printf("new result, Height: %v, index: %v, Tx: %s, Result: [Code: %v, Data: %X, Log: %v]\n",
		res.Height, res.Index, string(bytes), res.Result.Code, res.Result.Data, res.Result.Log)
}

func get() {
	key, _ := hex.DecodeString(*txHash)
	res, err := LoadTxResult(key)
	if err != nil {
		panic(err)
	}
	if res == nil {
		panic("TxResult is nil")
	}

	bytes, _ := json.Marshal(res.Tx)
	fmt.Printf("Height: %v, index: %v, Tx: %s, Result: [Code: %v, Data: %X, Log: %v]\n",
		res.Height, res.Index, string(bytes), res.Result.Code, res.Result.Data, res.Result.Log)

	var obj msq.WriteRequest
	err = msq.UnmarshalMessage(res.Tx, &obj)
	if err != nil {
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

func LoadTxResult(hash []byte) (*types.TxResult, error) {
	rawBytes := db.Get(hash)
	if rawBytes == nil {
		return nil, nil
	}

	r := bytes.NewReader(rawBytes)
	var n int
	var err error
	txResult := wire.ReadBinary(&types.TxResult{}, r, 0, &n, &err).(*types.TxResult)
	if err != nil {
		return nil, fmt.Errorf("Error reading TxResult: %v", err)
	}

	return txResult, nil
}

func SaveTxResult(result *types.TxResult) {
	rawBytes := wire.BinaryBytes(result)
	db.Set(result.Tx.Hash(), rawBytes)
}
