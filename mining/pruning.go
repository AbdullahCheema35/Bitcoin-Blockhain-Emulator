package mining

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func PruneTransactionList(block types.Block) {
	transactionPool := nodestate.LockTransactionPool()
	defer func() {
		nodestate.UnlockTransactionPool(transactionPool)
	}()

	for _, transaction := range block.Body.Transactions.Transactions {
		transactionPool.RemoveTransaction(transaction)
	}
}
