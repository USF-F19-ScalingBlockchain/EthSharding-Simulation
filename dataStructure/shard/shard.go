package shard

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"golang.org/x/crypto/sha3"
	"log"
	"strings"
	"time"
)

type Shard struct {
	ShardChainRoot     string                             `json:"shardChainRoot"`
	Timestamp          time.Time                          `json:"timestamp"`
	Id                 string                             `json:"id"`
	ProposerNode       string                             `json:"proposerNode"`
	OpenTransactionSet map[string]transaction.Transaction `json:openTransactionSet`
}

func (s *Shard) Show() string {
	sb := strings.Builder{}
	sb.WriteString("\nShard:")
	sb.WriteString("\nshard id: " + s.Id +
		"\nshard Root: " + s.ShardChainRoot +
		"\nshard Proposer: " + s.ProposerNode +
		"\nshard Time :" + s.Timestamp.String() + "\n")
	sb.WriteString("Open Transaction Set:\n")
	for _, val := range s.OpenTransactionSet {
		sb.WriteString(string(val.TransactionToJsonByteArray()))
	}

	return sb.String()
}

func NewShard(ShardChainRoot string, Timestamp time.Time, ProposerNode string) Shard {
	shard := Shard{
		ShardChainRoot:     ShardChainRoot,
		Timestamp:          Timestamp,
		ProposerNode:       ProposerNode,
		OpenTransactionSet: make(map[string]transaction.Transaction),
	}
	shard.Id = shard.genId()
	return shard
}

func (shard *Shard) genId() string {
	str := shard.Show()
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func (shard *Shard) CreateShSig(id transaction.Identity) []byte {
	return id.GenSignature(shard.ShardToJsonByteArray())
}

func VerifyShSig(id *rsa.PublicKey, shard Shard, sign []byte) bool {
	return transaction.VerifySingature(id, shard.ShardToJsonByteArray(), sign)
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
