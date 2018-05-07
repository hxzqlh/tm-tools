package types

import (
	"errors"

	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire/data"
)

var (
	ErrVoteUnexpectedStep          = errors.New("Unexpected step")
	ErrVoteInvalidValidatorIndex   = errors.New("Invalid round vote validator index")
	ErrVoteInvalidValidatorAddress = errors.New("Invalid round vote validator address")
	ErrVoteInvalidSignature        = errors.New("Invalid round vote signature")
	ErrVoteInvalidBlockHash        = errors.New("Invalid block hash")
)

type ErrVoteConflictingVotes struct {
	VoteA *Vote
	VoteB *Vote
}

// Types of votes
// TODO Make a new type "VoteType"
const (
	VoteTypePrevote   = byte(0x01)
	VoteTypePrecommit = byte(0x02)
)

// Represents a prevote, precommit, or commit vote from validators for consensus.
type Vote struct {
	ValidatorAddress data.Bytes       `json:"validator_address"`
	ValidatorIndex   int              `json:"validator_index"`
	Height           int              `json:"height"`
	Round            int              `json:"round"`
	Type             byte             `json:"type"`
	BlockID          BlockID          `json:"block_id"` // zero if vote is nil.
	Signature        crypto.Signature `json:"signature"`
}
