package peerList

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
)

type PeerList struct {
	shardId   int32
	peerMap   map[string]bool // <Host:Port, shardId>
	mux       sync.Mutex
}

/**
Creates new PeerList
*/
func NewPeerList(id int32) PeerList {
	return PeerList{id, map[string]bool{}, sync.Mutex{}}
}

/**
Add peer to peer map. Takes address and id of peer as input.
*/
func (peers *PeerList) Add(addr string) {
	peers.mux.Lock()
	peers.peerMap[addr] = true
	peers.mux.Unlock()
}

/**
Deletes a peer from peerList. Takes address to identify which peer to delete.
*/
func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}

/**
It converts the peermap into readable format.
*/
func (peers *PeerList) Show() string {
	rs := "ShardId: " + strconv.Itoa(int(peers.shardId)) + "\n"
	var keyList []string
	for key := range peers.peerMap {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	rs += "\nPeerList: \n[ \n"
	for _, key := range keyList {
		rs += key + "\n"
	}
	rs += "]\n"
	return rs
}

/**
Sets the self id in peerList.
*/
func (peers *PeerList) Register(id int32) {
	peers.shardId = id
	fmt.Printf("SelfId=%v\n", id)
}

/**
Copies the peermap to new buffer and returns
the current picture of peermap.
*/
func (peers *PeerList) Copy() map[string]bool {
	newMap := make(map[string]bool);
	for k, v := range peers.peerMap {
		newMap[k] = v
	}
	return newMap
}

/**
It will get self id from PeerList.
*/
func (peers *PeerList) GetShardId() int32 {
	return peers.shardId
}

/**
Converts peerMap to JSON String.
*/
func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	s, err := json.Marshal(peers.peerMap)
	return string(s), err
}

/**
Converts peerMapJson string to peerMap and add
peers to peerMap.
*/
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	newPeerMap := map[string]bool{}
	err := json.Unmarshal([]byte(peerMapJsonStr), &newPeerMap)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range newPeerMap {
		if k != selfAddr {
			peers.peerMap[k] = v
		}
	}
}
