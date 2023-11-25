package types

import "time"

// BlockHeader struct representing the header of a block
type BlockHeader struct {
	PreviousBlockHash string
	Nonce             int
	Height            int
	MerkleRootHash    string
	Timestamp         time.Time
	DifficultyTarget  int
}
