package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	RegisterToServer(REGISTRATION_SERVER +"/register/")
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
			go GenShardBlock()
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
	go Broadcast(message, "/shard/transaction/")
}

func Broadcast(message dataStructure.Message, uri string) {
	//if message.Verify() && message.HopCount > 0 {
	if message.HopCount > 0 {
		message.HopCount = message.HopCount - 1
		messageJson, err := json.Marshal(message)
		if err == nil {
			for k, _ := range sameShardPeers.Copy() {
				if k != message.NodeId {
					_, err := http.Post(k + uri, "application/json", bytes.NewBuffer(messageJson))
					if err != nil {
						sameShardPeers.Delete(k)
					}
				}
			}
		}
	}
}

func ShowAllTransactionsInPool(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(transactionPool.Show()))
}

func AddShardBlock(w http.ResponseWriter, r *http.Request) {
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
	transactionPool.DeleteTransactions(message.Block.Value)
	sbc.Insert(message.Block)
	go Broadcast(message, "/shard/block/")
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
				parentHash := "genesis"
				if latestBlocks != nil {
					parentHash = latestBlock.Header.Hash
				}
				block := sbc.GenBlock(sbc.GetLength()+1, parentHash, mpt, "0", identity.PublicKey, blockchain.TRANSACTION)
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
				go Broadcast(message, "/shard/block/")
			}
		}
	}
}

func UploadBlockchain(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := sbc.BlockChainToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(blockChainJson))
	}
}

func ShowShard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", sameShardPeers.Show(), sbc.Show())
}
