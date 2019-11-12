package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/utils"
	"net/http"
)

var REGISTRATION_SERVER = "http://localhost:6689"
var SELF_ADDR = "http://localhost:6689"

// start shard server
var SHARD_ID = uint32(utils.BEACON_ID)
var sameShardPeers = peerList.NewPeerList(SHARD_ID)
var transactionPool = transaction.NewTransactionPool(SHARD_ID)
var sbc blockchain.SyncBlockChain
// end shard sever

// start registration server
var beaconPeers = peerList.NewPeerList(uint32(utils.BEACON_ID)) // also used by beacon miner
var shardPeers = map[uint32]peerList.PeerList{}                 // also used by beacon miner
var identity transaction.Identity
// end registration server

func RegisterToServer(url string) {
	registerInfo := RegisterInfo{SELF_ADDR, SHARD_ID}
	registerInfoJson, err := json.Marshal(registerInfo)
	if err == nil {
		http.Post(url, "application/json", bytes.NewBuffer([]byte(registerInfoJson)))
	}
}
