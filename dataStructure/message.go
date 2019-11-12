package dataStructure

import (
	"crypto"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/shard"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
)

type MessageType int

// Time ToDo: make separate message structures for each message.
const (
	TRANSACTION MessageType = 0
	BLOCK       MessageType = 1
	SHARD       MessageType = 2
)

type Message struct {
	Type        MessageType
	Transaction transaction.Transaction
	Block       blockchain.Block
	Shard       shard.Shard
	HopCount    int32
	Signature   string // signature of miner
	PublicKey   *crypto.PublicKey
}
