package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"github.com/EthSharding-Simulation/utils"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type Input struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Result  []transaction.EthTx `json:"result"`
}

var beaconMiner = "http://localhost:8000"

func main() {
	data, err := ioutil.ReadFile("input/" + utils.Dataset + "/raw/transactions.json")
	if err == nil {
		input := Input{}
		err = json.Unmarshal(data, &input)
		if err != nil {
			fmt.Println("Parsing Error, ", err)
		}
		convertToTransaction(input)
	}
}

func convertToTransaction(input Input) {
	sort.Slice(input.Result, func(i, j int) bool {
		val1, _ := strconv.ParseInt(input.Result[i].TimeStamp, 10, 32)
		val2, _ := strconv.ParseInt(input.Result[j].TimeStamp, 10, 32)
		return val1 < val2
	})
	for _, val := range input.Result {
		if len(val.To) != 0 && len(val.From) != 0 {
			value, err := strconv.ParseFloat(val.Value, 64)
			if err == nil {
				tx := transaction.NewTransaction(val.Hash, val.From, val.To, value)
				txJson, err := json.Marshal(tx)
				fmt.Println(string(txJson))
				if err == nil {
					http.Post(beaconMiner+"/beacon/Tx/receive/", "application/json", bytes.NewBuffer(txJson))
				}
			}
		}
	}
}
