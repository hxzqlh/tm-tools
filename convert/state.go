package convert

import (
	"github.com/hxzqlh/tm-tools/old"
	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
)

var nState *state.State

func InitState() {
	oState := util.LoadOldState(oStateDb)
	nState = &state.State{}
	nState.AppHash = oState.AppHash
	nState.ChainID = oState.ChainID
	nState.GenesisDoc = NewGenesisDoc(oState.GenesisDoc)
	nState.LastBlockHeight = oState.LastBlockHeight
	nState.LastBlockTime = oState.LastBlockTime
	nState.LastValidators = NewValidatorSet(oState.LastValidators)
	nState.Validators = NewValidatorSet(oState.Validators)
	// need set LastBlockID before save state
}

func SaveState(lastBlockID *types.BlockID) {
	nState.LastBlockID = *lastBlockID
	util.SaveNewState(nStateDb, nState)
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

func NewValidatorSet(oValidatorSet *old.ValidatorSet) *types.ValidatorSet {
	nValidatorSet := &types.ValidatorSet{
		Validators: []*types.Validator{},
	}

	// Validators
	for _, val := range oValidatorSet.Validators {
		one := &types.Validator{}
		one.Accum = val.Accum
		one.Address = val.Address
		one.PubKey = val.PubKey
		one.VotingPower = val.VotingPower

		nValidatorSet.Validators = append(nValidatorSet.Validators, one)
	}

	// Proposer
	// NOTE: no this fiedl in old version, we find by new api
	nValidatorSet.Proposer = nValidatorSet.GetProposer()

	return nValidatorSet
}
