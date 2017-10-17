package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hxzqlh/tm-tools/util"
	dbm "github.com/tendermint/tmlibs/db"
)

/*
Usage: $ tm_view -db /path/of/db [-a get|getall|block] [-q key] [-d] [-v new|old] [-t height]

// -db : db，Note: the db path cannot end with "/"
// [-a get|getall|block]： read the value of a key | output all keyes | read block info
// [-q key] ：key format
// [-d]: whether decode value，default is "false"
// [-v new|old] ：new(0.10.0), old(0.7.3), default is "new"
// [-t height]: block height，workes with "-a block" arg to read block info at height "N"

examples：
$ tm_view -db /path/of/blockstore.db -a getall
$ tm_view -db /path/of/blockstore.db -a block -t 1 -d
$ tm_view -db /path/of/blockstore.db -q "H:1" -d -v old
$ tm_view -db /path/of/state.db -q "stateKey" -d -v old

| key format | value type | examples |
| ---- |-----| ---- |
| `stateKey` | raw byte of state | |
| `abciResponsesKey` | raw byte of ABCI Responses | |
| `blockStore` | raw json |  "blockStore": {"Height":32} |
| `H:{height}` | raw byte of block meta | H:1 |
| `P:{height}:{index}`| raw byte of block part | P:1:0, P:32:0, P:32:1 |
| `SC:{height}` | raw byte of block seen commit | SC:1, SC:32 |
| `C:{height-1}` | raw byte of block commit | C:0, SC:31 |
*/

var dbpath = flag.String("db", os.ExpandEnv("$HOME/.tendermint")+"/trade.db", "database db")
var action = flag.String("a", "get", "get key from database")
var key = flag.String("q", "", "the query string of database")
var decode = flag.Bool("d", false, "whether decode data")
var limit = flag.Int("l", 0, "limit of query list")
var ver = flag.String("v", "new", "version of tendermint")
var height = flag.Int("t", 1, "block height")
var ldb dbm.DB

func main() {
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s -db /path/of/db [-a get|getall|block] [-q key] [-d] [-v new|old] [-t height]\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	ldb = dbm.NewDB(util.FileNameNoExt(*dbpath), "leveldb", filepath.Dir(*dbpath))
	defer ldb.Close()

	if *action == "get" {
		get()
	} else if *action == "getall" {
		getall()
	} else if *action == "block" {
		LoadBlock(*height)
	}
}

func get() {
	data := ldb.Get([]byte(*key))
	if len(data) == 0 {
		fmt.Println(*key, "not exist")
		return
	}

	if !*decode {
		fmt.Println(string(data))
		return
	}

	if *key == util.StateKey {
		LoadState()
	} else if *key == util.AbciResponsesKey {
		LoadAbciResponses()
	} else if (*key)[0] == 'H' {
		height, _ := strconv.Atoi(strings.Split(*key, ":")[1])
		LoadBlockMeta(height)
	} else if (*key)[0] == 'P' {
		height, _ := strconv.Atoi(strings.Split(*key, ":")[1])
		index, _ := strconv.Atoi(strings.Split(*key, ":")[2])
		LoadBlockPart(height, index)
	} else if (*key)[0] == 'C' || (*key)[:2] == "SC" {
		prefix := strings.Split(*key, ":")[0]
		height, _ := strconv.Atoi(strings.Split(*key, ":")[1])
		LoadBlockCommit(height, prefix)
	} else {
		fmt.Println(string(data))
	}
}

func getall() {
	prefix := *key
	level := ldb.(*dbm.GoLevelDB).DB()
	query := level.NewIterator(nil, nil)
	//query := level.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer query.Release()
	query.Seek([]byte(prefix))
	i := 0
	for {
		fmt.Printf("%s\n", string(query.Key()))
		i++
		if !query.Next() {
			break
		}
		if i == *limit {
			break
		}
	}
	if query.Error() != nil {
		panic(query.Error())
	}
}
