package convert

import (
	"github.com/hxzqlh/tm-tools/old"
	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/types"
)

// bugy
func convertBlockMeta(height int, lastBlockId *types.BlockID) *types.BlockID {
	oMeta := util.LoadOldBlockMeta(oBlockDb, height)

	nMeta := &types.BlockMeta{}
	// BlockID
	// TODO: was Hash need recompute?
	nMeta.BlockID.Hash = oMeta.Hash
	nMeta.BlockID.PartsHeader.Hash = oMeta.PartsHeader.Hash
	nMeta.BlockID.PartsHeader.Total = oMeta.PartsHeader.Total

	// Header
	nMeta.Header = NewHeader(oMeta.Header, lastBlockId)
	util.SaveNewBlockMeta(nBlockDb, height, nMeta)

	return &nMeta.BlockID
}

func NewHeader(o *old.Header, lastBlockId *types.BlockID) *types.Header {
	n := &types.Header{}
	// TODO: AppHash need reset?
	n.AppHash = o.AppHash
	n.ChainID = o.ChainID
	n.DataHash = o.DataHash
	n.Height = o.Height
	if lastBlockId != nil {
		n.LastBlockID = *lastBlockId
	}
	n.LastCommitHash = o.LastCommitHash
	n.NumTxs = o.NumTxs
	n.Time = o.Time
	n.ValidatorsHash = o.ValidatorsHash

	return n
}
