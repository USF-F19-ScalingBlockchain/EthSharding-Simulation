package transaction

import (
	"sync"
)

type OpenTransactionSet struct {
	openTransactions 	map[string]Transaction
	mux 				sync.Mutex
}

func NewOpenTransactionSet() OpenTransactionSet {
	return OpenTransactionSet{
		openTransactions: make(map[string]Transaction),
	}
}

func (openTxSet *OpenTransactionSet) AddTransaction(tx Transaction) {
	openTxSet.mux.Lock()
	defer openTxSet.mux.Unlock()
	openTxSet.openTransactions[tx.Id] = tx
}

func (openTxSet *OpenTransactionSet) DeleteTransaction(tx Transaction) {
	openTxSet.mux.Lock()
	defer openTxSet.mux.Unlock()
	delete(openTxSet.openTransactions, tx.Id)
}

func (openTxSet *OpenTransactionSet) CopyAndClear() map[string]Transaction {
	copy := make(map[string]Transaction)
	for i, val := range openTxSet.openTransactions {
		copy[i] = val
	}
	openTxSet.openTransactions = make(map[string]Transaction)
	return copy
}
