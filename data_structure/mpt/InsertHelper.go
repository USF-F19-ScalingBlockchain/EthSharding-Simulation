package mpt

import (
	"reflect"
)

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	//fmt.Println("Key - Value being inserted:", key, " - ", new_value)

	pathLeft := StringToHexArray(key)
	if mpt.Root == "" {
		//pathLeft = append(pathLeft, 16)
		//leafNode := Node{}
		//leafNode.node_type = 2
		//leafNode.flag_value.encoded_prefix = compact_encode(pathLeft)
		//
		//leafNode.flag_value.value = new_value
		//mpt.Root = mpt.insertHelper(leafNode, nil, "", nil)
		//>>
		mpt.Root = mpt.insertHelper1(Node{}, pathLeft, new_value)

	} else {

		//RootNodeHash := mpt.Root
		//mpt.Root = mpt.insertHelper(mpt.db[RootNodeHash], pathLeft, new_value, []uint8{})
		//>>
		mpt.Root = mpt.insertHelper1(mpt.db[mpt.Root], pathLeft, new_value)
	}

}

// input - pathLeft, new_value; output - newHashValue
func (mpt *MerklePatriciaTrie) insertHelper1(currentNode Node, pathLeft []uint8, newValue string) string {
	////fmt.Println("in Insert Helper 1,nodeType :", currentNode.node_type, " pathLeft : ", pathLeft)
	currentNodeHash := currentNode.hash_node()

	if currentNode.node_type == 0 && len(mpt.db) == 0 { //1st node
		newNodeHash, newNode := createLeafNode(pathLeft, newValue)

		mpt.db = make(map[string]Node)
		mpt.Root = ""
		mpt.db[newNodeHash] = newNode
		mpt.Root = newNodeHash

		return newNodeHash
	}

	if len(pathLeft) > 0 {
		////fmt.Println("pathLeft > 0")
		switch {
		case currentNode.node_type == 1: //branch
			////fmt.Println("pathLeft > 0 and in branch")
			if currentNode.branch_value[pathLeft[0]] != "" {
				nextNode := mpt.db[currentNode.branch_value[pathLeft[0]]]
				newPathLeft := pathLeft[1:]

				currentNode.branch_value[pathLeft[0]] = mpt.insertHelper1(nextNode, newPathLeft, newValue)

				////fmt.Println("currentNode.branch_value[pathLeft[0]] : ", currentNode.branch_value[pathLeft[0]])

			} else {
				//expand trie from branch   //create a new leaf node
				////fmt.Println("pathLeft > 0 and in branch - to create a new leaf node")
				newLeafNodeHash, newLeafNode := createLeafNode(pathLeft[1:], newValue)
				mpt.db[newLeafNodeHash] = newLeafNode

				currentNode.branch_value[pathLeft[0]] = newLeafNodeHash

			}

			currentNodeNewHash := currentNode.hash_node()
			delete(mpt.db, currentNodeHash)
			mpt.db[currentNodeNewHash] = currentNode
			return currentNodeNewHash

		case currentNode.node_type == 2 && is_ext_node(currentNode.flag_value.encoded_prefix) == true: //ext
			////fmt.Println("pathLeft > 0 and in ext")
			currentNodeHexArray := compact_decode(currentNode.flag_value.encoded_prefix)
			//if len(pathLeft) >= len(currentNodeHexArray) { //pathLeft >= currentNodeHexArray
			if reflect.DeepEqual(pathLeft[:len(currentNodeHexArray)], currentNodeHexArray) { //both path are same
				nextNode := mpt.db[currentNode.flag_value.value]
				newPathLeft := pathLeft[len(currentNodeHexArray):]

				if len(newPathLeft) == 0 {
					//fmt.Println("ok path gets over at ext node") //bug
				}
				currentNode.flag_value.value = mpt.insertHelper1(nextNode, newPathLeft, newValue)

				currentNodeNewHash := currentNode.hash_node()
				delete(mpt.db, currentNodeHash)
				mpt.db[currentNodeNewHash] = currentNode
				return currentNodeNewHash

			} else {
				////fmt.Println("two paths not equal")
				////fmt.Println("matching : ", pathLeft, " --- ", currentNodeHexArray)
				counter := 0
				if len(pathLeft) >= len(currentNodeHexArray) {
					for i := range currentNodeHexArray {
						if pathLeft[i] == currentNodeHexArray[i] {
							counter = counter + 1
						} else {
							break
						}
					}
				} else {
					for i := range pathLeft {
						if pathLeft[i] == currentNodeHexArray[i] {
							counter = counter + 1
						} else {
							break
						}
					}
				}
				//fmt.Println("counter -- ", counter)
				if counter == 0 { //, it is extension then
					//totally different
					//create a branch, exsting value of ext ->in branch
					// plus new value at seperate index of branch
					//ext node + new node = branch node ->leftextnode + new node(pathleft[1:])
					////fmt.Println("0 in counter : ", counter)
					newBranchNode := Node{}
					newBranchNode.node_type = 1

					//connecting old trie
					if len(currentNodeHexArray) == 1 {
						////fmt.Println("currentNodeHexArray : ", currentNodeHexArray)
						newBranchNode.branch_value[currentNodeHexArray[0]] = currentNode.flag_value.value

					} else { //create new extensiom if old extension has length of greater than 1
						newExtNode1 := Node{}
						newExtNode1.node_type = 2
						newExtNode1.flag_value.encoded_prefix = compact_encode(currentNodeHexArray[1:])
						////fmt.Println("newExtNode1.flag_value.encoded_prefix : ", newExtNode1.flag_value.encoded_prefix)
						newExtNode1.flag_value.value = currentNode.flag_value.value
						////fmt.Println("currentNode.flag_value.value :", currentNode.flag_value.value)
						newExtNode1Hash := newExtNode1.hash_node()

						newBranchNode.branch_value[currentNodeHexArray[0]] = newExtNode1Hash
						////fmt.Println("currentNodeHexArray[0] : ", currentNodeHexArray[0])
						//delete currentextension node
						delete(mpt.db, currentNodeHash)
						mpt.db[newExtNode1Hash] = newExtNode1
					}

					newLeafNodeLeftPath := pathLeft[1:]
					////fmt.Println("newLeafNodeLeftPath : ", newLeafNodeLeftPath)
					newLeafNodeHash, newLeafNode := createLeafNode(newLeafNodeLeftPath, newValue)
					newBranchNode.branch_value[pathLeft[0]] = newLeafNodeHash
					mpt.db[newLeafNodeHash] = newLeafNode

					mpt.db[newBranchNode.hash_node()] = newBranchNode
					return newBranchNode.hash_node()

				} else {
					//some common path and then diverge
					////fmt.Println("some common path and then diverge")
					newExtNode := Node{}
					newExtNode.node_type = 2
					newExtNode.flag_value.encoded_prefix = compact_encode(currentNodeHexArray[:counter])
					////fmt.Println("Ext Node hex path :", currentNodeHexArray[:counter])

					//create new branch node
					newBranchNode := Node{}
					newBranchNode.node_type = 1

					//connect new insertion
					////fmt.Println("pathLeft : ", pathLeft[counter:])
					if len(pathLeft[counter:]) == 0 {
						newBranchNode.branch_value[16] = newValue
					} else {
						newLeafNode1Hash, newLeafNode1 := createLeafNode(pathLeft[counter+1:], newValue)
						newBranchNode.branch_value[pathLeft[counter]] = newLeafNode1Hash
						////fmt.Println("new leaf node inserted at : ", pathLeft[counter])
						////fmt.Println("new leaf hex path : ", pathLeft[counter+1:])
						mpt.db[newLeafNode1Hash] = newLeafNode1
					}

					//for old connection
					if len(pathLeft) >= len(currentNodeHexArray) && counter+1 == len(currentNodeHexArray) {
						//connecting new Branch node at currentNodeHexArray[counter+1] to
						//node that was connected to current node
						// //fmt.Println("value to be added in branch only")
						// //fmt.Println("counter : ", counter)
						////fmt.Println("currentNodeHexArray : ", currentNodeHexArray[counter])

						newBranchNode.branch_value[currentNodeHexArray[counter]] = currentNode.flag_value.value

					} else {
						//create a new ext node
						////fmt.Println("value to be added in complex")
						newExtNode1 := Node{}
						newExtNode1.node_type = 2
						////fmt.Println("current node left hex to be put in ext : ", currentNodeHexArray[counter+1:])
						newExtNode1.flag_value.encoded_prefix = compact_encode(currentNodeHexArray[counter+1:])
						newExtNode1.flag_value.value = currentNode.flag_value.value //connecting new Ext1 node to
						// the node that was connected to current node
						newExtNode1Hash := newExtNode1.hash_node()
						mpt.db[newExtNode1Hash] = newExtNode1
						////fmt.Println("newExtNode1 to be added at index : ", currentNodeHexArray[counter])
						newBranchNode.branch_value[currentNodeHexArray[counter]] = newExtNode1Hash

					}
					newBranchNodeHash := newBranchNode.hash_node()
					mpt.db[newBranchNodeHash] = newBranchNode
					newExtNode.flag_value.value = newBranchNodeHash //connecting new Ext node to new Branch node

					newExtNodeHash := newExtNode.hash_node()
					mpt.db[newExtNodeHash] = newExtNode
					delete(mpt.db, currentNode.hash_node())
					return newExtNodeHash

				}
			}
			///}
		case currentNode.node_type == 2 && is_ext_node(currentNode.flag_value.encoded_prefix) == false: //leaf
			//fmt.Println("pathLeft > 0 and in leaf")
			currentNodeHexArray := compact_decode(currentNode.flag_value.encoded_prefix)
			if reflect.DeepEqual(currentNodeHexArray, pathLeft) {
				//update value
				//fmt.Println("deep equal : to update value")
				currentNode.flag_value.value = newValue
				//fmt.Println("Adding in same position : ", newValue)
				currentNodeNewHash := currentNode.hash_node()
				delete(mpt.db, currentNodeHash)
				mpt.db[currentNodeNewHash] = currentNode
				return currentNodeNewHash

			} else {
				//2 conditions - a. len(pathLeft) <= len(currentNodeHexArray)
				//               b. len(pathLeft) > len(currentNodeHexArray)
				//fmt.Println("two paths not equal")
				counter := 0

				if len(pathLeft) <= len(currentNodeHexArray) {
					for i := range pathLeft {
						if currentNodeHexArray[i] == pathLeft[i] {
							counter = counter + 1
						} else {
							break
						}
					}
				} else if len(currentNodeHexArray) > 0 {
					////fmt.Println("currentNodeHexArray < pathleft, and  currentNodeHexArray > 0 ")
					for i := range currentNodeHexArray {
						if currentNodeHexArray[i] == pathLeft[i] {
							counter = counter + 1
						} else {
							break
						}
					}
				}

				//fmt.Println("counter : ", counter)
				//fmt.Println("pathLeft : ", pathLeft)
				//fmt.Println("currentNodeHexArray : ", currentNodeHexArray)

				if counter == 0 { //do not match at all    //create branch - 2 leaves
					////fmt.Println("counter = 0")
					newBranchNode := Node{}
					newBranchNode.node_type = 1

					//for old leaf
					if len(currentNodeHexArray) == 0 {
						newBranchNode.branch_value[16] = currentNode.flag_value.value
					} else if len(currentNodeHexArray) > 0 {
						newLeafNode1Hash, HnewLeafNode1 := createLeafNode(currentNodeHexArray[1:], currentNode.flag_value.value)
						mpt.db[newLeafNode1Hash] = HnewLeafNode1
						//connecting new leaf to branch
						newBranchNode.branch_value[currentNodeHexArray[0]] = newLeafNode1Hash
					}

					//for new leaf
					newLeafNode2Hash, HnewLeafNode2 := createLeafNode(pathLeft[1:], newValue)
					mpt.db[newLeafNode2Hash] = HnewLeafNode2

					//connecting two leaves to branch
					//newBranchNode.branch_value[currentNodeHexArray[0]] = newLeafNode1Hash
					newBranchNode.branch_value[pathLeft[0]] = newLeafNode2Hash

					newBranchNodeHash := newBranchNode.hash_node()
					mpt.db[newBranchNodeHash] = newBranchNode
					return newBranchNodeHash

				} else { //counter > 0  //create extension - branch - 2 leaves
					////fmt.Println("counter > 0, counter is - ", counter)
					newExtNode := Node{}
					newExtNode.node_type = 2
					////fmt.Println("common path : ", currentNodeHexArray[:counter])
					newExtNode.flag_value.encoded_prefix = compact_encode(currentNodeHexArray[:counter])

					//create new branch node
					newBranchNode := Node{}
					newBranchNode.node_type = 1

					//connect old leaf
					if len(currentNodeHexArray[counter:]) == 0 {
						newBranchNode.branch_value[16] = currentNode.flag_value.value
					} else {
						//fmt.Println("currentNodeHexArray[counter+1 :] - ", len(currentNodeHexArray[counter:]))

						//newLeafNode2Hash, newLeafNode2 := createLeafNode(currentNodeHexArray[counter+1:], currentNode.flag_value.value)
						newLeafNode2Path := currentNodeHexArray[counter+1:]
						newLeafNode2Hash, newLeafNode2 := createLeafNode(newLeafNode2Path, currentNode.flag_value.value)

						newBranchNode.branch_value[currentNodeHexArray[counter]] = newLeafNode2Hash

						mpt.db[newLeafNode2Hash] = newLeafNode2
					}
					//connect new insertion
					if len(pathLeft) >= counter+1 {
						////fmt.Println("len(pathLeft) >= counter+1 and pathLeft for leaf node : ", pathLeft[counter+1:])
						newLeafNode1Hash, newLeafNode1 := createLeafNode(pathLeft[counter+1:], newValue)
						newBranchNode.branch_value[pathLeft[counter]] = newLeafNode1Hash
						mpt.db[newLeafNode1Hash] = newLeafNode1
					} else if len(pathLeft) == counter { // len(pathleft) < len(currentnodepath)
						//											and pathLeft mathches completely
						//newvalue should then be in 16 index of branch
						newBranchNode.branch_value[16] = newValue
					}

					//hash of branch
					newBranchNodeHash := newBranchNode.hash_node()
					mpt.db[newBranchNodeHash] = newBranchNode

					//connecting branch to extension
					newExtNode.flag_value.value = newBranchNodeHash

					//hash of ext node
					newExtNodeHash := newExtNode.hash_node()
					mpt.db[newExtNodeHash] = newExtNode
					//del current node
					delete(mpt.db, currentNodeHash)

					return newExtNodeHash

				}
			}

		default:
			return "unable to insert"
		}

	}

	if len(pathLeft) == 0 {
		//TODO here //
		switch {
		case currentNode.node_type == 1: //branch
			currentNode.branch_value[16] = newValue
			currentNodeNewHash := currentNode.hash_node()
			mpt.db[currentNodeNewHash] = currentNode
			return currentNodeNewHash
		case currentNode.node_type == 2 && is_ext_node(currentNode.flag_value.encoded_prefix) == true: //ext

		case currentNode.node_type == 2 && is_ext_node(currentNode.flag_value.encoded_prefix) == false: //leaf
			//fmt.Println("pathleft is 0 and in leaf")
			currentNodeHexArray := compact_decode(currentNode.flag_value.encoded_prefix)
			if len(currentNodeHexArray) == 0 {
				//fmt.Println("0,leaf - current node hexarray ", currentNodeHexArray)
				currentNode.flag_value.value = newValue
				currentNodeNewHash := currentNode.hash_node()
				mpt.db[currentNodeNewHash] = currentNode
				return currentNodeNewHash
			} else {
				//create branch node and put newValue at 16th, and connect old trie
				newBranchNode := Node{}
				newBranchNode.node_type = 1

				//adding new value
				newBranchNode.branch_value[16] = newValue

				//creating new leaf for oldvalue
				newLeafNodeLeftPath := currentNodeHexArray[1:]
				newLeafNodeHash, newLeafNode := createLeafNode(newLeafNodeLeftPath, currentNode.flag_value.value)
				mpt.db[newLeafNodeHash] = newLeafNode

				//connecting new leaf for oldvalue
				newBranchNode.branch_value[currentNodeHexArray[0]] = newLeafNodeHash
				newBranchNodeHash := newBranchNode.hash_node()
				mpt.db[newBranchNodeHash] = newBranchNode

				delete(mpt.db, currentNodeHash)
				return newBranchNodeHash

			}

		default:
			return "pathLeft = 0 and unable to insert"
		}
	}

	return ""
}

func createLeafNode(pathLeft []uint8, newValue string) (string, Node) {
	newLeafNode := Node{}
	newLeafNode.node_type = 2
	newPathLeft := append(pathLeft, 16)
	newLeafNode.flag_value.encoded_prefix = compact_encode(newPathLeft)
	newLeafNode.flag_value.value = newValue
	newLeafNodeHash := newLeafNode.hash_node()
	return newLeafNodeHash, newLeafNode
}
