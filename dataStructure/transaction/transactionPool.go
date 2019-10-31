package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

//this is for mining
type TransactionPool struct {
	Pool map[string]Transaction `json:"pool"`
	//Confirmed map[string]bool        `json:"confirmed"`
	mux sync.Mutex
}

type TransactionPoolJson struct {
	Pool map[string]Transaction `json:"pool"`
}

func NewTransactionPool() TransactionPool {
	return TransactionPool{
		Pool: make(map[string]Transaction),
		//Confirmed: make(map[string]bool),
	}
}

func (txp *TransactionPool) AddToTransactionPool(tx Transaction) { //duplicates in transactinon pool
	txp.mux.Lock()
	defer txp.mux.Unlock()

	if _, ok := txp.Pool[tx.Id]; !ok {
		log.Println("In AddToTransactionPool : Adding new")
		txp.Pool[tx.Id] = tx
	}
}

func (txp *TransactionPool) AddPoolToTransactionPool(recvTxp TransactionPool) { //duplicates in transactinon pool
	txp.mux.Lock()
	defer txp.mux.Unlock()

	for _, tx := range recvTxp.Pool {
		if _, ok := txp.Pool[tx.Id]; !ok {
			log.Println("In AddToTransactionPool : Adding new - tx.Id : ", tx.Id)
			txp.Pool[tx.Id] = tx
		}
	}
}

func (txp *TransactionPool) DeleteFromTransactionPool(txid string) {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	delete(txp.Pool, txid)
}

func (txp *TransactionPool) Show() string {
	var byteBuf bytes.Buffer

	for _, tx := range txp.Pool {
		byteBuf.WriteString(tx.Show() + "\n")
	}

	return byteBuf.String()
}

func (txp *TransactionPool) ReadFromTransactionPool(n int) map[string]Transaction {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	tempMap := make(map[string]Transaction)
	counter := 0
	for txid, tx := range txp.Pool {

		if counter >= n || counter >= len(txp.Pool) {
			break
		}

		//txp.Pool[txid] = tx
		tempMap[txid] = tx
		counter++

		//txp.DeleteFromTransactionPool(txid)

	}
	return tempMap
}

//func (txp *TransactionPool) EncodeToJson() string {
//	jsonBytes, err := json.Marshal(txp)
//	if err != nil {
//		log.Println("Error in encoding TransactionPool to json, err - ", err)
//	}
//	log.Println("TransactionPool jsonStr is =======> ", string(jsonBytes))
//
//	return string(jsonBytes)
//}

func (txpj *TransactionPoolJson) EncodeToJsonTransactionPoolJson() string {
	jsonBytes, err := json.Marshal(txpj)
	if err != nil {
		log.Println("Error in encoding TransactionPool to json, err - ", err)
	}
	log.Println("TransactionPoolJson jsonStr is =======> ", string(jsonBytes))

	return string(jsonBytes)
}

//func DecodeJsonToTransactionPool(jsonStr string) TransactionPool {
//	txp := TransactionPool{}
//
//	err := json.Unmarshal([]byte(jsonStr), &txp)
//	if err != nil {
//		log.Println("Error in decoding json to TransactionPool, err - ", err)
//		log.Println("TransactionPool jsonStr is =======> ", jsonStr)
//	}
//	return txp
//}

func DecodeJsonToTransactionPoolJson(jsonStr string) TransactionPoolJson {
	txpj := TransactionPoolJson{}

	err := json.Unmarshal([]byte(jsonStr), &txpj)
	if err != nil {
		log.Println("Error in decoding json to TransactionPoolJson, err - ", err)
		log.Println("TransactionPoolJson jsonStr is =======> ", jsonStr)
	}
	return txpj
}

//Copy func returns a copy of the peerMap
func (txp *TransactionPool) GetTransactionPoolJsonObj() TransactionPoolJson {

	txp.mux.Lock()
	defer txp.mux.Unlock()

	txpj := TransactionPoolJson{}
	txpj.Pool = make(map[string]Transaction)
	//copyOfTxPool := make(map[string]Transaction)
	for k := range txp.Pool {
		txpj.Pool[k] = txp.Pool[k]
	}

	fmt.Println("GetTransactionPoolJsonObj :::::::::::::::: json is ", txpj.EncodeToJsonTransactionPoolJson())
	return txpj
}

//func TransactionPoolTest() {
//	txp := TransactionPool{}
//	str := txp.EncodeToJson()
//	fmt.Println("json Transaction pool - ",str )
//	txp1 := DecodeJsonToTransactionPool(str)
//	str1 := txp1.EncodeToJson()
//	fmt.Println("json Transaction pool - ",str1 )
//}
