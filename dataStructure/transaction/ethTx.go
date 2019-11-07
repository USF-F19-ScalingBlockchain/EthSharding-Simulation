package transaction

import (
	"encoding/json"
	"log"
)

//{
//"blockNumber": "4774826",
//"timeStamp": "1513917067",
//"hash": "0x7e06a2ff18db8f68b1d1b865abab5ada99676194846d44acdb5e1f5bd2879e68",
//"nonce": "79",
//"blockHash": "0xbee7b455f5616aa8febc6ef0851c702e95575aa65acc03a8fb7b30fde51a3dd3",
//"transactionIndex": "45",
//"from": "0x28e37925e4e628f37332901a98cfbba66f314b23",
//"to": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
//"value": "17900000000000000000",
//"gas": "150000",
//"gasPrice": "40000000000",
//"isError": "0",
//"txreceipt_status": "1",
//"input": "0xd0e30db0",
//"contractAddress": "",
//"cumulativeGasUsed": "1806753",
//"gasUsed": "43346",
//"confirmations": "3806886"
//}

type EthTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	Txreceipt_status  string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

func (tx *EthTx) EthTxToJson() string {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return string(txJson)
}

func JsonToEthTx(txJson string) EthTx {
	tx := EthTx{}
	err := json.Unmarshal([]byte(txJson), &tx)
	if err != nil {
		log.Println("Error in unmarshalling Transaction, err - ", err)
		log.Println("String given to unmarshall Transaction, ================> \n ", txJson, "\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n")
	}

	return tx
}
