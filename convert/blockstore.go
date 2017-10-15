package convert

import (
	"fmt"

	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/types"
)

// simulate the BlockStore api of blockchain
func OnBlockStore(startHeight int) {
	if startHeight < 1 {
		panic("Invalid start height")
	}

	InitState()

	var lastBlockID *types.BlockID
	if startHeight == 1 {
		lastBlockID = &types.BlockID{}
	} else {
		nMeta := util.LoadNewBlockMeta(nBlockDb, startHeight-1)
		lastBlockID = &nMeta.BlockID
	}

	for i := startHeight; i <= totalHeight; i++ {
		fmt.Printf("convert height %v/%v\n", i, totalHeight)

		nBlock := NewBlockFromOld(i, lastBlockID)
		blockParts := nBlock.MakePartSet(types.DefaultBlockPartSize)
		nMeta := types.NewBlockMeta(nBlock, blockParts)
		// seen this BlockId's commit
		seenCommit := NewSeenCommit(i, &nMeta.BlockID)
		SaveBlock(nBlock, nMeta, seenCommit)

		// update lastBlockID
		lastBlockID = &nMeta.BlockID
	}

	SaveState(lastBlockID)
}

func SaveBlock(block *types.Block, blockMeta *types.BlockMeta, seenCommit *types.Commit) {
	height := block.Height
	batch := nBlockDb.NewBatch()

	util.SaveNewBlockMeta2(batch, height, blockMeta)
	util.SaveNewBlockParts2(batch, height, block)
	util.SaveNewCommit2(batch, height-1, "C", block.LastCommit)
	util.SaveNewCommit2(batch, height, "SC", seenCommit)
	util.SaveNewBlockStoreStateJSON2(batch, height)

	batch.Write()
}
