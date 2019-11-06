package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EthSharding-Simulation/dataStructure/mpt"
	t "github.com/EthSharding-Simulation/dataStructure/transaction"
	"golang.org/x/crypto/sha3"
)

// Header struct defines the header of each block
type Head struct {
	Height     int32  //`json:"height"`
	Timestamp  int64  //`json:"timestamp"`
	Hash       string //`json:"hash"`
	ParentHash string //`json:"parenthash"`
	Size       int32  // `json:"parenthash"`
	Nonce      string //'json:"nonce"'
	Miner      t.PublicIdentity
	BlockType  BlockType
}

type BlockType int

const(
	SHARD BlockType = 0
	TRANSACTION BlockType = 1
)

func (blockType BlockType) String() string {
	switch blockType {
	case SHARD:
		return "SHARD"
	case TRANSACTION:
		return "TRANSACTION"
	default:
		return ""
	}
}

// Block struct defines the block
type Block struct {
	Header Head
	Value  mpt.MerklePatriciaTrie
}

// BlockJson is a block struct for json
type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	Nonce      string            `json:"nonce"`
	Miner      t.PublicIdentity  `json:"miner"`
	MPT        map[string]string `json:"mpt"`
	BlockType BlockType			 `json:"blockType"`
}

// Initial function a Block initializes the block for height, parentHash and Value
func (block *Block) Initial(height int32, parentHash string, value mpt.MerklePatriciaTrie, nonce string, miner t.PublicIdentity, blockType BlockType) {

	block.Header.Timestamp = time.Now().Unix()
	block.Header.Height = height
	block.Header.ParentHash = parentHash
	block.Value = value
	block.Header.Size = int32(len([]byte(block.Value.String()))) // mpt converted to string and then to byte array
	block.Header.Nonce = nonce
	block.Header.Miner = miner
	block.Header.Hash = block.Hash()
	block.Header.BlockType = blockType

}

// DecodeFromJSON func takes json string of type blockJson and converts it into a Block // proxy for : DecodeFromJson
func DecodeFromJSON(jsonString string) Block {

	// block := Block{}
	blockJson := BlockJson{}

	err := json.Unmarshal([]byte(jsonString), &blockJson)
	if err != nil {
		fmt.Println("DecodeFromJSON  in Block.go : block Err : ", err)
		return Block{}
	}
	return DecodeToBlock(
		blockJson.Height,
		blockJson.Timestamp,
		blockJson.Hash,
		blockJson.ParentHash,
		blockJson.Size,
		blockJson.Nonce,
		blockJson.Miner,
		blockJson.MPT,
		blockJson.BlockType)
}

// DecodeToBlock func creates a type block from from all given parameters
func DecodeToBlock(height int32, timestamp int64, hash string, parentHash string, size int32, nonce string,
	miner t.PublicIdentity, keyValueMap map[string]string, blockType BlockType) Block {

	block := Block{}
	block.Header.Height = height
	block.Header.Timestamp = timestamp
	block.Header.Hash = hash
	block.Header.ParentHash = parentHash
	block.Header.Size = size
	block.Header.Nonce = nonce
	block.Header.Miner = miner
	block.Header.BlockType = blockType

	//creating mpt from key - value pairs
	blockMPT := mpt.MerklePatriciaTrie{}
	blockMPT.Initial()
	for k, v := range keyValueMap {
		blockMPT.Insert(k, v)
	}
	block.Value = blockMPT
	//fmt.Println("in DecodeToBlock of Block.go : root : ", block.Value.Root)
	return block
}

// EncodeToJSON func takes type Block and converts it into json string
func EncodeToJSON(block *Block) string {

	blockForJson := BlockJson{
		Height:     block.Header.Height,
		Timestamp:  block.Header.Timestamp,
		Hash:       block.Header.Hash,
		ParentHash: block.Header.ParentHash,
		Size:       block.Header.Size,
		Nonce:      block.Header.Nonce,
		Miner:      block.Header.Miner,
		MPT:        block.Value.GetAllKeyValuePairs(),
		BlockType:  block.Header.BlockType,
	}

	jsonByteArray, err := json.Marshal(blockForJson)
	if err != nil {
		return ""
	}
	//jsonString = string(jsonByteArray)
	return string(jsonByteArray) //empty jsonString if not encoded else some value
}

//Hash func takes an instance of block and hashes it
//hash_str := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash +
//     b.Value.Root + string(b.Header.Size) + block.Header.Nonce
func (block *Block) Hash() string {
	var hashStr string

	hashStr = string(block.Header.Height) + string(block.Header.Timestamp) + string(block.Header.ParentHash) +
		string(block.Value.Root) + string(block.Header.Size) + block.Header.Nonce + block.Header.Miner.PublicIdentityToJson() + block.Header.BlockType.String()

	sum := sha3.Sum256([]byte(hashStr))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}
