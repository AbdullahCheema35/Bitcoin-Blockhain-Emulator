package types

// Block struct containing both header and body
type Block struct {
	BlockHash string
	Header    BlockHeader
	Body      BlockBody
}
