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
		"/register/peers/{shardId}",
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
		handlers.GetPeerListForShard, // Shard server
	},
	Route{
		"PostTransactionToShard",
		"POST",
		"/shard/{shardId}/transaction/", // Shard server
		handlers.AddTransaction,
	},
	Route{
		"GetTransactionsInPool",
		"GET",
		"/shard/{shardId}/transactionPool/", // Shard server
		handlers.ShowAllTransactionsInPool,
	},
}
