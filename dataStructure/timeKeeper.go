package dataStructure

import (
	"github.com/EthSharding-Simulation/dataStructure/shard"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"strings"
	"sync"
	"time"
)

type TimeMap struct {
	keeper map[string]time.Time
	mux    sync.Mutex
}

func NewTimeMap() TimeMap {
	return TimeMap{
		keeper: make(map[string]time.Time),
		mux:    sync.Mutex{},
	}
}

func (tm *TimeMap) Copy() map[string]time.Time {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	newKeep := make(map[string]time.Time)
	for k, v := range tm.keeper {
		newKeep[k] = v
	}
	return newKeep
}

func (tm *TimeMap) GetKeeper() map[string]time.Time {
	return tm.Copy()
}

func (tm *TimeMap) AddToKeeper(txId string, time time.Time) {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	if _, ok := tm.keeper[txId]; !ok {
		tm.keeper[txId] = time
	}
}

func (tm *TimeMap) AddTxIdsFromShardToKeeper(openTxSet map[string]transaction.Transaction, time time.Time) {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	for txId, _ := range openTxSet {
		if _, ok := tm.keeper[txId]; !ok {
			tm.keeper[txId] = time
		}

	}
}

func (tm *TimeMap) AddTxIdsFromBeaconBlockToKeeper(shardMap map[string]string, time time.Time) {
	for _, shardStr := range shardMap {
		shard := shard.JsonToShard(shardStr)
		tm.AddTxIdsFromShardToKeeper(shard.OpenTransactionSet, time)
	}
}

func (tm *TimeMap) GetLength() int {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	return len(tm.keeper)
}

func (tm *TimeMap) ToString() string {
	sb := strings.Builder{}
	sb.WriteString("Showing time: \n")
	for k, v := range tm.Copy() {
		sb.WriteString(k + " : " + v.String() + "\n")
	}
	return sb.String()
}

/// duration

type DurationMap struct {
	keeper map[string]time.Duration
	mux    sync.Mutex
}

func NewDurationMap() DurationMap {
	return DurationMap{
		keeper: make(map[string]time.Duration),
		mux:    sync.Mutex{},
	}
}

func (dm *DurationMap) Copy() map[string]time.Duration {
	dm.mux.Lock()
	defer dm.mux.Unlock()
	newKeep := make(map[string]time.Duration)
	for k, v := range dm.keeper {
		newKeep[k] = v
	}
	return newKeep
}

func (dm *DurationMap) GetKeeper() map[string]time.Duration {
	return dm.Copy()
}

func (dm *DurationMap) AddToKeeper(txId string, time time.Duration) {
	dm.mux.Lock()
	defer dm.mux.Unlock()
	if _, ok := dm.keeper[txId]; !ok {
		dm.keeper[txId] = time
	}
}

func (dm *DurationMap) GetLength() int {
	dm.mux.Lock()
	defer dm.mux.Unlock()
	return len(dm.keeper)
}

func (dm *DurationMap) ToString() string {
	sb := strings.Builder{}
	sb.WriteString("Showing duration: \n")
	for k, v := range dm.Copy() {
		sb.WriteString(k + " : " + v.String() + "\n")
	}
	return sb.String()
}
