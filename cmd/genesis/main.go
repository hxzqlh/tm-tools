package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hxzqlh/tm-tools/old"
	"github.com/tendermint/tendermint/types"
)

// convert tendermint v0.7.3 genesis.json to version of v0.10.0 and output result to stdout
// Usage: ./tm_genesis old_genesis_path
/*
Replace [TypeByte, Xxx] with {"type": "some-type", "data": Xxx} in RPC and all .json files by using go-wire/data.

For instance, a pubkey old verison is:
"pub_key": {
	1,
	"83DDF8775937A4A12A2704269E2729FCFCD491B933C4B0A7FFE37FE41D7760D0"
}

now is:
"pub_key": {
	"type": "ed25519",
	"data": "83DDF8775937A4A12A2704269E2729FCFCD491B933C4B0A7FFE37FE41D7760D0"
}
*/

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %s old_genesis.json \n", os.Args[0])
		os.Exit(0)
	}

	jsonBytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	oGen := old.GenesisDocFromJSON(jsonBytes)
	nGen := NewGenesisDoc(oGen)
	bytes, _ := json.Marshal(nGen)
	fmt.Println(string(bytes))
}

func NewGenesisDoc(old *old.GenesisDoc) *types.GenesisDoc {
	newGenesisDoc := &types.GenesisDoc{
		AppHash:     old.AppHash,
		ChainID:     old.ChainID,
		GenesisTime: old.GenesisTime,
		Validators:  []types.GenesisValidator{},
	}
	for _, val := range old.Validators {
		one := types.GenesisValidator{}
		one.Amount = val.Amount
		one.Name = val.Name
		one.PubKey = val.PubKey

		newGenesisDoc.Validators = append(newGenesisDoc.Validators, one)
	}

	return newGenesisDoc
}
