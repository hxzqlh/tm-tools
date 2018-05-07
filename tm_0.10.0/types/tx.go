package types

import (
	abci "github.com/hxzqlh/tm-tools/tm_0.10.0/types/abci/types"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/tmlibs/merkle"
)

type Tx []byte

type Txs []Tx

type TxProof struct {
	Index, Total int
	RootHash     data.Bytes
	Data         Tx
	Proof        merkle.SimpleProof
}

// TxResult contains results of executing the transaction.
//
// One usage is indexing transaction results.
type TxResult struct {
	Height uint64                 `json:"height"`
	Index  uint32                 `json:"index"`
	Tx     Tx                     `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}
