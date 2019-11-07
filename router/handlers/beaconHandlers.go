package handlers

import (
	"bytes"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func InitBeaconHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
	//beaconPeers
	//shardPeers
}

func StartBeaconMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer()
	//todo : get all peers for beacon
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newBeaconPeers := peerList.NewPeerList(SHARD_ID)
			newBeaconPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for k, _ := range newBeaconPeers.Copy() {
				RegisterToPeers(k)
				beaconPeers.Add(k)
			}
		}
	}

}

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
		//todo :get all peers for that shardId
		// if good response
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
