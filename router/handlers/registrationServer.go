package handlers

import (
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type RegisterInfo struct {
	Address string `json:"address"`
	ShardId uint32 `json:"shardId"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	registerInfo := RegisterInfo{}
	err = json.Unmarshal(reqBody, &registerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if registerInfo.ShardId == 9999 {
		beaconPeers.Add(registerInfo.Address)
	} else {
		if _, ok := shardPeers[registerInfo.ShardId]; !ok {
			shardPeers[registerInfo.ShardId] = peerList.NewPeerList(registerInfo.ShardId)
		}
		sp := shardPeers[registerInfo.ShardId]
		sp.Add(registerInfo.Address)
	}
	w.WriteHeader(http.StatusOK)
}

func GetPeers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shardId, err := strconv.Atoi(vars["shardId"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if uint32(shardId) == utils.BEACON_ID {
		beaconPeersJson, err := beaconPeers.PeerMapToJson()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(beaconPeersJson))
	} else {
		if val, ok := shardPeers[uint32(shardId)]; ok {
			shardPeersJson, err := val.PeerMapToJson()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(shardPeersJson))
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
