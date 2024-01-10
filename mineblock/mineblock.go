package mineblock

import (
	"log"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/validation"
)

func MineBlock(b types.Block, miningProcessAbortedChan chan bool) {
	// log.Println("Mining block...")

	difficultyTarget := b.Header.DifficultyTarget
	_, targetHashBytes := validation.GenerateTargetHash(difficultyTarget)

	header := b.Header
	header.Nonce = 0

	for {
		_, hashBytes := header.RecalculateBlockHeaderHash()
		if validation.CompareWithTargetHash(hashBytes, targetHashBytes) {
			// log.Printf("Block mined! Nonce: %d\n", header.Nonce)
			break
		}

		select {
		case <-miningProcessAbortedChan:
			log.Println("Mining process aborted")
			return
		default:
			header.Nonce++
		}
	}

	b.Header = header
	b.RecalculateBlockHash()

	selfNode := configuration.GetSelfServerAddress()

	propagation.HandleReceivedBlock(b, selfNode)
	// log.Println("Mining process completed (successfully)")
}
