# tm-tools

There are 2 tools for tendermint data:

* **tm-migrator**: migrate tendermint data from `v0.7.3` to `v0.10.0`
* **tm-viewer**: view tendermint data in `blockstore.db` or `state.db`

## tm-migrator

```
Usage: tm_migrator -old tmroot -new tmroot -priv priv_dir [-s startHeight]

	-old tmroot: dir of old tendermint root
	-new tmroot: dir of new tendermint root to store converted data
	-priv priv_dir: dir to place other validators's old `priv_validator.json`
	-s startHeight: from which height to convert tendermint data, default is `1`
```

Q: Why need `priv` arg?

A: A blockchain may consistes of one or more nodes. For every block, each node will verify `SeenCommit` and `Commit` which was produced by all validators before adding to blockchain. While `SeenCommit` and `Commit` were signed by validator, without other validators' `priv_validator.json` config info, this validaotr cannot reconstruct `SeenCommit` and `Commit` infos of the block.

## tm-viewer

```
Usage: $ tm_view -db /path/of/db [-a get|getall|block] [-q key] [-d] [-v new|old] [-t height]

// -db : db，Note: the db path cannot end with "/"
// [-a get|getall|block]： read the value of a key | output all keyes | read block info
// [-q key] ：key format, please see following "Tendermint data" section
// [-d]: whether decode value，default is "false"
// [-v new|old] ：new(0.10.0), old(0.7.3), default is "new"
// [-t height]: block height，workes with "-a block" arg to read block info at height "N"

examples：
$ tm_view -db /path/of/blockstore.db -a getall 
$ tm_view -db /path/of/blockstore.db -a block -t 1 -d 
$ tm_view -db /path/of/blockstore.db -q "H:1" -d -v old 
$ tm_view -db /path/of/state.db -q "stateKey" -d -v old 
```

## Tendermint data

In tendermint root dir(by default, is `~/.tendermint`), tendermint data is placed in its subdir `data/blockstore.db` or `data/state.db`. 

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

`state.db` is used to store `state` info, the leveldb keyes are:

* "stateKey"
* "abciResponsesKey"(v0.10.0 only)

`blockStore.db` is used to store `block` info, the leveldb keyes are(assuming the blockchain height is 32):

* "blockStore": blockchain height info. {"Height":32}
* "H:1"   ... "H:32": block meta info, 
* "P:1:0" ... "P:32:0": block part info. Block may be sliced to several parts, for each part, the key is "P:{height}:{partIndex}", partIndex start from `0`.
* "SC:1"  ... "SC:32": block seen commit info.
* "C:0"   ... "C:31": block commit info.

| key format | value type | examples | 
| ---- |-----| ---- |
| `stateKey` | raw byte of state | | 
| `abciResponsesKey` | raw byte of ABCI Responses | | 
| `blockStore` | raw json |  "blockStore": {"Height":32} | 
| `H:{height}` | raw byte of block meta | H:1 |
| `P:{height}:{index}`| raw byte of block part | P:1:0, P:32:0, P:32:1 |
| `SC:{height}` | raw byte of block seen commit | SC:1, SC:32 | 
| `C:{height-1}` | raw byte of block commit | C:0, SC:31 | 


## Tendermint Data Struct Differences

### Commit & SeenCommit

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

### BlockMeta

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

### Part

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

### Block

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

### State

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

## Further Improvement

Wating To Be Done.