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
	"github.com/EthSharding-Simulation/utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InitBeaconHandler(host string, port int32, shardId uint32) {
	//SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
	//identity = transaction.NewIdentity()
	//beaconPeers
	//shardPeers
	//shardPool = beacon.NewShardPool()
}

func StartBeaconMiner(w http.ResponseWriter, r *http.Request) {
	beaconSbc = blockchain.NewBlockChain()
	RegisterToServer(REGISTRATION_SERVER+"/register/", utils.BEACON_ID) // register itself to registration server
	//get all peers for beacon
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(utils.BEACON_ID))) // get all peers for 9999
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newBeaconPeers := peerList.NewPeerList(utils.BEACON_ID)
			newBeaconPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			for peer, _ := range newBeaconPeers.Copy() {
				fmt.Println("peer: ", peer)
				go RegisterToServer(peer+"/beacon/peers/", utils.BEACON_ID) // announce it self to all its peers
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
	go GenerateBeaconBlocks() //GenerateBeaconBlocks

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
			beaconSbc.UpdateEntireBlockChain(string(respBody))
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
	//making message of tx
	msg := dataStructure.Message{
		Type:        dataStructure.TRANSACTION,
		Transaction: tx,
		HopCount:    1,
		NodeId:      SELF_ADDR,
		TimeStamp:   time.Now(),
	}
	msg.Sign(identity)
	msgJson, err := json.Marshal(&msg)
	if err != nil {
		log.Println("cannot convert message to json")
	}
	// end of making message of tx

	if sendTxPostReq(msgJson, shardId) { //send tx to shard miner
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
	shardRecvTime := time.Now() //t2
	fmt.Println("Recved shard from shard miner")
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()

	message := dataStructure.Message{}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if message.Type != dataStructure.SHARD {
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println("yaha aa rahe hai!")
	shardPool.AddToShardPool(message.Shard)
	for _, j := range message.Shard.OpenTransactionSet {
		fmt.Println("Open Tx Set: " + j.Show())
	}
	//put t2 in tkShardRecv
	tkShardRecv.AddTxIdsFromShardToKeeper(message.Shard.OpenTransactionSet, shardRecvTime, message.Shard.TxProcessingTime) //t2

	go BroadcastMessage(message, "/beacon/shard/", beaconPeers.Copy()) // BroadcastShardMessage

}

func RecvBeaconBlock(w http.ResponseWriter, r *http.Request) {
	shardInBlockTime := time.Now() //t3

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	message := dataStructure.Message{}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if message.Type != dataStructure.BLOCK {
		w.WriteHeader(http.StatusBadRequest)
	}
	shardPool.DeleteShards(message.Block.Value)
	//check for parent hash
	if !beaconSbc.CheckParentHash(message.Block) && message.Block.Header.Height-1 > 0 {
		if AskForBeaconBlock(message.Block.Header.Height-1, message.Block.Header.ParentHash) {
			beaconSbc.Insert(message.Block)
		}
	} else {
		beaconSbc.Insert(message.Block)
	}

	go BroadcastMessage(message, "/beacon/block/", beaconPeers.Copy()) // BroadcastBeaconBlockMessage
	message.HopCount = 1
	message.Sign(identity)
	go BroadcastMessageToShardMiners(message, "/shard/beacon/", sameShardPeers.Copy()) //acting as shard miner

	//add t3 to tkShardsInBlock
	tkShardsInBlock.AddTxIdsFromBeaconBlockToKeeper(message.Block.Value.Raw_db, shardInBlockTime) //t3
}

func UploadBeaconChain(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := beaconSbc.BlockChainToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockChainJson))
	}
}

func UploadBeaconBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, _ := strconv.Atoi(vars["height"])
	hash := vars["hash"]

	blk, found := beaconSbc.GetBlock(int32(height), hash)
	if !found {
		w.WriteHeader(http.StatusNoContent)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(blk.EncodeToJSON()))
}

func ShowBeaconChain(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(beaconSbc.Show()))
}

func Showt2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tkShardRecv.ToString()))
}

func Showt3(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tkShardsInBlock.ToString()))
}

func Showt3t2(w http.ResponseWriter, r *http.Request) {

	//duration

	if tkShardsInBlock.GetLength() != dkt3t2.GetLength() {
		t2Copy := tkShardRecv.Copy()
		t3Copy := tkShardsInBlock.Copy()
		for t3k, t3v := range t3Copy {
			t2v := t2Copy[t3k]
			dkt3t2.AddToKeeper(t3k, t3v.Sub(t2v))
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dkt3t2.ToString("b")))
}

//////////////
//helper funcs

func BroadcastMessage(message dataStructure.Message, api string, peers map[string]bool) {
	if message.HopCount > 0 {
		message.HopCount = 0 //message.HopCount - 1 // broadcasting
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range peers {
				_, err := http.Post(k+api, "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					beaconPeers.Delete(k)
				}
			}
		}
	}
}

func BroadcastMessageToShardMiners(message dataStructure.Message, api string, peers map[string]bool) {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1 // broadcasting
		messageJson, err := json.Marshal(message)
		if err == nil {
			if len(peers) > 0 {
				for k, _ := range peers {
					_, err := http.Post(k+api, "application/json", bytes.NewBuffer(messageJson))
					if err != nil {
						//todo : get new miners for the shard from register or something else
					}
				}
			}

		}
	}
}

func sendTxPostReq(reqBody []byte, shardId uint32) bool {
	//client := http.Client{}

	sid := strconv.Itoa(int(shardId))

	//find a peer for the given shard
	if peerAdd, ok := shardPeersForBeacon[shardId]; ok { //if found
		return sendTxMessageToShard(peerAdd, reqBody)

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
				return sendTxMessageToShard(k, reqBody)
			}

		}
		return false
	} // end of else (if peer not found)

}

func sendTxMessageToShard(peerAdd string, reqBody []byte) bool {
	client := http.Client{}
	url := peerAdd + "/shard/transaction/" //shard/transaction/
	fmt.Println("Sending to shard : " + url)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print("Cannot create PostReq, err " + err.Error())
	}
	_, err = client.Do(req) // sending tx message to "one" shard miner
	if err != nil {
		return false
	}
	return true

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
		sid := fmt.Sprint(shardId)
		sb.WriteString("Shard Id: " + sid + ", Shard Peer: " + shardPeer + "\n")
	}

	return sb.String()
}

func AskForBeaconBlock(height int32, parentHash string) bool {
	peerList := beaconPeers.Copy()
	for i, _ := range peerList {
		resp, err := http.Get(i + "/beacon/block/" + strconv.Itoa(int(height)) + "/" + parentHash)
		if err == nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusInternalServerError {
			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 == nil {
				block := blockchain.DecodeFromJSON(string(body))
				if !beaconSbc.CheckParentHash(block) && block.Header.Height-1 > 0 {
					if AskForBlock(block.Header.Height-1, block.Header.ParentHash) {
						//IsOpenTransaction(block.Value)
						beaconSbc.Insert(block)
						return true
					}
				} else {
					//IsOpenTransaction(block.Value)
					beaconSbc.Insert(block)
					return true
				}
			}
		}
	}
	return false
}

func GenerateBeaconBlocks() {
	for SHARD_ID == utils.BEACON_ID {
		mpt, flag := shardPool.BuildMpt()
		if flag {
			latestBlocks := beaconSbc.GetLatestBlocks()
			var latestBlock blockchain.Block
			if latestBlocks != nil && len(latestBlocks) != 0 {
				latestBlock = latestBlocks[rand.Intn(len(latestBlocks))]
			}
			for latestBlocks == nil || latestBlock.Header.Height == beaconSbc.GetLatestBlocks()[0].Header.Height {
				parentHash := "genesis"
				if latestBlocks != nil {
					parentHash = latestBlock.Header.Hash
				}
				block := beaconSbc.GenBlock(beaconSbc.GetLength()+1, parentHash, mpt, "0", identity.PublicKey, blockchain.SHARD)

				shardInBlockTime := time.Now() //t3

				shardPool.DeleteShards(block.Value)
				fmt.Println("Size of shard pool : " + shardPool.Show())
				beaconSbc.Insert(block)
				message := dataStructure.Message{
					Type: dataStructure.BLOCK,
					//Transaction: transaction.Transaction{},
					Block:    block,
					HopCount: 1,
					NodeId:   SELF_ADDR,
				}
				message.Sign(identity)
				latestBlocks = beaconSbc.GetLatestBlocks()
				fmt.Println(latestBlock.Header.Height)
				latestBlocks = beaconSbc.GetLatestBlocks()                         //getting latest block
				go BroadcastMessage(message, "/beacon/block/", beaconPeers.Copy()) // BroadcastBeaconBlockMessage
				//broadcast to all miner of all shards
				go BroadcastMessageToShardMiners(message, "/shard/beacon/", sameShardPeers.Copy())

				//add t3 to tkShardsInBlock
				tkShardsInBlock.AddTxIdsFromBeaconBlockToKeeper(message.Block.Value.Raw_db, shardInBlockTime) //t3
			}
		}

		random := 10 //rand.Intn(5) + 7
		time.Sleep(time.Second * time.Duration(random))
	}
}
