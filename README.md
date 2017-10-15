# tm-tools

tendermint data

```
Allen@MacBook-Pro:~ ls -l ~/.tendermint.v0.10.0/data/
drwxr-xr-x  8 Allen  staff  272 Oct 15 20:23 blockstore.db
drwx------  3 Allen  staff  102 Oct 15 20:23 cs.wal
drwx------  3 Allen  staff  102 Oct 15 20:23 mempool.wal
drwxr-xr-x  8 Allen  staff  272 Oct 15 20:23 state.db
drwxr-xr-x  7 Allen  staff  238 Oct 15 20:23 tx_index.db
```

```
Allen@MacBook-Pro:~ ls -l ~/.tendermint.v0.7.3/data/
drwxr-xr-x   8 Allen  staff     272 Oct 15 20:22 blockstore.db
-rw-------   1 Allen  staff  829784 Oct 15 20:12 cswal
-rw-------   1 Allen  staff       3 Oct 15 20:10 mempool_wal
drwxr-xr-x  10 Allen  staff     340 Oct 15 20:22 state.db
```

```
state.db
stateKey:....
abciResponsesKey(0.10.0 only)

blockStore.db
blockStore:{"Height":32}
C:0 ....C:31
H:1 ... H:32
P:1:0 ....P:32:0
SC:1 ....SC:32

SC:1 ....SC:14
C:0 ....C:13
```

Commit & SeenCommit

```
// v0.10.0
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
```

```
// v0.7.3
type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	Precommits []*Vote `json:"precommits"`

	// Volatile
	firstPrecommit *Vote
	hash           []byte
	bitArray       *BitArray
}
```

H:1 ... H:14

```
// v0.10.0
type BlockMeta struct {
    BlockID BlockID `json:"block_id"` // the block hash and partsethash
    Header  *Header `json:"header"`   // The block's Header
}

type BlockID struct {
    Hash        data.Bytes    `json:"hash"`
    PartsHeader PartSetHeader `json:"parts"`
}
```

```
// v0.7.3
type BlockMeta struct {
	Hash        []byte        `json:"hash"`         // The block hash
	Header      *Header       `json:"header"`       // The block's Header
	PartsHeader PartSetHeader `json:"parts_header"` // The PartSetHeader, for transfer
}
```

P:1:0 ....P:14:0 for block

```
// v0.10.0
type Part struct {
    Index int                `json:"index"`
    Bytes data.Bytes         `json:"bytes"`
    Proof merkle.SimpleProof `json:"proof"`

    // Cache
    hash []byte
}
```

```
// v0.7.3
type Part struct {
	Index int                `json:"index"`
	Bytes []byte             `json:"bytes"`
	Proof merkle.SimpleProof `json:"proof"`

	// Cache
	hash []byte
}
```

block

```
// v0.10.0
type Block struct {
    *Header    `json:"header"`
    *Data      `json:"data"`
    LastCommit *Commit `json:"last_commit"`
}
```

```
// v0.7.3
type Block struct {
    *Header    `json:"header"`
    *Data      `json:"data"`
    LastCommit *Commit `json:"last_commit"`
}
```

state

```
// v0.10.0
type State struct {
    // mtx for writing to db
    mtx sync.Mutex
    db  dbm.DB

    // should not change
    GenesisDoc *types.GenesisDoc
    ChainID    string

    // updated at end of SetBlockAndValidators
    LastBlockHeight int // Genesis state has this set to 0.  So, Block(H=0) does not exist.
    LastBlockID     types.BlockID
    LastBlockTime   time.Time
    Validators      *types.ValidatorSet
    LastValidators  *types.ValidatorSet // block.LastCommit validated against this

    // AppHash is updated after Commit
    AppHash []byte

    TxIndexer txindex.TxIndexer `json:"-"` // Transaction indexer.

    // Intermediate results from processing
    // Persisted separately from the state
    abciResponses *ABCIResponses

    logger log.Logger
}
```

```
// v0.7.3
type State struct {
	mtx             sync.Mutex
	db              dbm.DB
	GenesisDoc      *types.GenesisDoc
	ChainID         string
	LastBlockHeight int // Genesis state has this set to 0.  So, Block(H=0) does not exist.
	LastBlockHash   []byte
	LastBlockParts  types.PartSetHeader
	LastBlockTime   time.Time
	Validators      *types.ValidatorSet
	LastValidators  *types.ValidatorSet
	AppHash         []byte
}
```