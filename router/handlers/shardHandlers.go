package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/router"
	"github.com/EthSharding-Simulation/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

var SHARD_ID = uint32(utils.BEACON_ID)
var REGISTRATION_SERVER = "http://localhost:6689"
var SELF_ADDR = "http://localhost:6689"

var sameShardPeers = peerList.NewPeerList(SHARD_ID)
var transactionPool = transaction.NewTransactionPool(SHARD_ID)

func InitShardHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
	transactionPool = transaction.NewTransactionPool(SHARD_ID)
	sameShardPeers = peerList.NewPeerList(SHARD_ID)
}

func StartShardMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer()
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newPeers := peerList.NewPeerList(SHARD_ID)
			newPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for k, _ := range newPeers.Copy() {
				RegisterToPeers(k)
				sameShardPeers.Add(k)
			}
		}
	}
}

func RegisterToServer() {
	registerInfo := RegisterInfo{SELF_ADDR, SHARD_ID}
	registerInfoJson, err := json.Marshal(registerInfo)
	if err == nil {
		http.Post(REGISTRATION_SERVER + "/register/", "application/json", bytes.NewBuffer([]byte(registerInfoJson)))
	}
}

func RegisterToPeers(server string) {
	registerInfo := RegisterInfo{SELF_ADDR, SHARD_ID}
	registerInfoJson, err := json.Marshal(registerInfo)
	if err == nil {
		http.Post(server + "/shard/peers/", "application/json", bytes.NewBuffer([]byte(registerInfoJson)))
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
	message := router.Message{}
	if message.Type != router.TRANSACTION {
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

func BroadcastTransaction(message router.Message) {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range sameShardPeers.Copy() {
				_, err := http.Post(k + "/shard/" + strconv.Itoa(int(SHARD_ID)) + "/transaction/", "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					sameShardPeers.Delete(k)
				}
			}
		}
	}
}

func ShowAllTransactionsInPool(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(transactionPool.Show()))
}
