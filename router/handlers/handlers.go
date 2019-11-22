package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/shard"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/utils"
	"net/http"
	"sync"
	"time"
)

var REGISTRATION_SERVER = "http://localhost:6689"
var SELF_ADDR = "http://localhost:6689"

// start shard server
var SHARD_ID = utils.BEACON_ID
var sameShardPeers = peerList.NewPeerList(SHARD_ID)
var transactionPool = transaction.NewTransactionPool(SHARD_ID)
var sbc blockchain.SyncBlockChain       //for both shard
var beaconSbc blockchain.SyncBlockChain // and beacon chain

// end shard sever

// start registration server
var beaconPeers = peerList.NewPeerList(utils.BEACON_ID) // also used by beacon miner
var shardPeers = map[uint32]peerList.PeerList{}
var identity = transaction.NewIdentity()
var recvLock = sync.Mutex{}
var recvTime = map[string]time.Time{}
var finalizeLock = sync.Mutex{}
var finalizeTime = map[string]time.Time{}
var openTransactionSet transaction.OpenTransactionSet
// end registration server

// start Beacon server
var shardPeersForBeacon = map[uint32]string{} // each shard one peer for beacon server
var shardPool = shard.NewShardPool()

var tkShardRecv = dataStructure.NewTimeMap()     //t2
var tkShardsInBlock = dataStructure.NewTimeMap() //t3
var dkt3t2 = dataStructure.NewDurationMap()      //t3-t2

// end of variables for beacon server

// functions
func RegisterToServer(url string, selfId uint32) {
	registerInfo := RegisterInfo{SELF_ADDR, selfId}
	registerInfoJson, err := json.Marshal(registerInfo)
	if err == nil {
		http.Post(url, "application/json", bytes.NewBuffer([]byte(registerInfoJson)))
	}
}
