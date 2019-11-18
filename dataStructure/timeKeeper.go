package dataStructure

import (
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
	tm.mux.Lock()
	defer tm.mux.Unlock()
	return tm.Copy()
}

func (tm *TimeMap) AddToKeeper(shardId string, time time.Time) {
	tm.mux.Lock()
	defer tm.mux.Unlock()
	tm.keeper[shardId] = time
}
