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
		"/register",
		handlers.Register, // Registration server
	},
	Route{
		"Register",
		"GET",
		"/register/peers/{shardId}",
		handlers.GetPeers, // Registration server
	},
	Route{
		"BeaconReceive",
		"POST",
		"/beacon/block/receive",
		handlers.BeaconReceive, // beacon server
	},


}
