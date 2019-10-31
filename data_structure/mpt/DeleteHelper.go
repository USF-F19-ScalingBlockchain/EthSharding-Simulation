package mpt

import (
	"errors"
	"fmt"
	"reflect"
)

//Delete func
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	pathLeft := StringToHexArray(key)
	currentNode := mpt.db[mpt.Root]

	value, err := mpt.delHelper(mpt.Root, currentNode, pathLeft, []string{})
	fmt.Println("Key - Value deleted : ", key, " - ", value)
	if err != nil {
		return "path_not_found"
	}
	return value
}

func (mpt *MerklePatriciaTrie) delHelper(nodeKey string, currentNode Node, pathLeft []uint8, hashStack []string) (string, error) {

	if len(pathLeft) > 0 && currentNode.node_type != 0 { //path length >0

		if currentNode.node_type == 1 { // branch and pathleft >0

			if currentNode.branch_value[pathLeft[0]] != "" {

				hash := currentNode.branch_value[pathLeft[0]]

				oldhash := currentNode.hash_node()
				index := pathLeft[0]
				pathLeft = pathLeft[1:]

				nextNode := mpt.db[hash]

				if nextNode.node_type == 2 {

					currentNodeHexArray := AsciiArrayToHexArray(nextNode.flag_value.encoded_prefix)
					if (currentNodeHexArray[0] == 2) || (currentNodeHexArray[0] == 3) { //leaf

						if len(pathLeft) == 0 && len(compact_decode(nextNode.flag_value.encoded_prefix)) == 0 {

							currentNode.branch_value[index] = ""
							mpt.db[oldhash] = currentNode

						} else if len(pathLeft) > 0 && len(nextNode.flag_value.encoded_prefix) > 0 {

							if reflect.DeepEqual(pathLeft, compact_decode(nextNode.flag_value.encoded_prefix)) {

								currentNode.branch_value[index] = ""
								mpt.db[oldhash] = currentNode

							}
						}
					}
				}

				hashStack = append(hashStack, oldhash)
				return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)

			}

		} else if currentNode.node_type == 2 { //ext or leaf and pathleft >0
			currentNodeHexArray := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)

			oldhash := currentNode.hash_node()

			if (currentNodeHexArray[0] == 0) || (currentNodeHexArray[0] == 1) { //extension

				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)

				if reflect.DeepEqual(nodePath, pathLeft[:len(nodePath)]) {
					pathLeft = pathLeft[len(nodePath):]

					if len(pathLeft) == 0 { // pathleft is zero now

						hash := currentNode.flag_value.value

						hashStack = append(hashStack, oldhash)

						if mpt.db[hash].branch_value[16] != "" { //value found

							nextnode := mpt.db[hash]
							valuetoreturn := nextnode.branch_value[16]
							nextnode.branch_value[16] = ""
							mpt.db[hash] = nextnode

							hashStack = append(hashStack, hash)

							mpt.rearrangeDeletedTrie(hashStack)

							return valuetoreturn, nil
						}

					} else if len(pathLeft) > 0 { //add to hashstack call on next node

						hash := currentNode.flag_value.value

						hashStack = append(hashStack, oldhash) //adding current ext

						return mpt.delHelper(hash, mpt.db[hash], pathLeft, hashStack)
					}

				}

			} else if (currentNodeHexArray[0] == 2) || (currentNodeHexArray[0] == 3) { //leaf //pathleft >0

				// hex without prefix <(ascii with prefix)
				nodePath := compact_decode(currentNode.flag_value.encoded_prefix)

				if reflect.DeepEqual(nodePath, pathLeft) {

					delete(mpt.db, nodeKey)

					mpt.rearrangeDeletedTrie(hashStack)

					return currentNode.flag_value.value, nil
				}

			}
		}
	} else if len(pathLeft) == 0 && currentNode.node_type != 0 { //pathlength ==0

		if currentNode.node_type == 1 { // branch
			if currentNode.branch_value[16] != "" {

				previousHash := currentNode.hash_node()
				previousValue := currentNode.branch_value[16]
				currentNode.branch_value[16] = ""
				mpt.db[previousHash] = currentNode
				hashStack = append(hashStack, previousHash) //adding current ext
				mpt.rearrangeDeletedTrie(hashStack)
				return previousValue, nil
			}

		} else if currentNode.node_type == 2 { //ext or leaf
			//extension

			nodeHexPath := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)

			if nodeHexPath[0] == 2 || nodeHexPath[0] == 3 { //leaf
				nodeHexPath := AsciiArrayToHexArray(currentNode.flag_value.encoded_prefix)

				if reflect.DeepEqual(nodeHexPath, []uint8{2, 0}) {

					value := currentNode.flag_value.value
					delete(mpt.db, currentNode.hash_node())
					mpt.rearrangeDeletedTrie(hashStack)
					return value, nil
				}

			}
		}
	}

	return "", errors.New("path_not_found")
}

// rearrangeDeletedTrie rearranges mpt
func (mpt *MerklePatriciaTrie) rearrangeDeletedTrie(hashStack []string) {

	counter := len(hashStack) - 1
	//rearranged :=
	mpt.rearrangeDeletedTrieHelper(hashStack, counter, "")

}

func (mpt *MerklePatriciaTrie) rearrangeDeletedTrieHelper(hashStack []string, counter int, currenthash string) bool {

	if counter == -1 {
		mpt.Root = currenthash
		return true
	}

	if counter == len(hashStack)-1 {

		if len(hashStack) == 1 {
			if mpt.db[hashStack[counter]].node_type == 1 { //just one node in stack and should be a branch
				currNode := mpt.db[hashStack[counter]]
				numValues := 0
				n := 0 //mpt.db[currNode.branch_value]
				for i := 0; i < 17; i++ {
					if currNode.branch_value[i] != "" {
						numValues++
						if numValues == 1 {
							n = i
							// 	nextnode := mpt.db[currNode.branch_value[i]]
						}
					}
				}
				if numValues == 1 {
					nextnode := mpt.db[currNode.branch_value[n]]
					nodetype := nextnode.node_type
					if nodetype == 1 {
						nodeE := Node{}
						nodeE.node_type = 2
						nodeE.flag_value.encoded_prefix = compact_encode([]uint8{uint8(n)})
						delete(mpt.db, hashStack[counter])
						hash := nodeE.hash_node()
						mpt.db[hash] = nodeE
						return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
					} else if nodetype == 2 {
						hex_prefix_array := AsciiArrayToHexArray(nextnode.flag_value.encoded_prefix)

						decoded_hex_array := compact_decode(nextnode.flag_value.encoded_prefix)

						if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
							nodeE := Node{}
							nodeE.node_type = 2
							combinedHexArray := append([]uint8{uint8(n)}, compact_decode(nextnode.flag_value.encoded_prefix)...)
							nodeE.flag_value.encoded_prefix = compact_encode(combinedHexArray)
							nodeE.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							delete(mpt.db, hashStack[counter-1])
							hash := nodeE.hash_node()
							mpt.db[hash] = nodeE
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
						} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
							nodeL := Node{}
							nodeL.node_type = 2
							decoded_hex_array = append([]uint8{uint8(n)}, decoded_hex_array...)
							decoded_hex_array = append(decoded_hex_array, 16)
							nodeL.flag_value.encoded_prefix = compact_encode(decoded_hex_array)
							nodeL.flag_value.value = nextnode.flag_value.value
							delete(mpt.db, hashStack[counter])
							hash := nodeL.hash_node()
							mpt.db[hash] = nodeL
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, nodeL.hash_node())
						}
					}

				} else if numValues > 1 {
					node := mpt.db[hashStack[counter]]
					hash := node.hash_node()
					mpt.db[hash] = node
					return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
				}
			}

		} else if len(hashStack) > 1 {
			if mpt.db[hashStack[counter]].node_type == 1 { // branch
				currNode := mpt.db[hashStack[counter]]
				numValues := 0
				n := 0
				for i := 0; i < 17; i++ {
					if currNode.branch_value[i] != "" {
						numValues++
						if numValues == 1 {
							n = i
						}
					}
				}
				if numValues == 1 {

					if n == 16 {
						//fmt.Println("It is the branchvalue[16] field n=16")
						//convert to leaf and store value with key = ""
						//put the value in branch[16] of the following branch
						prevnode := mpt.db[hashStack[counter-1]]
						prevnodetype := prevnode.node_type
						if prevnodetype == 1 { // branch
							// create a empty(key) leaf node with value (branch_value[16])

							nodeL := Node{}
							nodeL.node_type = 2
							nodeL.flag_value.encoded_prefix = compact_encode([]uint8{16})
							nodeL.flag_value.value = currNode.branch_value[16]
							delete(mpt.db, hashStack[counter])
							hash := nodeL.hash_node()
							mpt.db[hash] = nodeL
							return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
						} else if prevnodetype == 2 {

							hex_prefix_array := AsciiArrayToHexArray(prevnode.flag_value.encoded_prefix)

							if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension
								// add the value to prev Ext node - and make a leaf
								hex_array := compact_decode(prevnode.flag_value.encoded_prefix)
								hex_array = append(hex_array, 16)
								nodeL := Node{}
								nodeL.node_type = 2
								nodeL.flag_value.encoded_prefix = compact_encode(hex_array)
								nodeL.flag_value.value = currNode.branch_value[16]
								delete(mpt.db, hashStack[counter])
								delete(mpt.db, hashStack[counter-1])
								hash := nodeL.hash_node()
								mpt.db[hash] = nodeL
								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
							}
						}

					} else if n < 16 && n >= 0 {
						nextnode := mpt.db[currNode.branch_value[n]]
						nodetype := nextnode.node_type
						u := uint8(n)

						if nodetype == 1 {
							//Creating Extension node
							var new_array []uint8
							var flag int
							previous_node := mpt.db[hashStack[counter-1]]
							if previous_node.node_type == 2 {
								hex_with_prefix_array := AsciiArrayToHexArray(previous_node.flag_value.encoded_prefix)
								if hex_with_prefix_array[0] == 0 || hex_with_prefix_array[0] == 1 {

									flag = 1
									new_array = compact_decode(previous_node.flag_value.encoded_prefix)
								}
							}
							new_array = append(new_array, []uint8{u}...)
							nodeE := Node{}
							nodeE.node_type = 2
							nodeE.flag_value.encoded_prefix = compact_encode(new_array)
							nodeE.flag_value.value = currNode.branch_value[n]
							delete(mpt.db, hashStack[counter])
							hash := nodeE.hash_node()
							mpt.db[hash] = nodeE
							if flag == 1 {
								delete(mpt.db, hashStack[counter-1])

								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
							} else {

								return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
							}
						} else if nodetype == 2 {

							hex_prefix_array := AsciiArrayToHexArray(nextnode.flag_value.encoded_prefix)

							if (hex_prefix_array[0] == 0) || (hex_prefix_array[0] == 1) { //extension

								var new_array []uint8
								previous_node := mpt.db[hashStack[counter-1]]
								var flag int
								if previous_node.node_type == 2 {
									hex_with_prefix_array := AsciiArrayToHexArray(previous_node.flag_value.encoded_prefix)
									if hex_with_prefix_array[0] == 0 || hex_with_prefix_array[0] == 1 { // extension
										flag = 1

										new_array = compact_decode(previous_node.flag_value.encoded_prefix)
									}
								}

								new_array = append(new_array, []uint8{u}...)

								hex_array := compact_decode(nextnode.flag_value.encoded_prefix)
								new_array = append(new_array, hex_array...)
								next_hex := currNode.branch_value[n]

								nextnode.flag_value.encoded_prefix = compact_encode(new_array)
								delete(mpt.db, hashStack[counter])
								delete(mpt.db, next_hex)

								hash := nextnode.hash_node()
								mpt.db[hash] = nextnode
								if flag == 1 {
									delete(mpt.db, hashStack[counter-1])
									return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
								} else {
									return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
								}
							} else if (hex_prefix_array[0] == 2) || (hex_prefix_array[0] == 3) { //leaf
								//check if the previous is extension - if so club it with the leaf
								var flag int
								var hex_array []uint8
								previousNode := mpt.db[hashStack[counter-1]]
								if previousNode.node_type == 2 {
									hex_with_prefix_array := AsciiArrayToHexArray(previousNode.flag_value.encoded_prefix)
									if hex_with_prefix_array[0] == 0 || hex_with_prefix_array[0] == 1 {
										flag = 1
										hex_array = compact_decode(previousNode.flag_value.encoded_prefix)
									}
								}

								hex_array = append(hex_array, []uint8{u}...)
								to_be_added_array := compact_decode(nextnode.flag_value.encoded_prefix)
								hex_array = append(hex_array, to_be_added_array...)
								hex_array = append(hex_array, 16)

								nextnode.flag_value.encoded_prefix = compact_encode(hex_array)
								delete(mpt.db, hashStack[counter])
								hash := nextnode.hash_node()
								mpt.db[hash] = nextnode
								if flag == 1 {
									delete(mpt.db, hashStack[counter-1])
									return mpt.rearrangeDeletedTrieHelper(hashStack, counter-2, hash)
								} else {
									return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
								}
							}
						}
					}
				} else if numValues > 1 {
					node := mpt.db[hashStack[counter]]

					hash := node.hash_node()

					delete(mpt.db, hashStack[counter])

					counter = counter - 1
					mpt.db[hash] = node
					return mpt.rearrangeDeletedTrieHelper(hashStack, counter, hash)
				}
			}
		}
	} else {

		node := mpt.db[hashStack[counter]]
		if node.node_type == 1 {

			for i := 0; i < 17; i++ {
				if node.branch_value[i] == hashStack[counter+1] {
					node.branch_value[i] = currenthash
				}
			}
			delete(mpt.db, hashStack[counter])
			hash := node.hash_node()
			mpt.db[hash] = node
			return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
		} else if node.node_type == 2 {

			node.flag_value.value = currenthash
			delete(mpt.db, hashStack[counter])
			hash := node.hash_node()
			mpt.db[hash] = node
			return mpt.rearrangeDeletedTrieHelper(hashStack, counter-1, hash)
		}

	}

	return false
}
