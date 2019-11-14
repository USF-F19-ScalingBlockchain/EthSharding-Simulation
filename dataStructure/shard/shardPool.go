package shard

import (
	"bytes"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/mpt"
	"github.com/EthSharding-Simulation/utils"
	"log"
	"sync"
)

//this is for mining
type ShardPool struct {
	pool map[string]Shard `json:"pool"`
	//Confirmed map[string]bool        `json:"confirmed"`
	shardId uint32
	mux     sync.Mutex
}

func NewShardPool() ShardPool {
	return ShardPool{
		pool:    make(map[string]Shard),
		shardId: utils.BEACON_ID,
		//Confirmed: make(map[string]bool),
	}
}

func (sp *ShardPool) ContainsInShardPool(s Shard) bool {
	if _, ok := sp.pool[s.Id]; ok {
		return true
	}
	return false
}

func (sp *ShardPool) AddToShardPool(s Shard) {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	if _, ok := sp.pool[s.Id]; !ok {
		log.Println("In AddToTransactionPool : Adding new")
		sp.pool[s.Id] = s
	}
}

func (sp *ShardPool) DeleteFromShardPool(shardId string) {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	delete(sp.pool, shardId)
}

func (sp *ShardPool) DeleteShards(mpt mpt.MerklePatriciaTrie) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	for key, _ := range sp.pool {
		if _, ok := mpt.Raw_db[key]; ok {
			delete(sp.pool, key)
		}
	}
}

func (sp *ShardPool) Show() string {
	var byteBuf bytes.Buffer

	for _, s := range sp.pool {
		byteBuf.WriteString(s.Show() + "\n")
	}

	return byteBuf.String()
}

//
//// ToDo: Should we delete shard from shard pool after building mpt?
func (shardPool *ShardPool) BuildMpt() (mpt.MerklePatriciaTrie, bool) {
	shardPool.mux.Lock()
	shardPool.mux.Unlock()
	shardMpt := mpt.MerklePatriciaTrie{}
	shardMpt.Initial()
	if len(shardPool.pool) < utils.MIN_TX_POOL_SIZE {
		return shardMpt, false
	}
	for i, _ := range shardPool.pool {
		shardsJson, err := json.Marshal(shardPool.pool[i])
		if err == nil {
			shardMpt.Insert(shardPool.pool[i].Id, string(shardsJson))
		}
	}
	return shardMpt, true
}

//old
//
//func (txp *TransactionPool) ReadFromTransactionPool(n int) map[string]Transaction {
//	txp.mux.Lock()
//	defer txp.mux.Unlock()
//
//	tempMap := make(map[string]Transaction)
//	counter := 0
//	for txid, tx := range txp.pool {
//
//		if counter >= n || counter >= len(txp.pool) {
//			break
//		}
//
//		//txp.Pool[txid] = tx
//		tempMap[txid] = tx
//		counter++
//
//		//txp.DeleteFromTransactionPool(txid)
//
//	}
//	return tempMap
//}
//
//func (txp *TransactionPool) GetShardId(fromField string) uint32 {
//	h := fnv.New32a()
//	h.Write([]byte(fromField))
//	return h.Sum32() % utils.TOTAL_SHARDS
//}
//
////func (txp *TransactionPool) matchShard(transaction Transaction) bool {
////	h := fnv.New32a()
////	h.Write([]byte(transaction.From))
////	return h.Sum32() % utils.TOTAL_SHARDS == txp.shardId
////}
//
//func (txp *TransactionPool) IsOpenTransaction(transaction Transaction) bool {
//	h := fnv.New32a()
//	h.Write([]byte(transaction.To))
//	return h.Sum32()%utils.TOTAL_SHARDS != txp.shardId
//}
//
//func (txp *TransactionPool) BuildMpt() (mpt.MerklePatriciaTrie, bool) {
//	txp.mux.Lock()
//	txp.mux.Unlock()
//	txMpt := mpt.MerklePatriciaTrie{}
//	txMpt.Initial()
//	if len(txp.pool) < utils.MIN_TX_POOL_SIZE {
//		return txMpt, false
//	}
//	for i, _ := range txp.pool {
//		transJson, err := json.Marshal(txp.pool[i])
//		if err == nil {
//			txMpt.Insert(i, string(transJson))
//			delete(txp.pool, i)
//		}
//	}
//	return txMpt, true
//}
//
////func (txp *TransactionPoolJson) EncodeToJsonTransactionPoolJson() string {
////	jsonBytes, err := json.Marshal(txp)
////	if err != nil {
////		log.Println("Error in encoding TransactionPool to json, err - ", err)
////	}
////	log.Println("TransactionPoolJson jsonStr is =======> ", string(jsonBytes))
////
////	return string(jsonBytes)
////}
////
////func DecodeJsonToTransactionPoolJson(jsonStr string) TransactionPoolJson {
////	txp := TransactionPoolJson{}
////
////	err := json.Unmarshal([]byte(jsonStr), &txp)
////	if err != nil {
////		log.Println("Error in decoding json to TransactionPoolJson, err - ", err)
////		log.Println("TransactionPoolJson jsonStr is =======> ", jsonStr)
////	}
////	return txp
////}
////
//////Copy func returns a copy of the peerMap
////func (txp *TransactionPool) GetTransactionPoolJsonObj() TransactionPoolJson {
////
////	txp.mux.Lock()
////	defer txp.mux.Unlock()
////
////	txpj := TransactionPoolJson{}
////	txpj.Pool = make(map[string]Transaction)
////	//copyOfTxPool := make(map[string]Transaction)
////	for k := range txp.Pool {
////		txpj.Pool[k] = txp.Pool[k]
////	}
////
////	fmt.Println("GetTransactionPoolJsonObj :::::::::::::::: json is ", txpj.EncodeToJsonTransactionPoolJson())
////	return txpj
////}
//
