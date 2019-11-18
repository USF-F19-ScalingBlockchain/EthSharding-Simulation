package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure"
	"github.com/EthSharding-Simulation/dataStructure/blockchain"
	"github.com/EthSharding-Simulation/dataStructure/mpt"
	"github.com/EthSharding-Simulation/dataStructure/peerList"
	"github.com/EthSharding-Simulation/dataStructure/shard"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/utils"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func InitShardHandler(host string, port int32, shardId uint32) {
	SHARD_ID = shardId
	SELF_ADDR = host + ":" + strconv.Itoa(int(port))
}

func StartShardMiner(w http.ResponseWriter, r *http.Request) {
	transactionPool = transaction.NewTransactionPool(SHARD_ID)
	sameShardPeers = peerList.NewPeerList(SHARD_ID)
	sbc = blockchain.NewBlockChain()
	openTransactionSet = transaction.NewOpenTransactionSet()
	RegisterToServer(REGISTRATION_SERVER + "/register/", SHARD_ID)
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(SHARD_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			newPeers := peerList.NewPeerList(SHARD_ID)
			newPeers.InjectPeerMapJson(respBody, SELF_ADDR)
			flag := false
			for server, _ := range newPeers.Copy() {
				go RegisterToServer(server + "/shard/peers/", SHARD_ID)
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
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if message.Type != dataStructure.TRANSACTION {
		w.WriteHeader(http.StatusBadRequest)
	}
	recvLock.Lock()
	recvTime[message.Transaction.Id] = time.Now()
	recvLock.Unlock()
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
					_, err := http.Post(k+uri, "application/json", bytes.NewBuffer(messageJson))
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
	err = json.Unmarshal(reqBody, &message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	}
	if message.Type != dataStructure.BLOCK {
		w.WriteHeader(http.StatusBadRequest)
	}
	transactionPool.DeleteTransactions(message.Block.Value)
	if !sbc.CheckParentHash(message.Block) && message.Block.Header.Height-1 > 0 {
		if AskForBlock(message.Block.Header.Height-1, message.Block.Header.ParentHash) {
			IsOpenTransaction(message.Block.Value, false)
			sbc.Insert(message.Block)
		}
	} else {
		IsOpenTransaction(message.Block.Value, false)
		sbc.Insert(message.Block)
	}
	go Broadcast(message, "/shard/block/")
}

func AskForBlock(height int32, parentHash string) bool {
	peerList := sameShardPeers.Copy()
	for i, _ := range peerList {
		resp, err := http.Get(i + "/block/" + strconv.Itoa(int(height)) + "/" + parentHash)
		if err == nil && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusInternalServerError {
			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 == nil {
				block := blockchain.DecodeFromJSON(string(body))
				if !sbc.CheckParentHash(block) && block.Header.Height-1 > 0 {
					if AskForBlock(block.Header.Height-1, block.Header.ParentHash) {
						IsOpenTransaction(block.Value, false)
						sbc.Insert(block)
						return true
					}
				} else {
					IsOpenTransaction(block.Value, false)
					sbc.Insert(block)
					return true
				}
			}
		}
	}
	return false
}

func UploadShardBlock(w http.ResponseWriter, r *http.Request) {
	pathParam := strings.Split(r.URL.Path, "/")
	height, err := strconv.ParseInt(pathParam[3], 10, 32)
	fmt.Println(height)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: InternalServerError. " + err.Error()))
	} else {
		hash := pathParam[4]
		fmt.Println(hash)
		block, flag := sbc.GetBlock(int32(height), hash)

		if !flag {
			w.WriteHeader(http.StatusNoContent)
		} else {
			blockJson := block.EncodeToJSON()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(blockJson))
		}
	}
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
				IsOpenTransaction(block.Value, false)
				sbc.Insert(block)
				message := dataStructure.Message{
					Type:        dataStructure.BLOCK,
					Transaction: transaction.Transaction{},
					Block:       block,
					HopCount:    1,
					NodeId:      SELF_ADDR,
					TimeStamp:   time.Now(),
				}
				message.Sign(identity)
				if sbc.GetLength()%utils.SHARD_INTERVAL == 0 {
					fmt.Println("Submitting to beacon")
					go SubmitToBeacon()
				}
				latestBlocks = sbc.GetLatestBlocks()
				go Broadcast(message, "/shard/block/")
			}
		}
		random := rand.Intn(3) + 4
		time.Sleep(time.Second * time.Duration(random))
	}
}

func IsOpenTransaction(mpt mpt.MerklePatriciaTrie, ignoreFlag bool) {
	for k, v := range mpt.Raw_db {
		tx := transaction.JsonToTransaction(v)
		if transactionPool.IsOpenTransaction(tx) {
			if !ignoreFlag {
				openTransactionSet.AddTransaction(tx)	
			}
		} else {
			finalizeLock.Lock()
			finalizeTime[k] = time.Now() // ToDo: Replace time.now with tx.Timestamp
			finalizeLock.Unlock()
		}
	}
}

func SubmitToBeacon() {
	if len(beaconPeers.Copy()) == 0 {
		UpdateBeaconPeer()
		fmt.Println("len(beaconPeers.Copy()) ", len(beaconPeers.Copy()))
	}
	flag := true
	for flag {
		for k, _ := range beaconPeers.Copy() {
			shard := shard.NewShard("abc", time.Now(), SELF_ADDR, openTransactionSet.CopyAndClear())
			message := dataStructure.Message{
				Type:      dataStructure.SHARD,
				Shard:     shard,
				HopCount:  1,
				NodeId:    SELF_ADDR,
				TimeStamp: time.Time{},
			}
			message.Sign(identity)
			messageJson, err := json.Marshal(message)
			if err == nil {
				resp, err := http.Post(k+"/beacon/shard/", "application/json", bytes.NewBuffer(messageJson))
				if err != nil || resp.StatusCode != http.StatusOK {
					fmt.Println("Error in ", k+"/beacon/shard/")
					beaconPeers.Delete(k)
					UpdateBeaconPeer()
				} else {
					flag = false
					break
				}
			}
		}
	}
}

func UpdateBeaconPeer() {
	resp, err := http.Get(REGISTRATION_SERVER + "/register/peers/" + strconv.Itoa(int(utils.BEACON_ID)))
	if err == nil && resp.StatusCode != http.StatusBadRequest {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			respBody := string(respBytes)
			beaconPeers.InjectBeaconPeerMapJson(respBody)
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

func GetBeaconBlock(w http.ResponseWriter, r *http.Request) {
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
	IsOpenTransaction(message.Block.Value, true)
}

func ShowReceivedTimes(w http.ResponseWriter, r *http.Request) {
	recvLock.Lock()
	defer recvLock.Unlock()
	recvTimeJson, err := json.Marshal(recvTime)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(recvTimeJson)
	}
}

func ShowFinalizeTimes(w http.ResponseWriter, r *http.Request) {
	finalizeLock.Lock()
	defer finalizeLock.Unlock()
	finalizeTimeJson, err := json.Marshal(finalizeTime)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(finalizeTimeJson)
	}
}

func GetFinalTimePerTransaction(w http.ResponseWriter, r *http.Request) {
	finalizeLock.Lock()
	recvLock.Lock()
	timeDiff := make(map[string]string)
	for k, f := range finalizeTime {
		if r, ok := recvTime[k]; ok {
			timeDiff[k] = f.Sub(r).String()
		}
	}
	timeDiffJson, err := json.Marshal(timeDiff)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(timeDiffJson)
	defer func() {
		recvLock.Unlock()
		finalizeLock.Unlock()
	}()
}
