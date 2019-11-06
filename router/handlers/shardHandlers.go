package handlers

import (
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/utils"
	"strconv"
)

var SHARD_ID uint32 = uint32(utils.BEACON_ID)
var REGISTRATION_SERVER = "loaclhost:6689"
var SELF_SERVER = "localhost:6689"

var sameShardPeers = peerList.NewPeerList(int32(SHARD_ID))
var transactionPool = transaction.NewTransactionPool(SHARD_ID)
var isStarted = false

func InitShardHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_SERVER = host + ":" + strconv.Itoa(int(port))
	transactionPool = transaction.NewTransactionPool(SHARD_ID)
	sameShardPeers = peerList.NewPeerList(int32(SHARD_ID))
}