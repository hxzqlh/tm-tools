package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	old "github.com/hxzqlh/tm-tools/tm_0.10.0/types"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	stateKey = "stateKey"
)

// 0.10.0  --> 0.16.0
// ./state -old dir -new dir
var oldDir = flag.String("old", os.ExpandEnv("$HOME/.tendermint")+"/data/state.db", "tendermint state db")
var newDir = flag.String("new", "", "tendermint state db")
var oldDb dbm.DB
var newDb dbm.DB

func main() {
	if len(os.Args) <= 4 {
		fmt.Printf("Usage: %s -old dir -new dir\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	oldDb = dbm.NewDB("state", "leveldb", filepath.Dir(*oldDir))
	defer oldDb.Close()

	newDb = dbm.NewDB("state", "leveldb", filepath.Dir(*newDir))
	defer newDb.Close()

	oldState := old.LoadState(oldDb)
	res, _ := json.Marshal(oldState)
	fmt.Println(string(res))

	newState := convert(oldState)
	res, _ = json.Marshal(newState)
	fmt.Println(string(res))

	state.SaveState(newDb, *newState)
}

func convert(s *old.State) *state.State {
	ss := &state.State{}
	ss.ChainID = s.ChainID
	ss.LastBlockHeight = int64(s.LastBlockHeight)
	//ss.LastBlockTotalTx
	ss.LastBlockID = copyLastBlockID(s.LastBlockID)
	ss.LastBlockTime = s.LastBlockTime
	ss.Validators = NewValidatorSet(s.Validators)
	ss.LastValidators = NewValidatorSet(s.LastValidators)
	ss.AppHash = s.AppHash
	return ss
}

func copyLastBlockID(b old.BlockID) types.BlockID {
	return types.BlockID{
		Hash: cmn.HexBytes(b.Hash),
		PartsHeader: types.PartSetHeader{
			Total: b.PartsHeader.Total,
			Hash:  cmn.HexBytes(b.PartsHeader.Hash),
		},
	}
}

func NewValidatorSet(oValidatorSet *old.ValidatorSet) *types.ValidatorSet {
	nValidatorSet := &types.ValidatorSet{
		Validators: []*types.Validator{},
	}

	// Validators
	for _, val := range oValidatorSet.Validators {
		one := &types.Validator{}
		one.Accum = val.Accum
		one.Address = cmn.HexBytes(val.Address)
		one.PubKey = val.PubKey
		one.VotingPower = val.VotingPower

		nValidatorSet.Validators = append(nValidatorSet.Validators, one)
	}

	// Proposer
	// NOTE: no this fiedl in old version, we find by new api
	nValidatorSet.Proposer = nValidatorSet.GetProposer()

	return nValidatorSet
}
