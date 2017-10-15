package main

import (
	"encoding/json"
	"fmt"

	"github.com/hxzqlh/tm-tools/util"
)

func LoadBlock(height int) {
	var res []byte

	switch *ver {
	case "new":
		block := util.LoadNewBlock(ldb, height)
		res, _ = json.Marshal(block)
	case "old":
		block := util.LoadOldBlock(ldb, height)
		res, _ = json.Marshal(block)
	default:
		panic(ver)
	}

	fmt.Println(string(res))
}

//state
func LoadState() {
	var res []byte

	switch *ver {
	case "new":
		s := util.LoadNewState(ldb)
		res, _ = json.Marshal(s)
	case "old":
		s := util.LoadOldState(ldb)
		res, _ = json.Marshal(s)
	default:
		panic(ver)
	}

	fmt.Println(string(res))
}

//meta
func LoadBlockMeta(height int) {
	var res []byte

	switch *ver {
	case "new":
		meta := util.LoadNewBlockMeta(ldb, height)
		res, _ = json.Marshal(meta)
	case "old":
		meta := util.LoadOldBlockMeta(ldb, height)
		res, _ = json.Marshal(meta)
	default:
		panic(ver)
	}

	fmt.Println(string(res))
}

//Part
func LoadBlockPart(height int, index int) {
	var res []byte

	switch *ver {
	case "new":
		part := util.LoadNewBlockPart(ldb, height, index)
		res, _ = json.Marshal(part)
	case "old":
		part := util.LoadOldBlockPart(ldb, height, index)
		res, _ = json.Marshal(part)
	default:
		panic(ver)
	}

	fmt.Println(string(res))
}

//commit
func LoadBlockCommit(height int, prefix string) {
	var res []byte

	switch *ver {
	case "new":
		commit := util.LoadNewBlockCommit(ldb, height, prefix)
		res, _ = json.Marshal(commit)
	case "old":
		commit := util.LoadOldBlockCommit(ldb, height, prefix)
		res, _ = json.Marshal(commit)
	default:
		panic(ver)
	}

	fmt.Println(string(res))
}

func LoadAbciResponses() {
	abciResps := util.LoadAbciResps(ldb)
	res, _ := json.Marshal(abciResps)
	fmt.Println(string(res))
}
