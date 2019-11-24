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
		//"http://localhost:8002/start/beacon/",
		//"http://localhost:8003/start/beacon/",
		//"http://localhost:8004/start/beacon/",
		//"http://localhost:8005/start/beacon/",
		//"http://localhost:8006/start/beacon/",
		//"http://localhost:8007/start/beacon/",
		//"http://localhost:8008/start/beacon/",
		//"http://localhost:8009/start/beacon/",
		//"http://localhost:8010/start/beacon/",
		//"http://localhost:8011/start/beacon/",
		//"http://localhost:8012/start/beacon/",
		//"http://localhost:8013/start/beacon/",
		//"http://localhost:8014/start/beacon/",
		//"http://localhost:8015/start/beacon/",
		//"http://localhost:8016/start/beacon/",
		"http://localhost:8001/start/shard/",
		//"http://localhost:8002/start/shard/",
		//"http://localhost:8003/start/shard/",
		//"http://localhost:8004/start/shard/",
		//"http://localhost:8005/start/shard/",
		//"http://localhost:8006/start/shard/",
		//"http://localhost:8007/start/shard/",
		//"http://localhost:8008/start/shard/",
		//"http://localhost:8009/start/shard/",
		//"http://localhost:8010/start/shard/",
		//"http://localhost:8011/start/shard/",
		//"http://localhost:8012/start/shard/",
		//"http://localhost:8013/start/shard/",
		//"http://localhost:8014/start/shard/",
		//"http://localhost:8015/start/shard/",
		//"http://localhost:8016/start/shard/",
	}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(url + resp.Status)
	}
}
