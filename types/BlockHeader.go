package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// BlockHeader struct representing the header of a block
type BlockHeader struct {
	PreviousBlockHash string
	Nonce             int
	Height            int
	MerkleRootHash    string
	Timestamp         time.Time
	DifficultyTarget  int
}

// Function to calculate the hash of the block
func (bh *BlockHeader) RecalculateBlockHeaderHash() (string, []byte) {
	hashInput := bh.PreviousBlockHash + fmt.Sprintf("%d", bh.Nonce) + fmt.Sprintf("%d", bh.Timestamp.UnixMicro()) + bh.MerkleRootHash + fmt.Sprintf("%d", bh.DifficultyTarget)
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:]), hash[:]
}
