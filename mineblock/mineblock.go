package mineblock

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/validation"
)

func MineBlock(b types.Block, miningProcessAbortedChan chan bool) {
	// fmt.Println("Inside the function MineBlock")
	// fmt.Println("Mining block...")

	difficultyTarget := b.Header.DifficultyTarget
	_, targetHashBytes := validation.GenerateTargetHash(difficultyTarget)

	// // fmt.Println("Difficulty target:", difficultyTarget)
	// // fmt.Println("Target hash:", targetHash)

	header := b.Header
	header.Nonce = 0

	for {
		_, hashBytes := header.RecalculateBlockHeaderHash()
		if validation.CompareWithTargetHash(hashBytes, targetHashBytes) {
			// fmt.Printf("Block mined! Nonce: %d\n", header.Nonce)
			break
		}

		select {
		case <-miningProcessAbortedChan:
			// fmt.Println("Mining process aborted")
			return
		default:
			header.Nonce++
		}
		// // Sleep for random amount of time less than 100 milliseconds
		// time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}

	b.Header = header
	b.RecalculateBlockHash()

	// Get self node
	selfNode := configuration.GetSelfServerAddress()

	// Call the function to handle the mined block
	propagation.HandleReceivedBlock(b, selfNode)
}
