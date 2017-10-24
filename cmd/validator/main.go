package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hxzqlh/tm-tools/old"
	"github.com/tendermint/tendermint/types"
)

// convert tendermint v0.7.3 priv_validator.json to version of v0.10.0 and output result to stdout
// Usage: ./tm_priv_validator old_priv_validator_path
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
		fmt.Printf("Usage: %s old_priv_validator.json \n", os.Args[0])
		os.Exit(0)
	}

	nVali := NewPrivValidator(os.Args[1])
	bytes, _ := json.Marshal(nVali)
	fmt.Println(string(bytes))
}

func NewPrivValidator(oPath string) *types.PrivValidator {
	privVali := &types.PrivValidator{}
	old := old.LoadPrivValidator(oPath)
	privVali.Address = old.Address
	privVali.LastHeight = old.LastHeight
	privVali.LastRound = old.LastRound
	privVali.LastSignature = old.LastSignature
	privVali.LastSignBytes = old.LastSignBytes
	privVali.LastStep = old.LastStep
	privVali.PrivKey = old.PrivKey
	privVali.PubKey = old.PubKey
	return privVali
}
