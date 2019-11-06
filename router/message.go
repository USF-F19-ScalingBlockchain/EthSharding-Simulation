package router

import (
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
)

type MessageType int

const(
	TRANSACTION MessageType = 0
	BLOCK 	MessageType = 1
)

type Message struct {
	Type		MessageType
	Transaction transaction.Transaction
	Block 		blockchain.Block
	HopCount 	int32
}
