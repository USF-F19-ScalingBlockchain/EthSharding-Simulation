package dataStructure

import (
	"crypto/rsa"
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
	Shard       shard.Shard //
	HopCount    int32
	Signature   string // signature of miner
	PublicKey   *rsa.PublicKey
	NodeId		string
}

func (message *Message) Sign(identity transaction.Identity) {
	if message.Type == TRANSACTION {
		message.Signature = string(message.Transaction.CreateTxSig(identity))
		message.PublicKey = identity.PublicKey
	} else if message.Type == BLOCK {
		message.Signature = string(message.Block.CreateBlockSig(identity))
		message.PublicKey = identity.PublicKey
	} else if message.Type == SHARD {
		message.Signature = string(message.Shard.CreateShSig(identity))
		message.PublicKey = identity.PublicKey
	}
}

func (message *Message) Verify() bool {
	if message.Type == TRANSACTION {
		return transaction.VerifyTxSig(message.PublicKey, message.Transaction, []byte(message.Signature))
	} else if message.Type == BLOCK {
		return blockchain.VerifyBlockSig(message.PublicKey, message.Block, message.Signature)
	} else if message.Type == SHARD {
		return shard.VerifyShSig(message.PublicKey, message.Shard, []byte(message.Signature))
	}
	return false
}
