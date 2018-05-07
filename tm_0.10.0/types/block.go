package types

import (
	"time"

	"github.com/tendermint/go-wire/data"
	. "github.com/tendermint/tmlibs/common"
)

const (
	MaxBlockSize         = 22020096 // 21MB TODO make it configurable
	DefaultBlockPartSize = 65536    // 64kB TODO: put part size in parts header?
)

type Block struct {
	*Header    `json:"header"`
	*Data      `json:"data"`
	LastCommit *Commit `json:"last_commit"`
}

type Header struct {
	ChainID        string     `json:"chain_id"`
	Height         int        `json:"height"`
	Time           time.Time  `json:"time"`
	NumTxs         int        `json:"num_txs"` // XXX: Can we get rid of this?
	LastBlockID    BlockID    `json:"last_block_id"`
	LastCommitHash data.Bytes `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       data.Bytes `json:"data_hash"`        // transactions
	ValidatorsHash data.Bytes `json:"validators_hash"`  // validators for the current block
	AppHash        data.Bytes `json:"app_hash"`         // state after txs from the previous block
}

// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID    BlockID `json:"blockID"`
	Precommits []*Vote `json:"precommits"`

	// Volatile
	firstPrecommit *Vote
	hash           data.Bytes
	bitArray       *BitArray
}

type Data struct {

	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs Txs `json:"txs"`

	// Volatile
	hash data.Bytes
}

type BlockID struct {
	Hash        data.Bytes    `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}
