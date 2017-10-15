package convert

import (
	"github.com/hxzqlh/tm-tools/old"
	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/types"
)

func convertPart(height int, lastBlockID *types.BlockID) {
	nBlock := NewBlockFromOld(height, lastBlockID)
	util.SaveNewBlockParts(nBlockDb, height, nBlock)
}

func NewBlockFromOld(height int, lastBlockID *types.BlockID) *types.Block {
	oBlock := util.LoadOldBlock(oBlockDb, height)

	nBlock := &types.Block{}
	nBlock.Data = NewData(oBlock.Data)
	nBlock.LastCommit = NewCommit(oBlock.LastCommit, lastBlockID)
	nBlock.Header = NewHeader(oBlock.Header, lastBlockID)

	return nBlock
}

func NewBlockFromOld2(height int, lastBlockID *types.BlockID) (*types.Block, *types.PartSet) {
	oBlock := util.LoadOldBlock(oBlockDb, height)
	commit := NewCommit(oBlock.LastCommit, lastBlockID)
	txs := []types.Tx{}
	for _, tx := range oBlock.Txs {
		txs = append(txs, []byte(tx))
	}

	return types.MakeBlock(height, oBlock.ChainID, txs, commit,
		*lastBlockID, oBlock.ValidatorsHash, oBlock.AppHash, types.DefaultBlockPartSize)
}

func NewData(o *old.Data) *types.Data {
	nData := &types.Data{}
	txs := []types.Tx{}
	for _, tx := range o.Txs {
		txs = append(txs, []byte(tx))
	}
	nData.Txs = txs
	return nData
}
