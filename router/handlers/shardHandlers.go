package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

func InitShardHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
	transactionPool = transaction.NewTransactionPool(SHARD_ID)
	sameShardPeers = peerList.NewPeerList(SHARD_ID)
	sbc = blockchain.NewBlockChain()
}

func StartShardMiner(w http.ResponseWriter, r *http.Request) {
	RegisterToServer(REGISTRATION_SERVER+"/register/")
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newPeers := peerList.NewPeerList(SHARD_ID)
			newPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			flag := false
			for server, _ := range newPeers.Copy() {
				go RegisterToServer(server+"/shard/peers/")
				sameShardPeers.Add(server)
				if !flag {
					DownloadBlockchain(server)
					flag = true
				}
			}
		}
	}
}

func DownloadBlockchain(server string) {
	resp, err := http.Get(server + "/shard/upload/")
	if err == nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			sbc.UpdateEntireBlockChain(string(respBody))
		}
	}
}

func RegisterShard(w http.ResponseWriter, r *http.Request) {
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
	sameShardPeers.Add(registerInfo.Address)
	w.WriteHeader(http.StatusOK)
}

func GetPeerListForShard(w http.ResponseWriter, r *http.Request) {
	respBody, err := sameShardPeers.PeerMapToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respBody))
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	message := dataStructure.Message{}
	if message.Type != dataStructure.TRANSACTION {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	transactionPool.AddToTransactionPool(message.Transaction)
	go BroadcastTransaction(message)
}

func BroadcastTransaction(message dataStructure.Message) {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range sameShardPeers.Copy() {
				_, err := http.Post(k+"/shard/transaction/", "application/json", bytes.NewBuffer(messageJson))
				if err != nil {
					sameShardPeers.Delete(k)
				}
			}
		}
	}
}

func ShowAllTransactionsInPool(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(transactionPool.Show()))
}

func GetBlock(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	defer r.Body.Close()
	message := dataStructure.Message{}
	if message.Type != dataStructure.BLOCK {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	sbc.Insert(message.Block)
	go BroadcastBlock(message)
}

func GenShardBlock() {
	for {
		mpt, flag := transactionPool.BuildMpt()
		if flag {
			latestBlocks := sbc.GetLatestBlocks()
			var latestBlock blockchain.Block
			if latestBlocks != nil && len(latestBlocks) != 0 {
				latestBlock = latestBlocks[rand.Intn(len(latestBlocks))]
			}
			for latestBlocks == nil || latestBlock.Header.Height == sbc.GetLatestBlocks()[0].Header.Height {
				parentHash := "Genesis"
				if latestBlocks != nil {
					parentHash = latestBlock.Header.Hash
				}
				block := sbc.GenBlock(sbc.GetLength()+1, parentHash, mpt, "0", identity.PublicKey, blockchain.TRANSACTION)
				sbc.Insert(block)
				message := dataStructure.Message{
					Type:        dataStructure.BLOCK,
					Transaction: transaction.Transaction{},
					Block:       block,
					HopCount:    1,
				}
				message.Sign(identity)
			}
		}
	}
}

func BroadcastBlock(message dataStructure.Message) {
	// ToDo: broadcast block to peers with same shard id.
}

func UploadBlockchain(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := sbc.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockChainJson))
	}
}
