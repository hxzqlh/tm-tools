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

// -db : db 目录，路径末尾不能有 "/"
// [-a get|getall|block]： 读取某个 key 对应的 value | 获取所有 key 列表 | 读取区块信息
// [-q key] ：字符串 key
// [-d]: 是否解码二进制数据，默认 false
// [-v new|old] ：new(新版0.10.0), old 老版(0.7.3),  默认 new
// [-t height]: 区块高度，搭配 "-a block"参数使用， 获取某个高度的区块信息

tendermint 的 blockstore.db 保存了所有 block 及 commit 信息，有 4 类 key
H:{height} // 第 height 个块的头部， height>=1
P:{height}:{index} // 第 height 个块的第 index 个分片, index>=0
C:{height} // 第 height 个块的 commit 信息， height>=0
SC:{height}  //第 height 个块的 seen commit 信息，height>=1

tendermint 的 state.db 保存了区块链的最新信息, 只有 2 个 key
stateKey
abciResponsesKey

examples：
$ tm_view -db /path/of/db -a getall //输出所有 key
$ tm_view -db /path/of/blockstore.db -a block -t 1 -d //解码新版tendermint第一个块信息
$ tm_view -db /path/of/blockstore.db -q "H:1" -d -v old //解码老版tendermint第一个块的头部信息
$ tm_view -db /path/of/app/db -q "lastblock" -d //解码 app 中 key 为 lastblock 的数据
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
