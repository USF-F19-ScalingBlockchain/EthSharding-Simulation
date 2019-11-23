package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	urls := []string{ //8000 //8001 0 // 8002 1 // 8003 2 // 8004 3
		"http://localhost:8000/start/beacon/",
		"http://localhost:8001/start/beacon/",
		"http://localhost:8002/start/beacon/",
		//"http://localhost:8003/start/beacon/",
		//"http://localhost:8004/start/beacon/",
		"http://localhost:8001/start/shard/",
		"http://localhost:8002/start/shard/",
		//"http://localhost:8003/start/shard/",
		//"http://localhost:8004/start/shard/",
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(url + resp.Status)
	}
}
