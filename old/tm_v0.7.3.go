package old

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/merkle"
)

// following structs are based on tendermint v0.7.3

type Block struct {
	*Header    `json:"header"`
	*Data      `json:"data"`
	LastCommit *Commit `json:"last_commit"`
}

type BlockStoreStateJSON struct {
	Height int
}

type BlockMeta struct {
	Hash        []byte        `json:"hash"`         // The block hash
	Header      *Header       `json:"header"`       // The block's Header
	PartsHeader PartSetHeader `json:"parts_header"` // The PartSetHeader, for transfer
}

type Header struct {
	ChainID        string        `json:"chain_id"`
	Height         int           `json:"height"`
	Time           time.Time     `json:"time"`
	NumTxs         int           `json:"num_txs"`
	LastBlockHash  []byte        `json:"last_block_hash"`
	LastBlockParts PartSetHeader `json:"last_block_parts"`
	LastCommitHash []byte        `json:"last_commit_hash"`
	DataHash       []byte        `json:"data_hash"`
	ValidatorsHash []byte        `json:"validators_hash"`
	AppHash        []byte        `json:"app_hash"` // state merkle root of txs from the previous block
}

type PartSetHeader struct {
	Total int    `json:"total"`
	Hash  []byte `json:"hash"`
}

type Data struct {
	// Txs that will be applied by state @ block.Height+1.
	// NOTE: not all txs here are valid.  We're just agreeing on the order first.
	// This means that block.AppHash does not include these txs.
	Txs Txs `json:"txs"`
	// Volatile
	hash []byte
}

type Tx []byte
type Txs []Tx

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

type BitArray struct {
	mtx   sync.Mutex
	Bits  int      `json:"bits"`  // NOTE: persisted via reflect, must be exported
	Elems []uint64 `json:"elems"` // NOTE: persisted via reflect, must be exported
}

type Vote struct {
	Height           int                     `json:"height"`
	Round            int                     `json:"round"`
	Type             byte                    `json:"type"`
	BlockHash        []byte                  `json:"block_hash"`         // empty if vote is nil.
	BlockPartsHeader PartSetHeader           `json:"block_parts_header"` // zero if vote is nil.
	Signature        crypto.SignatureEd25519 `json:"signature"`
}

type Part struct {
	Index int                `json:"index"`
	Bytes []byte             `json:"bytes"`
	Proof merkle.SimpleProof `json:"proof"`
	// Cache
	hash []byte
}

type State struct {
	mtx             sync.Mutex
	db              dbm.DB
	GenesisDoc      *GenesisDoc
	ChainID         string
	LastBlockHeight int // Genesis state has this set to 0.  So, Block(H=0) does not exist.
	LastBlockHash   []byte
	LastBlockParts  PartSetHeader
	LastBlockTime   time.Time
	Validators      *ValidatorSet
	LastValidators  *ValidatorSet
	AppHash         []byte
}

type GenesisDoc struct {
	GenesisTime time.Time          `json:"genesis_time"`
	ChainID     string             `json:"chain_id"`
	Validators  []GenesisValidator `json:"validators"`
	AppHash     []byte             `json:"app_hash"`
}

type GenesisValidator struct {
	PubKey crypto.PubKey `json:"pub_key"`
	Amount int64         `json:"amount"`
	Name   string        `json:"name"`
}

type ValidatorSet struct {
	Validators []*Validator // NOTE: persisted via reflect, must be exported.
	// cached (unexported)
	proposer         *Validator
	totalVotingPower int64
}

type Validator struct {
	Address          []byte        `json:"address"`
	PubKey           crypto.PubKey `json:"pub_key"`
	LastCommitHeight int           `json:"last_commit_height"`
	VotingPower      int64         `json:"voting_power"`
	Accum            int64         `json:"accum"`
}

type PrivValidator struct {
	Address       []byte           `json:"address"`
	PubKey        crypto.PubKey    `json:"pub_key"`
	LastHeight    int              `json:"last_height"`
	LastRound     int              `json:"last_round"`
	LastStep      int8             `json:"last_step"`
	LastSignature crypto.Signature `json:"last_signature"` // so we dont lose signatures
	LastSignBytes []byte           `json:"last_signbytes"` // so we dont lose signatures
	// PrivKey should be empty if a Signer other than the default is being used.
	PrivKey crypto.PrivKey `json:"priv_key"`
	Signer  `json:"-"`
	// For persistence.
	// Overloaded for testing.
	filePath string
	mtx      sync.Mutex
}

type Signer interface {
	Sign(msg []byte) crypto.Signature
}

// Implements Signer
type DefaultSigner struct {
	priv crypto.PrivKey
}

func NewDefaultSigner(priv crypto.PrivKey) *DefaultSigner {
	return &DefaultSigner{priv: priv}
}

// Implements Signer
func (ds *DefaultSigner) Sign(msg []byte) crypto.Signature {
	return ds.priv.Sign(msg)
}

func GenesisDocFromJSON(jsonBlob []byte) (genState *GenesisDoc) {
	var err error
	wire.ReadJSONPtr(&genState, jsonBlob, &err)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read GenesisDoc: %v", err))
	}
	return
}

func LoadPrivValidator(filePath string) *PrivValidator {
	privValJSONBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err.Error())
	}
	privVal := wire.ReadJSON(&PrivValidator{}, privValJSONBytes, &err).(*PrivValidator)
	if err != nil {
		panic(fmt.Sprintf("Error reading PrivValidator from %v: %v\n", filePath, err))
	}
	privVal.filePath = filePath
	privVal.Signer = NewDefaultSigner(privVal.PrivKey)
	return privVal
}

// cswal
type ConsensusLogMessage struct {
	Time time.Time                    `json:"time"`
	Msg  ConsensusLogMessageInterface `json:"msg"`
}

type ConsensusLogMessageInterface interface{}

// 0x01
type EventDataRoundState struct {
	Height int    `json:"height"`
	Round  int    `json:"round"`
	Step   string `json:"step"`
	// private, not exposed to websockets
	RoundState interface{} `json:"-"`
}

//0x02
type MsgInfo struct {
	Msg     ConsensusMessage `json:"msg"`
	PeerKey string           `json:"peer_key"`
}

type ConsensusMessage interface{}

// 0x03
// internally generated messages which may update the state
type TimeoutInfo struct {
	Duration time.Duration `json:"duration"`
	Height   int           `json:"height"`
	Round    int           `json:"round"`
	Step     RoundStepType `json:"step"`
}

type RoundStepType uint8 // These must be numeric, ordered.
