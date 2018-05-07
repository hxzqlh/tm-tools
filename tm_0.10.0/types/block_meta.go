package types

type BlockMeta struct {
	BlockID BlockID `json:"block_id"` // the block hash and partsethash
	Header  *Header `json:"header"`   // The block's Header
}
