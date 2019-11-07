package main

import (
	"github.com/EthSharding-Simulation/router"
	"github.com/EthSharding-Simulation/router/handlers"
	"github.com/EthSharding-Simulation/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	router := router.NewRouter()
	if len(os.Args) == 3 {
		// shard id miner
		port, _ := strconv.Atoi(os.Args[1])
		shardId, _ := strconv.Atoi(os.Args[2])
		handlers.InitShardHandler("http://localhost", int32(port), uint32(shardId))
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))

	} else if len(os.Args) == 2 {
		port, _ := strconv.Atoi(os.Args[1])
		handlers.InitBeaconHandler("http://localhost", int32(port), uint32(utils.BEACON_ID))
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))

	} else {
		log.Fatal(http.ListenAndServe(":6689", router))
	}
}
