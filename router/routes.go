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
		"/shard/{shardId}/transactionPool/", // Shard server
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
		handlers.GetBlock,
	},
	//Route{
	//	"GetBeaconChain",
	//	"GET",
	//	"/beacon/upload/",
	//	handlers.UploadBeaconChain, // Beacon server todo
	//},
	Route{
		"ShowShardPool",
		"GET",
		"/beacon/shardpool/",
		handlers.ShowShardPool, // Shard server
	},
	Route{
		"AddToShardPool",
		"POST",
		"/beacon/shardpool/",
		handlers.AddToShardPool, // Shard server todo
	},
}
