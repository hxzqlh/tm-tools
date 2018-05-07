package types

import (
	"bytes"
	"errors"
	"sync"

	"github.com/tendermint/go-wire/data"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/merkle"
)

var (
	ErrPartSetUnexpectedIndex = errors.New("Error part set unexpected index")
	ErrPartSetInvalidProof    = errors.New("Error part set invalid proof")
)

type Part struct {
	Index int                `json:"index"`
	Bytes data.Bytes         `json:"bytes"`
	Proof merkle.SimpleProof `json:"proof"`

	// Cache
	hash []byte
}

type PartSetHeader struct {
	Total int        `json:"total"`
	Hash  data.Bytes `json:"hash"`
}

type PartSet struct {
	total int
	hash  []byte

	mtx           sync.Mutex
	parts         []*Part
	partsBitArray *cmn.BitArray
	count         int
}

type PartSetReader struct {
	i      int
	parts  []*Part
	reader *bytes.Reader
}
