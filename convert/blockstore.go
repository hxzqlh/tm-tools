package convert

import (
	"log"

	"github.com/hxzqlh/tm-tools/util"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tmlibs/db"
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

	cnt := 0
	limit := 1000
	batch := nBlockDb.NewBatch()
	for i := startHeight; i <= totalHeight; i++ {
		cnt++

		nBlock := NewBlockFromOld(i, lastBlockID)
		blockParts := nBlock.MakePartSet(types.DefaultBlockPartSize)
		nMeta := types.NewBlockMeta(nBlock, blockParts)
		// seen this BlockId's commit
		seenCommit := NewSeenCommit(i, &nMeta.BlockID)

		SaveBlock(batch, nBlock, nMeta, seenCommit)
		if cnt%limit == 0 {
			log.Printf("batch write %v/%v\n", cnt, totalHeight)
			batch.Write()
			batch = nBlockDb.NewBatch()
		}

		// update lastBlockID
		lastBlockID = &nMeta.BlockID
	}
	if cnt%limit != 0 {
		log.Printf("batch write %v/%v\n", cnt, totalHeight)
		batch.Write()
	}

	SaveState(lastBlockID)
}

func SaveBlock(batch dbm.Batch, block *types.Block, blockMeta *types.BlockMeta, seenCommit *types.Commit) {
	height := block.Height
	util.SaveNewBlockMeta2(batch, height, blockMeta)
	util.SaveNewBlockParts2(batch, height, block)
	util.SaveNewCommit2(batch, height-1, "C", block.LastCommit)
	util.SaveNewCommit2(batch, height, "SC", seenCommit)
	util.SaveNewBlockStoreStateJSON2(batch, height)
}
