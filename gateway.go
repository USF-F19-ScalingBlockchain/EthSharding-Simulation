package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/EthSharding-Simulation/dataStructure/transaction"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Input struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Result  []transaction.EthTx `json:"result"`
}

var beaconMiner = "http://localhost:8000"
var dataset = [7]string{ "0x", "chainlink", "cryptokitties", "dice2win", "fairwin", "makerdao", "tether"}

func main() {
	if os.Args[1] == "-r" {
		input := Input{}
		for _, v := range dataset {
			newInput := doData(v)
			for _, v := range newInput.Result {
				input.Result = append(input.Result, v)
			}
		}
		convertToTransaction(input)
	} else {
		input := doData(os.Args[1])
		convertToTransaction(input)
	}

}

func doData(dataset string) Input {
	data, err := ioutil.ReadFile("input/"+ dataset +"/raw/transactions.json")
	fmt.Println("input/"+ dataset +"/raw/transactions.json")
	input := Input{}
	if err == nil {
		err = json.Unmarshal(data, &input)
		if err != nil {
			fmt.Println("Parsing Error, ", err)
		}
	}
	return input
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
