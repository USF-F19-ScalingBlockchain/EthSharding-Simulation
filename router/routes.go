package router

import (
	"github.com/EthSharding-Simulation/router/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Register",
		"POST",
		"/register/",
		handlers.Register, // Registration server
	},
	Route{
		"GetPeers",
		"GET",
		"/register/peers/{shardId}/",
		handlers.GetPeers, // Registration server
	},
	Route{
		"RegisterShard",
		"POST",
		"/shard/peers/",
		handlers.RegisterShard, // Shard server
	},
	Route{
		"startShardMiner",
		"GET",
		"/start/shard/",
		handlers.StartShardMiner, // Shard server
	},
	Route{
		"GetPeerListForShard",
		"GET",
		"/shard/peers/",
		handlers.GetPeerListForShard, // Shard server - for UI
	},
	Route{
		"PostTransactionToShard",
		"POST",
		"/shard/transaction/", // Shard server
		handlers.AddTransaction,
	},
	Route{
		"GetTransactionsInPool",
		"GET",
		"/shard/transactionPool/", // Shard server
		handlers.ShowAllTransactionsInPool,
	},
	Route{
		"TxReceive",
		"POST",
		"/beacon/Tx/receive/",
		handlers.TxReceive, // beacon server
	},
	Route{
		"StartBeaconMiner",
		"GET",
		"/start/beacon/",
		handlers.StartBeaconMiner, // beacon server
	},
	Route{
		"RegisterBeaconPeer",
		"POST",
		"/beacon/peers/",
		handlers.RegisterBeaconPeer, // beacon server
	},
	Route{
		"GetBeaconPeers",
		"GET",
		"/beacon/peers/",
		handlers.GetPeerListForBeacon, // beacon server
	},
	Route{
		"GetShardChain",
		"GET",
		"/shard/upload/",
		handlers.UploadBlockchain, // Shard server
	},
	Route{
		"PostBlockToShard",
		"POST",
		"/shard/block/",
		handlers.AddShardBlock, // Shard server
	},
	Route{
		"ShowShard",
		"GET",
		"/shard/show/",
		handlers.ShowShard, // Shard server
	},
	Route{
		"GetShardBlockAtHeight",
		"GET",
		"/shard/block/{height}/{hash}",
		handlers.UploadShardBlock,
	},
	Route{
		"GetBeaconChain",
		"GET",
		"/beacon/upload/",
		handlers.UploadBeaconChain, // Beacon server
	},
	Route{
		"ShowShardPool",
		"GET",
		"/beacon/shardpool/",
		handlers.ShowShardPool, // Beacon server
	},
	Route{
		"RecvShard",
		"POST",
		"/beacon/shard/",
		handlers.RecvShardStuff, // Beacon server // will recv from shard miner
	},
	Route{
		"RecvBeaconBlock",
		"POST",
		"/beacon/block/",
		handlers.RecvBeaconBlock, // Beacon server // will recv from beacon miner
	},
	Route{
		"ShowBeaconChain",
		"GET",
		"/beacon/show/",
		handlers.ShowBeaconChain, // Beacon server
	},
	Route{
		"GetBeaconBlock",
		"GET",
		"/shard/beacon/",
		handlers.GetBeaconBlock, // Shard server
	},
	//Route{
	//	"ShowOpenTransactionSet",
	//	"GET",
	//	"/shard/openTransactionSet/",
	//	handlers.ShowOpenTransactionSet,
	//},
	Route{
		"GetFinalTimePerTransaction",
		"GET",
		"/shard/finalityTime",
		handlers.GetFinalTimePerTransaction,
	},
	Route{
		"ShowReceivedTimes",
		"GET",
		"/shard/recvTimes",
		handlers.ShowReceivedTimes,
	},
	Route{
		"ShowFinalizeTimes",
		"GET",
		"/shard/finalTimes",
		handlers.ShowFinalizeTimes,
	},
	Route{
		"GetBeaconBlockAtHeight",
		"GET",
		"/beacon/block/{height}/{hash}",
		handlers.UploadBeaconBlock, //beacon server
	},
	Route{
		"Showt2",
		"GET",
		"/beacon/showt2",
		handlers.Showt2, //beacon server
	},
	Route{
		"Showt3",
		"GET",
		"/beacon/showt3",
		handlers.Showt3, //beacon server
	},
	Route{
		"Showt3t2",
		"GET",
		"/beacon/showt3t2",
		handlers.Showt3t2, //beacon server
	},
}
