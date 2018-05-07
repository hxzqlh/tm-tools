package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hxzqlh/tm-tools/old"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/tendermint/blockchain"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	BlockStoreKey    = "blockStore"
	StateKey         = "stateKey"
	AbciResponsesKey = "abciResponsesKey"
)

// blockstore json
func LoadOldBlockStoreStateJSON(ldb dbm.DB) old.BlockStoreStateJSON {
	bytes := ldb.Get([]byte(BlockStoreKey))
	if bytes == nil {
		return old.BlockStoreStateJSON{
			Height: 0,
		}
	}
	bsj := old.BlockStoreStateJSON{}
	err := json.Unmarshal(bytes, &bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal bytes: %X", bytes))
	}
	return bsj
}

func SaveNewBlockStoreStateJSON(ldb dbm.DB, totalHeight int) {
	bsj := blockchain.BlockStoreStateJSON{
		Height: int64(totalHeight),
	}
	bytes, err := json.Marshal(bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal state bytes: %v", err))
	}
	ldb.Set([]byte(BlockStoreKey), bytes)
}

func SaveNewBlockStoreStateJSON2(batch dbm.Batch, totalHeight int) {
	bsj := blockchain.BlockStoreStateJSON{
		Height: int64(totalHeight),
	}
	bytes, err := json.Marshal(bsj)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal state bytes: %v", err))
	}
	batch.Set([]byte(BlockStoreKey), bytes)
}

// block
func LoadNewBlock(ldb dbm.DB, height int) *types.Block {
	var n int
	var err error
	bytez := []byte{}
	meta := LoadNewBlockMeta(ldb, height)
	for i := 0; i < meta.BlockID.PartsHeader.Total; i++ {
		part := LoadNewBlockPart(ldb, height, i)
		bytez = append(bytez, part.Bytes...)
	}
	block := wire.ReadBinary(&types.Block{}, bytes.NewReader(bytez), 0, &n, &err).(*types.Block)
	if err != nil {
		panic(err)
	}
	return block
}

func LoadOldBlock(ldb dbm.DB, height int) *old.Block {
	var n int
	var err error
	bytez := []byte{}
	meta := LoadOldBlockMeta(ldb, height)
	for i := 0; i < meta.PartsHeader.Total; i++ {
		part := LoadOldBlockPart(ldb, height, i)
		bytez = append(bytez, part.Bytes...)
	}
	block := wire.ReadBinary(&old.Block{}, bytes.NewReader(bytez), 0, &n, &err).(*old.Block)
	if err != nil {
		panic(err)
	}
	return block
}

// header
func LoadNewBlockMeta(ldb dbm.DB, height int) *types.BlockMeta {
	var n int
	var err error
	buf := ldb.Get(calcBlockMetaKey(height))
	meta := wire.ReadBinary(&types.BlockMeta{}, bytes.NewReader(buf), 0, &n, &err).(*types.BlockMeta)
	if err != nil {
		panic(err)
	}
	return meta
}

func LoadOldBlockMeta(ldb dbm.DB, height int) *old.BlockMeta {
	var n int
	var err error
	buf := ldb.Get(calcBlockMetaKey(height))
	meta := wire.ReadBinary(&old.BlockMeta{}, bytes.NewReader(buf), 0, &n, &err).(*old.BlockMeta)
	if err != nil {
		panic(err)
	}
	return meta
}

func SaveNewBlockMeta(ldb dbm.DB, height int, blockMeta *types.BlockMeta) {
	metaBytes := wire.BinaryBytes(blockMeta)
	ldb.Set(calcBlockMetaKey(height), metaBytes)
}

func SaveNewBlockMeta2(batch dbm.Batch, height int, blockMeta *types.BlockMeta) {
	metaBytes := wire.BinaryBytes(blockMeta)
	batch.Set(calcBlockMetaKey(height), metaBytes)
}

// part
func LoadNewBlockPart(ldb dbm.DB, height int, index int) *types.Part {
	buf := ldb.Get(calcBlockPartKey(height, index))
	r, n, err := bytes.NewReader(buf), new(int), new(error)
	part := wire.ReadBinary(&types.Part{}, r, 0, n, err).(*types.Part)
	if *err != nil {
		panic(*err)
	}
	return part
}

func LoadOldBlockPart(ldb dbm.DB, height int, index int) *old.Part {
	buf := ldb.Get(calcBlockPartKey(height, index))
	r, n, err := bytes.NewReader(buf), new(int), new(error)
	part := wire.ReadBinary(&old.Part{}, r, 0, n, err).(*old.Part)
	if *err != nil {
		panic(*err)
	}
	return part
}

func SaveNewBlockParts(ldb dbm.DB, height int, block *types.Block) {
	blockParts := block.MakePartSet(types.DefaultBlockPartSize)
	for index := 0; index < blockParts.Total(); index++ {
		SaveNewBlockPart(ldb, height, index, blockParts.GetPart(index))
	}
}

func SaveNewBlockParts2(batch dbm.Batch, height int, block *types.Block) {
	blockParts := block.MakePartSet(types.DefaultBlockPartSize)
	for index := 0; index < blockParts.Total(); index++ {
		SaveNewBlockPart2(batch, height, index, blockParts.GetPart(index))
	}
}

func SaveNewBlockPart(ldb dbm.DB, height int, index int, part *types.Part) {
	partBytes := wire.BinaryBytes(part)
	ldb.Set(calcBlockPartKey(height, index), partBytes)
}

func SaveNewBlockPart2(batch dbm.Batch, height int, index int, part *types.Part) {
	partBytes := wire.BinaryBytes(part)
	batch.Set(calcBlockPartKey(height, index), partBytes)
}

// commit
func LoadNewBlockCommit(ldb dbm.DB, height int, prefix string) *types.Commit {
	var buf []byte

	if prefix == "C" {
		buf = ldb.Get(calcBlockCommitKey(height))
	} else if prefix == "SC" {
		buf = ldb.Get(calcSeenCommitKey(height))
	}

	blockCommit := &types.Commit{}
	r, n, err := bytes.NewReader(buf), new(int), new(error)
	wire.ReadBinaryPtr(&blockCommit, r, 0, n, err)
	if *err != nil {
		panic(err)
	}
	return blockCommit
}

func LoadOldBlockCommit(ldb dbm.DB, height int, prefix string) *old.Commit {
	var buf []byte

	if prefix == "C" {
		buf = ldb.Get(calcBlockCommitKey(height))
	} else if prefix == "SC" {
		buf = ldb.Get(calcSeenCommitKey(height))
	}

	blockCommit := &old.Commit{}

	r, n, err := bytes.NewReader(buf), new(int), new(error)
	wire.ReadBinaryPtr(&blockCommit, r, 0, n, err)
	if *err != nil {
		panic(err)
	}
	return blockCommit
}

func SaveNewCommit(ldb dbm.DB, height int, prefix string, commit *types.Commit) {
	var key []byte
	switch prefix {
	case "C":
		key = calcBlockCommitKey(height)
	case "SC":
		key = calcSeenCommitKey(height)
	default:
		panic(prefix)
	}

	buf := wire.BinaryBytes(commit)
	ldb.Set(key, buf)
}

func SaveNewCommit2(batch dbm.Batch, height int, prefix string, commit *types.Commit) {
	var key []byte
	switch prefix {
	case "C":
		key = calcBlockCommitKey(height)
	case "SC":
		key = calcSeenCommitKey(height)
	default:
		panic(prefix)
	}

	buf := wire.BinaryBytes(commit)
	batch.Set(key, buf)
}

// state
func LoadOldState(ldb dbm.DB) *old.State {
	buf := ldb.Get([]byte(StateKey))
	s := &old.State{}
	r, n, err := bytes.NewReader(buf), new(int), new(error)
	wire.ReadBinaryPtr(&s, r, 0, n, err)
	if *err != nil {
		panic(err)
	}
	return s
}

func LoadNewState(ldb dbm.DB) *state.State {
	buf := ldb.Get([]byte(StateKey))
	s := &state.State{}
	r, n, err := bytes.NewReader(buf), new(int), new(error)
	wire.ReadBinaryPtr(&s, r, 0, n, err)
	if *err != nil {
		panic(err)
	}
	return s
}

func SaveNewState(ldb dbm.DB, s *state.State) {
	buf := s.Bytes()
	ldb.Set([]byte(StateKey), buf)
}

// abciResps
func LoadAbciResps(ldb dbm.DB) *state.ABCIResponses {
	buf := ldb.Get([]byte(AbciResponsesKey))
	resps := &state.ABCIResponses{}

	r, n, err := bytes.NewReader(buf), new(int), new(error)
	wire.ReadBinaryPtr(resps, r, 0, n, err)
	if *err != nil {
		panic(err)
	}
	return resps
}

//==============================================================================

func calcBlockMetaKey(height int) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}

func calcBlockPartKey(height int, partIndex int) []byte {
	return []byte(fmt.Sprintf("P:%v:%v", height, partIndex))
}

func calcBlockCommitKey(height int) []byte {
	return []byte(fmt.Sprintf("C:%v", height))
}

func calcSeenCommitKey(height int) []byte {
	return []byte(fmt.Sprintf("SC:%v", height))
}
