package blockchain

import (
	"crypto/rsa"
	"sync"

	"github.com/EthSharding-Simulation/dataStructure/mpt"
)

//SyncBlockChain struct is main - shared - common - datastu
type SyncBlockChain struct {
	bc  Blockchain `json:"blockchain"`
	mux sync.Mutex `json:"mux"`
}

// NewBlockChain func generates a new syncBlockchain
func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{bc: NewBlockchain()}
}

//Get func takes height as input and returns list of block at that height
func (sbc *SyncBlockChain) GetLength() int32 {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}

//Get func takes height as input and returns list of block at that height
func (sbc *SyncBlockChain) Get(height int32) ([]Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height), false
}

// GetBlock func takes height and hash as parameter and returns a block
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (Block, bool) {

	sbc.mux.Lock()
	defer sbc.mux.Unlock()

	blks, found := sbc.Get(height)
	//blocks, found := sbc.bc.Get(height)
	if found == true {
		for _, b := range blks {
			if b.Header.Hash == hash {
				return b, true
			}
		}
	}
	return Block{}, false
}

//Insert func inserts a block into blockchain in safe way
func (sbc *SyncBlockChain) Insert(block Block) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	sbc.bc.Insert(block)
}

// CheckParentHash func takes a block and checks if parent hash exists and return true or false
func (sbc *SyncBlockChain) CheckParentHash(insertBlock Block) bool {
	//sbc.mux.Lock()
	//defer sbc.mux.Unlock() apr 4

	if insertBlock.Header.Height > 1 { // good coz genesis created
		pblocks, found := sbc.Get(insertBlock.Header.Height - 1)
		if found == true {
			for _, pb := range pblocks {
				if pb.Header.Hash == insertBlock.Header.ParentHash {
					//log.Println("Parent Hash found at height :", pb.Header.Height)
					return true
				}
			}
		}
	}
	return false
}

// UpdateEntireBlockChain func takes a json and updates the existing blockchain
func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	DecodeBlockchainFromJSON(&sbc.bc, blockChainJson)
}

// BlockChainToJson converts blockchain to json string
func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return EncodeBlockchainToJSON(&sbc.bc), nil
}

// GenBlock finc takes in a mpt and returns a block for the node
// takes parentat list[0] in random height
func (sbc *SyncBlockChain) GenBlock(height int32, parentHash string, mpt mpt.MerklePatriciaTrie, nonce string, miner *rsa.PublicKey, blockType BlockType) Block {

	var newBlock Block
	newBlock.Initial(height, parentHash, mpt, nonce, miner, blockType)

	//fmt.Println(" blockHash : ", newBlock.Header.Hash)
	return newBlock
}

//// GenBlock finc takes in a mpt and returns a block for the node
//// takes parentat list[0] in random height
//func (sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie, nonce string) block.Block {
//
//	var parentHash string
//	var parentHeight int32
//	var blockList []block.Block
//	var found bool
//	currHeight := sbc.bc.Length
//
//	if currHeight == 0 {
//		parentHash = "genesis"
//	}
//	for currHeight >= 1 {
//		//sun
//		//blockList, found = sbc.Get(currHeight) //todo here // todo
//		//port, _ := strconv.Atoi(os.Args[1])
//		//if port == 7001 {
//		//	random := rand.Int() % 10
//		//	blockList, found = sbc.Get(currHeight - int32(random))
//		//} else {
//		blockList, found = sbc.Get(currHeight)
//		//}
//		//
//		//sun
//		if found == true {
//			random := 0
//			if len(blockList) > 1 {
//				random = rand.Int() % (len(blockList) - 1)
//			}
//			parentHash = blockList[random].Header.Hash
//			parentHeight = blockList[random].Header.Height
//			break
//		} else {
//			currHeight--
//			if currHeight == 0 { //apr5 ///////
//				parentHash = "genesis"
//				parentHeight = 1
//				break //apr5 //////////////////
//			}
//		}
//	}
//
//	//fmt.Println(" Current Height : ", currHeight)
//	//fmt.Println(" parentHash : ", parentHash)
//
//	var newBlock block.Block
//	newBlock.Initial(parentHeight+1, parentHash, mpt, nonce)
//
//	//fmt.Println(" blockHash : ", newBlock.Header.Hash)
//	return newBlock
//}

// Show func returns blockchain in displayable format
func (sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}

func (sbc *SyncBlockChain) GetLatestBlocks() []Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks() //blockchain.Chain[blockchain.Length]
}

func (sbc *SyncBlockChain) GetParentBlock(blk Block) Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(blk)
}
