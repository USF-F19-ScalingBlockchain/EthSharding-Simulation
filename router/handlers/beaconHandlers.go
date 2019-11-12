package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func InitBeaconHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))

	//beaconPeers
	//shardPeers
}

func StartBeaconMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer(REGISTRATION_SERVER + "/register/") // register itself to registration server
	//todo : get all peers for beacon
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID))) // get all peers for 9999
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newBeaconPeers := peerList.NewPeerList(SHARD_ID)
			newBeaconPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for peer, _ := range newBeaconPeers.Copy() {
				fmt.Println("peer: ", peer)
				go RegisterToServer(peer + "/beacon/peers/")
				beaconPeers.Add(peer)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("::: Started BeaconMiner :::\n")
	sb.WriteString(getBeaconMinerPeers())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(sb.String()))

}

func GetPeerListForBeacon(w http.ResponseWriter, r *http.Request) {
	peersJson := getBeaconMinerPeers()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peersJson))
}

func RegisterBeaconPeer(w http.ResponseWriter, r *http.Request) {
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

	beaconPeers.Add(registerInfo.Address)
	w.WriteHeader(http.StatusOK)

}

func TxReceive(w http.ResponseWriter, r *http.Request) {
	reqBody := readRequestBody(w, r)
	tx := transaction.JsonToTransaction(string(reqBody))
	shardId := transactionPool.GetShardId(tx.From)
	//todo : create a message from submitted transaction
	if sendTxPostReq(reqBody, shardId) { //send tx to shard miner
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func sendTxPostReq(reqBody []byte, shardId uint32) bool {
	client := http.Client{}

	//find a peer for the given shard
	if peerAdd, ok := shardPeersForBeacon[shardId]; ok { //if found
		url := peerAdd + "/shard/" + strconv.Itoa(int(shardId)) + "/transaction" //shard/{shardId}/transaction/
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Print("Cannot create PostReq, err " + err.Error())
		}
		client.Do(req) // sending tx message to "one" shard miner
		return true
	} else { //if peer not found
		//todo :get all peers for that shardId
		// if good response
		return false
	}
}

//helper funcs
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

func getBeaconMinerPeers() string {
	var sb strings.Builder
	sb.WriteString("Beacon Peers : \n")
	beaconPeerJson, err := beaconPeers.PeerMapToJson()
	if err != nil {
		sb.WriteString(err.Error())
	} else {
		sb.WriteString(beaconPeerJson)
	}

	sb.WriteString("\nShard Peers : \n")
	for shardId, shardPeer := range shardPeersForBeacon {
		sb.WriteString("Shard Id: " + string(shardId) + ", Shard Peer: " + shardPeer + "\n")
	}

	return sb.String()
}
