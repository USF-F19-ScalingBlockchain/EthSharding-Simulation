package main

import (
	"github.com/EthSharding-Simulation/router"
	"log"
	"net/http"
	"os"
)

func main() {
	router := router.NewRouter()
	if len(os.Args) == 3 {
		// shard id miner
		// initShardMiner()
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else if len(os.Args) == 2 {
		// initBeaconMiner()
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6689", router))
	}
}
