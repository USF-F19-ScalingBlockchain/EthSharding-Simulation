package shard

import (
	"encoding/json"
	"github.com/EthSharding-Simulation/dataStructure/mpt"
	"github.com/EthSharding-Simulation/utils"
	"sync"
)

// ToDo: Check for duplicate shard before inserting.

type ShardPool struct {
	Pool 	map[string]Shard
	mux  	sync.Mutex
}

func NewShardPool() ShardPool {
	return ShardPool{
		Pool: make(map[string]Shard),
		mux:  sync.Mutex{},
	}
}

func (shardPool *ShardPool) AddToShardPool(shard Shard) { //duplicates in transactinon pool
	shardPool.mux.Lock()
	defer shardPool.mux.Unlock()

	if _, ok := shardPool.Pool[shard.Id]; !ok {
		shardPool.Pool[shard.Id] = shard
	}
}

func (shardPool *ShardPool) DeleteFromTransactionPool(shadId string) {
	shardPool.mux.Lock()
	defer shardPool.mux.Unlock()

	delete(shardPool.Pool, shadId)
}

// ToDo: Should we delete shard from shard pool after building mpt?
func (shardPool *ShardPool) BuildMpt() (mpt.MerklePatriciaTrie, bool) {
	shardPool.mux.Lock()
	shardPool.mux.Unlock()
	shardMpt := mpt.MerklePatriciaTrie{}
	shardMpt.Initial()
	if len(shardPool.Pool) < utils.MIN_TX_POOL_SIZE {
		return shardMpt, false
	}
	for i, _ := range shardPool.Pool {
		shardsJson, err := json.Marshal(shardPool.Pool[i])
		if err == nil {
			shardMpt.Insert(shardPool.Pool[i].Id, string(shardsJson))
		}
	}
	return shardMpt, true
}