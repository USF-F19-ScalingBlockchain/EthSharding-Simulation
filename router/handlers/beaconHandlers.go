package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func InitBeaconHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))

	//beaconPeers
	//shardPeers
	//shardPool = beacon.NewShardPool()
}

func StartBeaconMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer(REGISTRATION_SERVER + "/register/") // register itself to registration server
	//get all peers for beacon
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID))) // get all peers for 9999
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newBeaconPeers := peerList.NewPeerList(SHARD_ID)
			newBeaconPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for peer, _ := range newBeaconPeers.Copy() {
				fmt.Println("peer: ", peer)
				go RegisterToServer(peer + "/beacon/peers/") // announce it self to all its peers
				beaconPeers.Add(peer)
			}
		}
	}

	bpc := beaconPeers.Copy()
	if len(bpc) > 0 {
		for peer, _ := range bpc {
			fmt.Println("peer: ", peer)
			DownloadBeaconChain(peer)
		}
	}
	go GenerateBeaconBlocks() //todo : anurag from here : GenerateBeaconBlocks

	var sb strings.Builder
	sb.WriteString("::: Started BeaconMiner :::\n")
	sb.WriteString(getBeaconMinerPeers())

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(sb.String()))

}

func DownloadBeaconChain(peer string) {
	resp, err := http.Get(peer + "/beacon/upload/")
	if err == nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			sbc.UpdateEntireBlockChain(string(respBody))
		}
	}
}

func GetPeerListForBeacon(w http.ResponseWriter, r *http.Request) {
	peersJson := getBeaconMinerPeers()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peersJson))
}

func RegisterBeaconPeer(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	registerInfo := RegisterInfo{}
	err = json.Unmarshal(reqBody, &registerInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}

	beaconPeers.Add(registerInfo.Address)
	w.WriteHeader(http.StatusOK)

}

func TxReceive(w http.ResponseWriter, r *http.Request) {
	reqBody := readRequestBody(w, r)
	tx := transaction.JsonToTransaction(string(reqBody))
	shardId := transactionPool.GetShardId(tx.From)
	if sendTxPostReq(reqBody, shardId) { //send tx to shard miner
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ShowShardPool(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(shardPool.Show()))
	} else if r.Method == http.MethodPost {

	}
}

func RecvShardStuff(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	message := dataStructure.Message{}
	if message.Type != dataStructure.SHARD {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	shardPool.AddToShardPool(message.Shard)
	go BroadcastShard(message)
}

/// todo : here
func RecvBeaconBlock(w http.ResponseWriter, r *http.Request) {

}

func UploadBeaconChain(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := sbc.BlockChainToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockChainJson))
	}
}

func ShowBeaconChain(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sbc.Show()))
}

//////////////
//helper funcs

func BroadcastShard(message dataStructure.Message) {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1 // broadcasting
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range beaconPeers.Copy() {
				//_, err := http.Post(k+"/shard/transaction/", "application/json", bytes.NewBuffer(messageJson))
				_, err := http.Post(k+"/beacon/shard/", "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					beaconPeers.Delete(k)
				}
			}
		}
	}
}

func sendTxPostReq(reqBody []byte, shardId uint32) bool {
	client := http.Client{}

	sid := strconv.Itoa(int(shardId))

	//find a peer for the given shard
	if peerAdd, ok := shardPeersForBeacon[shardId]; ok { //if found
		url := peerAdd + "/shard/" + sid + "/transaction" //shard/{shardId}/transaction/
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			log.Print("Cannot create PostReq, err " + err.Error())
		}
		_, _ = client.Do(req) // sending tx message to "one" shard miner
		return true
	} else { //if peer not found
		resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + sid + "/") // /register/peers/{shardId}/
		if err != nil {
			log.Println(errors.New("cannot get valid response registration server"))
		}
		respBody := readResponseBody(resp)
		peers := peerList.JsonToPeerMap(respBody)
		log.Print("peers : " + string(respBody))

		if len(peers) > 0 {
			for k, _ := range peers {
				shardPeersForBeacon[shardId] = k
				url := k + "/shard/" + sid + "/transaction" //shard/{shardId}/transaction/
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
				if err != nil {
					log.Print("Cannot create PostReq, err " + err.Error())
				}
				_, err = client.Do(req) // sending tx message to "one" shard miner
				if err == nil {
					return true
				}
			}

		}
		return false
	} // end of else (if peer not found)

}

func readRequestBody(w http.ResponseWriter, r *http.Request) []byte {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//return nil, errors.New("503: ServiceUnavailable")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	return reqBody
}

func readResponseBody(r *http.Response) []byte {
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}
	}
	defer r.Body.Close()
	return respBody
}

func getBeaconMinerPeers() string {
	var sb strings.Builder
	sb.WriteString("Beacon Peers : \n")
	beaconPeerJson, err := beaconPeers.PeerMapToJson()
	if err != nil {
		sb.WriteString(err.Error())
	} else {
		sb.WriteString(beaconPeerJson)
	}

	sb.WriteString("\nShard Peers : \n")
	for shardId, shardPeer := range shardPeersForBeacon {
		sb.WriteString("Shard Id: " + string(shardId) + ", Shard Peer: " + shardPeer + "\n")
	}

	return sb.String()
}

func GenerateBeaconBlocks() {
	for {
		mpt, flag := shardPool.BuildMpt()
		if flag {
			latestBlocks := sbc.GetLatestBlocks()
			var latestBlock blockchain.Block
			if latestBlocks != nil && len(latestBlocks) != 0 {
				latestBlock = latestBlocks[rand.Intn(len(latestBlocks))]
			}
			for latestBlocks == nil || latestBlock.Header.Height == sbc.GetLatestBlocks()[0].Header.Height {
				parentHash := "genesis"
				if latestBlocks != nil {
					parentHash = latestBlock.Header.Hash
				}
				block := sbc.GenBlock(sbc.GetLength()+1, parentHash, mpt, "0", identity.PublicKey, blockchain.SHARD)
				transactionPool.DeleteTransactions(block.Value)
				sbc.Insert(block)
				message := dataStructure.Message{
					Type:        dataStructure.BLOCK,
					Transaction: transaction.Transaction{},
					Block:       block,
					HopCount:    1,
					NodeId:      SELF_ADDR,
				}
				message.Sign(identity)
				latestBlocks = sbc.GetLatestBlocks()
				fmt.Println(latestBlock.Header.Height)
				go Broadcast(message, "/beacon/block/")
			}
		}
	}
}
