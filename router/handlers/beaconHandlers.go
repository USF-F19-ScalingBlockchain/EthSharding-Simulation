package handlers

import (
	"bytes"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func TxReceive(w http.ResponseWriter, r *http.Request) {
	reqBody := readRequestBody(w, r)

	tx := transaction.JsonToTransaction(string(reqBody))
	shardId := transactionPool.GetShardId(tx.From)
	//todo : create a message from submitted transaction
	if sendPostReq(reqBody, shardId) { //send tx to shard miner
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func sendPostReq(reqBody []byte, shardId uint32) bool {
	client := http.Client{}

	//find a peer for the given shard
	if peers, ok := shardPeers[shardId]; ok { //if found
		peerAdd := peers.GetAPeer()
		url := peerAdd + "/shard/" + strconv.Itoa(int(shardId)) + "/transaction" //shard/{shardId}/transaction/
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Print("Cannot create PostReq, err " + err.Error())
		}
		client.Do(req) // sending tx message to "one" shard miner
		return true
	} else {
		return false
	}
}

func readRequestBody(w http.ResponseWriter, r *http.Request) []byte {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//return nil, errors.New("503: ServiceUnavailable")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	return reqBody
}
