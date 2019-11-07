package dataStructure

import (
	"crypto"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
)

type MessageType int

// Time ToDo: make separate message structures for each message.
const (
	TRANSACTION MessageType = 0
	BLOCK       MessageType = 1
)

type Message struct {
	Type        MessageType
	Transaction transaction.Transaction
	Block       blockchain.Block
	HopCount    int32
	Signature   string
	PublicKey   *crypto.PublicKey
}
