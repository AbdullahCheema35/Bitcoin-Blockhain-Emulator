package mining

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func MineBlock(b types.Block) {
	miningProcessAbortedChan := nodestate.InitMiningChan()
	fmt.Println("Mining block...")

	difficultyTarget := b.Header.DifficultyTarget
	_, targetHashBytes := GenerateTargetHash(difficultyTarget)

	header := b.Header
	header.Nonce = 0

	minedSuccessfully := false

	for {
		_, hashBytes := header.RecalculateBlockHeaderHash()
		if compareWithTargetHash(hashBytes, targetHashBytes) {
			minedSuccessfully = true
		}

		select {
		case <-miningProcessAbortedChan:
			return
		default:
			header.Nonce++
		}
		if minedSuccessfully {
			fmt.Printf("Block mined! Nonce: %d\n", header.Nonce)
			break
		}
		// Sleep for random amount of time less than 100 milliseconds
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	b.Header = header
	b.RecalculateBlockHash()
}

func GenerateTargetHash(difficultyTarget int) (string, []byte) {
	// returns a 256 bit hash
	// difficultyTarget is the number of leading zeros
	// hexadecimal string representation of the hash will contain 64 characters
	// 4 bits per character
	if difficultyTarget < 0 {
		difficultyTarget = 0
	} else if difficultyTarget > 256 {
		difficultyTarget = 256
	}
	totalCharacters := 64
	numZeros := difficultyTarget / 4
	leftover := difficultyTarget % 4

	// create a hex string with the required number of zeros
	// and the required number of leftover characters
	switch leftover {
	case 0:
		targetString := strings.Repeat("0", numZeros) + strings.Repeat("f", totalCharacters-numZeros)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 1:
		targetString := strings.Repeat("0", numZeros) + "7" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 2:
		targetString := strings.Repeat("0", numZeros) + "3" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	case 3:
		targetString := strings.Repeat("0", numZeros) + "1" + strings.Repeat("f", totalCharacters-numZeros-1)
		// convert the hex string to a byte array
		targetBytes, _ := hex.DecodeString(targetString)
		return targetString, targetBytes
	}
	return "", nil
}

// returns true if the hash is less than or equal to the target
// returns false if the hash is greater than the target
func compareWithTargetHash(hash []byte, target []byte) bool {
	// must be the same length
	if len(hash) != len(target) || len(hash) != 32 {
		fmt.Println("Hash and target are not the same length")
		return false
	}
	for i := 0; i < len(hash); i++ {
		if hash[i] < target[i] {
			return true
		} else if hash[i] > target[i] {
			return false
		}
	}
	return true
}

// func CompareWithTargetHashString(hashString string, targetString string) bool {
// 	if len(hashString) != len(targetString) || len(hashString) != 64 {
// 		fmt.Println("Hash and target are not the same length")
// 		return false
// 	}
// 	hash, _ := hex.DecodeString(hashString)
// 	target, _ := hex.DecodeString(targetString)
// 	return compareWithTargetHash(hash, target)
// }

func AddNewBlockToBlockChain(b types.Block) ReturnType {
	ret := VerifyBlock(b)
	if ret != BlockVerificationSuccessful {
		return NewBlockVerificationFailed
	}

	newHeight := b.Header.Height

	// Get and lock the the block chain
	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		return NewHeightLEQCurrentHeight
	}

	// Verify that the previous block hash matches the hash of the latest block in the blockchain
	latestBlockHash := bchain.GetLatestBlockHash()
	if b.Header.PreviousBlockHash != latestBlockHash {
		// Handle this case in the calling function such that the missing blocks are added to the blockchain
		return NewBlockPrevHashDontMatch
	}
	// Verify that the transactions in the block are not already present in the blockchain
	duplicate := VerifyDuplicateTransactionsInBlock(bchain, b)
	if duplicate {
		return NewBlockDuplicateTransactions
	}
	// Now finally adding new block to the blockchain
	// Stop the mining process
	nodestate.CloseMiningChan()
	// First prune the transactions from the transaction pool
	PruneTransactionList(b)
	// Add the block to the blockchain
	success := bchain.AddBlock(b)
	if !success {
		log.Println("\n\n\n\n\n\nSerious Error: Block was not added to the blockchain\n\n\n\n\n\n")
		os.Exit(1)
	}
	return NewBlockAddedSuccessfully
}
