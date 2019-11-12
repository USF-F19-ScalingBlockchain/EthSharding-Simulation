package shard

import (
	"encoding/hex"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"golang.org/x/crypto/sha3"
	"log"
	"strings"
	"time"
)

type Shard struct {
	ShardChainRoot     string
	Timestamp          time.Time
	Id                 string
	ProposerNode       string
	OpenTransactionSet map[transaction.Transaction]bool
}

func (s *Shard) Show() string {
	sb := strings.Builder{}
	sb.WriteString("\nShard:")
	sb.WriteString("\nshard id: " + s.Id +
		"\nshard Root: " + s.ShardChainRoot +
		"\nshard Proposer: " + s.ProposerNode +
		"\nshard Time :" + s.Timestamp.String() + "\n")
	sb.WriteString("Open Transaction Set:\n")
	for key, _ := range s.OpenTransactionSet {
		sb.WriteString(key.Show())
	}

	return sb.String()
}

func NewShard(ShardChainRoot string, Timestamp time.Time, ProposerNode string, OpenTransactionSet map[transaction.Transaction]bool) Shard {
	shard := Shard{
		ShardChainRoot:     ShardChainRoot,
		Timestamp:          Timestamp,
		ProposerNode:       ProposerNode,
		OpenTransactionSet: OpenTransactionSet,
	}
	shard.Id = shard.genId()
	return shard
}

func (shard *Shard) genId() string {
	str := shard.ShardChainRoot +
		shard.Id +
		shard.ProposerNode +
		mapToString(shard.OpenTransactionSet)
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func (shard *Shard) CreateShSig(id transaction.Identity) []byte {
	return id.GenSignature(shard.ShardToJsonByteArray())
}

func VerifyShSig(id transaction.Identity, shard Shard, sign []byte) bool {
	return transaction.VerifySingature(id.PublicKey, shard.ShardToJsonByteArray(), sign)
}

func (shard *Shard) ShardToJsonByteArray() []byte {
	txJson, err := json.Marshal(shard)
	if err != nil {
		log.Println("in ShardToJsonByteArray : Error in marshalling shard : ", err)
	}

	return txJson
}

func JsonToShard(shardJson string) Shard {
	shard := Shard{}
	err := json.Unmarshal([]byte(shardJson), &shard)
	if err != nil {
		log.Println("Error in unmarshalling Transaction, err - ", err)
		log.Println("String given to unmarshall Transaction, ================> \n ", shardJson, "\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n")
	}

	return shard
}

func mapToString(m map[transaction.Transaction]bool) string {
	s := ""
	for key, _ := range m {
		s += key.TransactionToJson()
	}
	return s
}
