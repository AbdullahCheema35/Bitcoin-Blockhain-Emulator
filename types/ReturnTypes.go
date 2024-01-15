package types

type ReturnType uint8

const (
	MiningProofFailed ReturnType = iota
	MerkleRootFailed
	AddBlockFailed
	NewBlockVerificationFailed
	PreviousBlockHashFailed
	BlockHeightFailed
	GenesisBlockHeightFailed
	DuplicateTransactionsFailed
	BlockChainVerificationSuccessful
	BlockVerificationSuccessful
	NewHeightLEQCurrentHeight
	NewBlockAddedSuccessfully
	NewBlockPrevHashDontMatch
	NewBlockDuplicateTransactions
	InitiateBroadcastBlockChainRequest
	InitiateBroadcastBlock
	DoNothing
)
