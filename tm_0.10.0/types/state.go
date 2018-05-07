package types

import (
	"bytes"
	"sync"
	"time"

	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	abci "github.com/hxzqlh/tm-tools/tm_0.10.0/types/abci/types"
	wire "github.com/tendermint/go-wire"
)

var (
	stateKey         = []byte("stateKey")
	abciResponsesKey = []byte("abciResponsesKey")
)

//-----------------------------------------------------------------------------

// NOTE: not goroutine-safe.
type State struct {
	// mtx for writing to db
	mtx sync.Mutex
	db  dbm.DB

	// should not change
	GenesisDoc *GenesisDoc
	ChainID    string

	// updated at end of SetBlockAndValidators
	LastBlockHeight int // Genesis state has this set to 0.  So, Block(H=0) does not exist.
	LastBlockID     BlockID
	LastBlockTime   time.Time
	Validators      *ValidatorSet
	LastValidators  *ValidatorSet // block.LastCommit validated against this

	// AppHash is updated after Commit
	AppHash []byte

	// Intermediate results from processing
	// Persisted separately from the state
	abciResponses *ABCIResponses

	logger log.Logger
}

func LoadState(db dbm.DB) *State {
	return loadState(db, stateKey)
}

func loadState(db dbm.DB, key []byte) *State {
	s := &State{db: db}
	buf := db.Get(key)
	if len(buf) == 0 {
		return nil
	} else {
		r, n, err := bytes.NewReader(buf), new(int), new(error)
		wire.ReadBinaryPtr(&s, r, 0, n, err)
		if *err != nil {
			// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
			cmn.Exit(cmn.Fmt("LoadState: Data has been corrupted or its spec has changed: %v\n", *err))
		}
		// TODO: ensure that buf is completely read.
	}
	return s
}

func (s *State) SetLogger(l log.Logger) {
	s.logger = l
}

func (s *State) Save() {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.db.SetSync(stateKey, s.Bytes())
}

func (s *State) Bytes() []byte {
	buf, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(s, buf, n, err)
	if *err != nil {
		cmn.PanicCrisis(*err)
	}
	return buf.Bytes()
}

//--------------------------------------------------
// ABCIResponses holds intermediate state during block processing

type ABCIResponses struct {
	Height int

	DeliverTx []*abci.ResponseDeliverTx
	EndBlock  abci.ResponseEndBlock

	txs Txs // reference for indexing results by hash
}
