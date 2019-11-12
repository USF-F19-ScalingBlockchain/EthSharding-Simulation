package transaction

import (
	"github.com/EthSharding-Simulation/dataStructure/shard"
	"sync"
)

//this is for mining
type ShardPool struct {
	pool map[string]shard.Shard `json:"pool"`
	//Confirmed map[string]bool        `json:"confirmed"`
	shardId uint32
	mux     sync.Mutex
}

//
//func NewTransactionPool(shardId uint32) TransactionPool {
//	return TransactionPool{
//		pool:    make(map[string]Transaction),
//		shardId: shardId,
//		//Confirmed: make(map[string]bool),
//	}
//}
//
//func (txp *TransactionPool) ContainsInTransactionPool(tx Transaction) bool {
//	if _, ok := txp.pool[tx.Id]; ok {
//		return true
//	}
//	return false
//}
//
//func (txp *TransactionPool) AddToTransactionPool(tx Transaction) { //duplicates in transactinon pool
//	txp.mux.Lock()
//	defer txp.mux.Unlock()
//
//	if _, ok := txp.pool[tx.Id]; !ok {
//		if txp.matchShard(tx) {
//			log.Println("In AddToTransactionPool : Adding new")
//			txp.pool[tx.Id] = tx
//		} else {
//			log.Println("Cannot add Transaction with Id: ", tx.Id, " in pool ", txp.shardId)
//		}
//	}
//}
//
//func (txp *TransactionPool) DeleteFromTransactionPool(txid string) {
//	txp.mux.Lock()
//	defer txp.mux.Unlock()
//
//	delete(txp.pool, txid)
//}
//
//func (txp *TransactionPool) Show() string {
//	var byteBuf bytes.Buffer
//
//	for _, tx := range txp.pool {
//		byteBuf.WriteString(tx.Show() + "\n")
//	}
//
//	return byteBuf.String()
//}
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
//func (txp *TransactionPool) matchShard(transaction Transaction) bool {
//	h := fnv.New32a()
//	h.Write([]byte(transaction.From))
//	return h.Sum32() % utils.TOTAL_SHARDS == txp.shardId
//}
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
