package convert

import (
	"github.com/hxzqlh/tm-tools/old"
	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tendermint/types"
)

func NewSeenCommit(height int, lastBlockID *types.BlockID) *types.Commit {
	oCommit := util.LoadOldBlockCommit(oBlockDb, height, "SC")
	return NewCommit(oCommit, lastBlockID)
}

func NewCommit(oCommit *old.Commit, lastBlockID *types.BlockID) *types.Commit {
	nCommit := &types.Commit{}

	preCommits := []*types.Vote{}
	for i := 0; i < len(oCommit.Precommits); i++ {
		v := oCommit.Precommits[i]
		// node's commit may be nil
		if v == nil {
			preCommits = append(preCommits, nil)
			continue
		}

		one := &types.Vote{}
		one.BlockID = *lastBlockID
		one.Height = v.Height
		one.Round = v.Round
		one.Type = v.Type
		one.ValidatorIndex = i
		one.ValidatorAddress = nState.Validators.Validators[i].Address

		one.Signature = SignaVote(i, nState.ChainID, one)
		preCommits = append(preCommits, one)
	}

	nCommit.BlockID = *lastBlockID
	nCommit.Precommits = preCommits

	return nCommit
}

func convertSeenCommit(height int, lastBlockID *types.BlockID) {
	oCommit := util.LoadOldBlockCommit(oBlockDb, height, "SC")
	nCommit := NewCommit(oCommit, lastBlockID)
	util.SaveNewCommit(nBlockDb, height, "SC", nCommit)
}

func convertCommit(height int, lastBlockID *types.BlockID) {
	oCommit := util.LoadOldBlockCommit(oBlockDb, height, "C")
	nCommit := NewCommit(oCommit, lastBlockID)
	util.SaveNewCommit(nBlockDb, height, "C", nCommit)
}

func SignaVote(index int, chainID string, vote *types.Vote) crypto.Signature {
	if index >= len(privValidators.Validators) {
		panic("privValidators index overflow")
	}

	bytez := types.SignBytes(chainID, vote)
	return privValidators.Validators[index].PrivKey.Sign(bytez)
}
