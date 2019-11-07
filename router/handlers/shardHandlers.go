package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"net/http"
	"strconv"
)

func InitShardHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
	transactionPool = transaction.NewTransactionPool(SHARD_ID)
	sameShardPeers = peerList.NewPeerList(SHARD_ID)
}

func StartShardMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer(REGISTRATION_SERVER+"/register/")
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newPeers := peerList.NewPeerList(SHARD_ID)
			newPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for server, _ := range newPeers.Copy() {
				RegisterToServer(server+"/shard/peers/")
				sameShardPeers.Add(server)
			}
		}
	}
}

func RegisterShard(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	registerInfo := RegisterInfo{}
	err = json.Unmarshal(reqBody, &registerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	sameShardPeers.Add(registerInfo.Address)
	w.WriteHeader(http.StatusOK)
}

func GetPeerListForShard(w http.ResponseWriter, r *http.Request) {
	respBody, err := sameShardPeers.PeerMapToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respBody))
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	message := dataStructure.Message{}
	if message.Type != dataStructure.TRANSACTION {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	transactionPool.AddToTransactionPool(message.Transaction)
	go BroadcastTransaction(message)
}

func BroadcastTransaction(message dataStructure.Message) {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range sameShardPeers.Copy() {
				_, err := http.Post(k+"/shard/"+strconv.Itoa(int(SHARD_ID))+"/transaction/", "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					sameShardPeers.Delete(k)
				}
			}
		}
	}
}

func ShowAllTransactionsInPool(w http.ResponseWriter, r *http.Request) {
	transaction := transaction.NewTransaction("abc", "def", 45.5)
	transactionPool.AddToTransactionPool(transaction)
	fmt.Println(transactionPool.Show())
	w.Write([]byte(transactionPool.Show()))
	w.WriteHeader(http.StatusOK)
}
