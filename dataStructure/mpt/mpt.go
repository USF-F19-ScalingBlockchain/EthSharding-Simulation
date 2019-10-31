package mpt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

/*
	Data Structure to represent value of Leaf and Extension Node
*/
type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

/*
	Data Structure of Node that is inserted in MPT
*/
type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

/*
	Data Structure for MPT
*/
type MerklePatriciaTrie struct {
	db   map[string]Node
	Root string
}

/*
	Prints the hash map of MPT
*/
func print_hash_map(maps map[string]Node) {
	for k, v := range maps {
		if v.node_type == 2 {
			fmt.Println("Key: ", k, " Decoded Prefix: ", compact_decode(v.flag_value.encoded_prefix), " Node: ", v)
		} else {
			fmt.Println("Key: ", k, " Node: ", v)
		}
	}
}

/*
	Insert() function takes a pair of <key, value> as arguments.
	It will traverse down the Merkle Patricia Trie, find the right
	place to insert the value, and do the insertion.
*/
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	var node = mpt.insert(mpt.db[mpt.Root], convert_to_hex(key), new_value)
	mpt.Root = node.hash_node()
}

/*
	An helper function for insert which is called when path of
	inserting node completely matches the current path of
	leaf/ext node.
*/
func (mpt *MerklePatriciaTrie) fullmatch_leaf_ext(new_value string, cur_node *Node, rem_prefix []uint8) {
	//fmt.Println("Case 0")
	var node_type = mpt.get_node_type(*cur_node)
	if node_type == "Leaf" {
		cur_node.flag_value.value = new_value
	} else {
		var new_node = mpt.insert(mpt.db[cur_node.flag_value.value], rem_prefix, new_value)
		cur_node.flag_value.value = new_node.hash_node()
	}
}

/*
	An helper function for insert which is called when path of
	inserting node has some unmatched suffix left but current
	path is completely matched with the prefix.
*/
func (mpt *MerklePatriciaTrie) fullmatch_cur_key(new_value string, cur_node *Node, rem_prefix []uint8, cur_prefix []uint8) {
	//fmt.Println("Case 1")
	var node_type = mpt.get_node_type(*cur_node)
	if node_type == "Extension" {
		var new_node = mpt.insert(mpt.db[cur_node.flag_value.value], rem_prefix, new_value)
		cur_node.flag_value.value = new_node.hash_node()
	} else {
		var new_node = Node{1, [17]string{""}, Flag_value{}}
		new_node.branch_value[16] = cur_node.flag_value.value
		var new_leaf_node = Node{2, [17]string{""}, Flag_value{compact_encode(append(rem_prefix[1:], 16)), new_value}}
		new_node.branch_value[rem_prefix[0]] = new_leaf_node.hash_node()
		mpt.db[new_leaf_node.hash_node()] = new_leaf_node
		mpt.db[new_node.hash_node()] = new_node
		if len(cur_prefix) != 0 {
			cur_node.flag_value.value = new_node.hash_node()
			cur_node.flag_value.encoded_prefix = compact_encode(cur_prefix)
			mpt.db[cur_node.hash_node()] = *cur_node
		} else {
			*cur_node = new_node
		}
	}
}

/*
	An helper function for insert which is called when path of
	inserting node completely matches with the prefix of current
	node but current node has some unmatched suffix left.
*/
func (mpt *MerklePatriciaTrie) fullmatch_prefix_key(new_value string, cur_node *Node, rem_cur_prefix []uint8, cur_prefix []uint8, i int) Node {
	//fmt.Println("Case 2")
	var node_type = mpt.get_node_type(*cur_node)
	var new_node = Node{1, [17]string{""}, Flag_value{}}
	new_node.branch_value[16] = new_value
	var new_cur_node Node
	if node_type == "Leaf" {
		new_cur_node = Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(append(rem_cur_prefix[1:], 16)), cur_node.flag_value.value}}
	} else {
		new_cur_node = Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(rem_cur_prefix[1:]), cur_node.flag_value.value}}
	}
	new_node.branch_value[rem_cur_prefix[0]] = new_cur_node.hash_node()
	if i != 0 {
		cur_node = &Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(cur_prefix[:i]), new_node.hash_node()}}
		mpt.db[cur_node.hash_node()] = *cur_node
		mpt.db[new_node.hash_node()] = new_node
		mpt.db[new_cur_node.hash_node()] = new_cur_node
		return *cur_node
	} else {
		new_node.branch_value[cur_prefix[0]] = cur_node.hash_node()
		mpt.db[new_node.hash_node()] = new_node
		mpt.db[new_cur_node.hash_node()] = new_cur_node
		return new_node
	}
}

/*
	An helper function for insert which is called when path of
	inserting node either don't match completely with the path
	of current node or partially matches the path of current
	node.
*/
func (mpt *MerklePatriciaTrie) no_match(new_value string, cur_node *Node, rem_cur_prefix []uint8, cur_prefix []uint8, rem_prefix []uint8, i int) Node {
	//fmt.Println("Case 3")
	var node_type = mpt.get_node_type(*cur_node)
	var new_node = Node{1, [17]string{""}, Flag_value{}}
	var new_leaf_node = Node{2, [17]string{""}, Flag_value{compact_encode(append(rem_prefix[1:], 16)), new_value}}
	new_node.branch_value[rem_prefix[0]] = new_leaf_node.hash_node()
	var new_cur_node Node
	if node_type == "Leaf" {
		new_cur_node = Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(append(rem_cur_prefix[1:], 16)), cur_node.flag_value.value}}
	} else {
		if len(rem_cur_prefix) == 1 {
			new_cur_node = mpt.db[cur_node.flag_value.value]
			delete(mpt.db, cur_node.hash_node())
		} else {
			new_cur_node = Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(rem_cur_prefix[1:]), cur_node.flag_value.value}}
		}
	}
	new_node.branch_value[rem_cur_prefix[0]] = new_cur_node.hash_node()
	mpt.db[new_node.hash_node()] = new_node
	mpt.db[new_leaf_node.hash_node()] = new_leaf_node
	mpt.db[new_cur_node.hash_node()] = new_cur_node
	if i == 0 {
		return new_node
	} else {
		cur_node = &Node{cur_node.node_type, cur_node.branch_value, Flag_value{compact_encode(cur_prefix[:i]), new_node.hash_node()}}
		mpt.db[cur_node.hash_node()] = *cur_node
		return *cur_node
	}
}

/*
	Main insert helper function which takes cur_node,
	prefix (i.e. path to insert new node) new_value
	which is value that we need to insert.
*/
func (mpt *MerklePatriciaTrie) insert(cur_node Node, prefix []uint8, new_value string) Node {
	var old_hash = cur_node.hash_node()
	if cur_node.node_type == 0 || mpt.Root == "" {
		cur_node = Node{2, [17]string{""}, Flag_value{compact_encode(append(prefix, 16)), new_value}}
	} else if cur_node.node_type == 1 {
		if len(prefix) == 0 {
			cur_node.branch_value[16] = new_value
		} else {
			var new_node = mpt.insert(mpt.db[cur_node.branch_value[prefix[0]]], prefix[1:], new_value)
			cur_node.branch_value[prefix[0]] = new_node.hash_node()
		}
	} else {
		var cur_prefix = compact_decode(cur_node.flag_value.encoded_prefix)
		var i = 0
		for ; i < len(prefix) && i < len(cur_prefix); i++ {
			if prefix[i] != cur_prefix[i] {
				break
			}
		}
		var rem_prefix = prefix[i:]
		var rem_cur_prefix = cur_prefix[i:]
		if len(rem_prefix) == 0 && len(rem_cur_prefix) == 0 { // full match leaf/ext node.
			mpt.fullmatch_leaf_ext(new_value, &cur_node, rem_prefix)
		} else if len(rem_cur_prefix) == 0 { // full match cur_key
			mpt.fullmatch_cur_key(new_value, &cur_node, rem_prefix, cur_prefix)
		} else if len(rem_prefix) == 0 { // full match prefix_key
			return mpt.fullmatch_prefix_key(new_value, &cur_node, rem_cur_prefix, cur_prefix, i)
		} else { // No match
			return mpt.no_match(new_value, &cur_node, rem_cur_prefix, cur_prefix, rem_prefix, i)
		}
	}
	if old_hash != cur_node.hash_node() {
		delete(mpt.db, old_hash)
		mpt.db[cur_node.hash_node()] = cur_node
	}
	return cur_node
}

/*
	The function returns the node type for
	leaf and extension nodes.
*/
func (mpt *MerklePatriciaTrie) get_node_type(node Node) string {
	if node.node_type == 2 {
		var decoded = compact_decode_helper(node.flag_value.encoded_prefix)
		if decoded[0] == 1 || decoded[0] == 0 {
			return "Extension"
		} else if decoded[0] == 2 || decoded[0] == 3 {
			return "Leaf"
		}
		return ""
	}
	return ""
}

/*
	The function converts string to bytes and
	passes it to compact_decode_helper
*/
func convert_to_hex(str string) []uint8 {
	return compact_decode_helper([]byte(str))
}

/*
	Get() function takes a key as the argument, traverse
	down the Merkle Patricia Trie and find the value.
	If the key doesn't exist, it will return an empty string
*/
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	var root = mpt.db[mpt.Root]
	var prefix = convert_to_hex(key)
	return mpt.get(prefix, root)
}

/*
	Get helper method that takes prefix (i.e. path of
	node that we are searching) and current node to
	traverse down the tree. It returns and error if
	path is not found.
*/
func (mpt *MerklePatriciaTrie) get(prefix []uint8, cur_node Node) (string, error) {
	if cur_node.node_type == 0 || mpt.Root == "" {
		return "", errors.New("path_not_found")
	} else {
		if cur_node.node_type == 1 { // Branch Node
			if len(prefix) == 0 {
				return cur_node.branch_value[16], nil
			}
			var new_node = mpt.db[cur_node.branch_value[prefix[0]]]
			return mpt.get(prefix[1:], new_node)
		} else { // Ext or Leaf Node
			var cur_prefix = compact_decode(cur_node.flag_value.encoded_prefix)
			var i = 0
			var node_type = mpt.get_node_type(cur_node)
			for ; i < len(prefix) && i < len(cur_prefix); i++ {
				if prefix[i] != cur_prefix[i] {
					break
				}
			}
			if i == len(prefix) && i == len(cur_prefix) && node_type == "Leaf" {
				return cur_node.flag_value.value, nil
			} else if i != len(cur_prefix) {
				return "", errors.New("path_not_found")
			}
			var new_node = mpt.db[cur_node.flag_value.value]
			return mpt.get(prefix[i:], new_node)
		}
	}
}

/*
	Delete() function takes a key as the argument, traverse
	down the Merkle Patricia Trie and find that key. If the
	key exists, delete the corresponding value and re-balance
	the trie if necessary; if the key doesn't exist,
	return "path_not_found".
*/
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	var old_hash = mpt.Root
	var node = mpt.delete(convert_to_hex(key), mpt.db[mpt.Root])
	if node.hash_node() == old_hash {
		return "path_not_found"
	}
	mpt.Root = node.hash_node()
	return ""
}

/*
	Delete helper function which takes prefix (i.e. path
	of the node that needs to be deleted) and current
	node to traverse down the tree to delete a node.
*/
func (mpt *MerklePatriciaTrie) delete(prefix []uint8, cur_node Node) Node {
	var old_hash = cur_node.hash_node()
	if cur_node.node_type == 1 {
		if len(prefix) == 0 {
			cur_node.branch_value[16] = ""
			cur_node = mpt.normalize_branch(cur_node)
		} else {
			var next_node = mpt.delete(prefix[1:], mpt.db[cur_node.branch_value[prefix[0]]])
			if next_node.hash_node() != cur_node.branch_value[prefix[0]] {
				if next_node.node_type == 0 {
					cur_node.branch_value[prefix[0]] = ""
					cur_node = mpt.normalize_branch(cur_node)
				} else {
					cur_node.branch_value[prefix[0]] = next_node.hash_node()
				}
			}
		}
	} else {
		var cur_prefix = compact_decode(cur_node.flag_value.encoded_prefix)
		var node_type = mpt.get_node_type(cur_node)
		var i = 0
		for ; i < len(prefix) && i < len(cur_prefix); i++ {
			if prefix[i] != cur_prefix[i] {
				break
			}
		}
		if i == len(prefix) && i == len(cur_prefix) && node_type == "Leaf" {
			cur_node = Node{}
		} else if i == len(cur_prefix) && node_type == "Extension" {
			mpt.delete_ext(&cur_node, prefix, cur_prefix, i)
		}
	}
	if old_hash != cur_node.hash_node() {
		delete(mpt.db, old_hash)
		if cur_node.node_type != 0 {
			mpt.db[cur_node.hash_node()] = cur_node
		}
	}
	return cur_node
}

/*
	Delete helper when node type is extension.
*/
func (mpt *MerklePatriciaTrie) delete_ext(cur_node *Node, prefix []uint8, cur_prefix []uint8, i int) {
	var next_node = mpt.delete(prefix[i:], mpt.db[cur_node.flag_value.value])
	if next_node.hash_node() != cur_node.flag_value.value {
		if next_node.node_type == 0 {
			*cur_node = Node{}
		} else if next_node.node_type == 1 {
			cur_node.flag_value.value = next_node.hash_node()
		} else {
			var node_type = mpt.get_node_type(next_node)
			var decoded_prefix = compact_decode(next_node.flag_value.encoded_prefix)
			if node_type == "Leaf" {
				decoded_prefix = append(decoded_prefix, 16)
			}
			cur_node.flag_value.encoded_prefix = compact_encode(append(cur_prefix, decoded_prefix...))
			cur_node.flag_value.value = next_node.flag_value.value
		}
	}
}

/*
	This function will normalize the branch node with just one value.
*/
func (mpt *MerklePatriciaTrie) normalize_branch(cur_node Node) Node {
	var sum = 0
	var index = 0
	for i := 0; i < len(cur_node.branch_value); i++ {
		if cur_node.branch_value[i] != "" {
			sum += 1
			index = i
		}
	}
	if sum > 1 {
		return cur_node
	} else if index == 16 {
		var node = Node{2, [17]string{""}, Flag_value{compact_encode([]uint8{16}), cur_node.branch_value[16]}}
		return node
	} else {
		var next_node = mpt.db[cur_node.branch_value[index]]
		if next_node.node_type == 2 {
			var decode_prefix = compact_decode(next_node.flag_value.encoded_prefix)
			decode_prefix = append([]uint8{uint8(index)}, decode_prefix...)
			if mpt.get_node_type(next_node) == "Leaf" {
				next_node.flag_value.encoded_prefix = compact_encode(append(decode_prefix, 16))
				return next_node
			}
			next_node.flag_value.encoded_prefix = compact_encode(decode_prefix)
			return next_node
		} else {
			cur_node.branch_value = [17]string{""}
			cur_node.node_type = 2
			cur_node.flag_value.encoded_prefix = compact_encode([]uint8{uint8(index)})
			cur_node.flag_value.value = next_node.hash_node()
			return cur_node
		}
	}
}

/*
	This function takes an array of HEX value as the input,
	mark the Node type(such as Branch, Leaf, Extension),
	make sure the length is even, and convert it into array
	of ASCII number as the output.
*/
func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if len(hex_array) != 0 && hex_array[len(hex_array)-1] == 16 {
		term = 1
	} else {
		term = 0
	}
	if term == 1 {
		hex_array = hex_array[:len(hex_array)-1]
	}
	var odd = len(hex_array) % 2
	var flag = uint8(2*term + odd)
	if odd == 1 {
		hex_array = append([]uint8{flag}, hex_array...)
	} else {
		hex_array = append([]uint8{flag, 0}, hex_array...)
	}
	var res []uint8
	for i := 0; i < len(hex_array); i += 2 {
		res = append(res, 16*hex_array[i]+hex_array[i+1])
	}
	return res
}

/*
	This function is helper function for compact_decode.
*/
func compact_decode_helper(encoded_arr []uint8) []uint8 {
	var res []uint8
	for i := 0; i < len(encoded_arr); i++ {
		var d1 = uint8(encoded_arr[i] % 16)
		var d2 = uint8((encoded_arr[i]) / 16)
		res = append(res, d2)
		res = append(res, d1)
	}
	return res
}

/*
	This function reverses the compact_encode() function.
*/
func compact_decode(encoded_arr []uint8) []uint8 {
	var res = compact_decode_helper(encoded_arr)
	if len(res) != 0 {
		if res[0] == 2 || res[0] == 0 {
			res = res[2:]
		} else {
			res = res[1:]
		}
	}
	return res
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

/*
	UDF for testing MPT
*/
func Tests() {
	var mpt = MerklePatriciaTrie{map[string]Node{}, ""}
	mpt.Insert("a", "10")
	mpt.Insert("p", "20")
	mpt.Insert("ab", "30")
	fmt.Println("Root: ", mpt.Root)
	print_hash_map(mpt.db)
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	var sum = sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.Root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}

// GetAllKeyValuePairs of mpt and put in map
func (mpt *MerklePatriciaTrie) GetAllKeyValuePairs() map[string]string {

	if len(mpt.db) == 0 {
		return nil //, errors.New("Empty MPT")
	}
	emptyKeyValuePairs := make(map[string]string)
	rootNode := mpt.db[mpt.Root]

	KeyValuePairs, err := mpt.GetAllKeyValuePairsHelper(emptyKeyValuePairs, rootNode, []uint8{})
	if err != nil {
		return emptyKeyValuePairs
	} else {
		return KeyValuePairs
	}

}

func (mpt *MerklePatriciaTrie) GetAllKeyValuePairsHelper(mptKeyValuePairs map[string]string, thisNode Node, hexPath []uint8) (map[string]string, error) {
	currentHexPath := hexPath

	switch {
	case thisNode.node_type == 1:
		for i := 0; i < 16; i++ {
			if thisNode.branch_value[i] != "" {
				newcurrentHexPath := append(currentHexPath, uint8(i)) //int should be treated as part of ascii path
				mpt.GetAllKeyValuePairsHelper(mptKeyValuePairs, mpt.db[thisNode.branch_value[i]], newcurrentHexPath)
			}
		}
		if thisNode.branch_value[16] != "" {
			key := HexArraytoString(currentHexPath)
			mptKeyValuePairs[key] = thisNode.branch_value[16]
		}

	case thisNode.node_type == 2 && is_ext_node(thisNode.flag_value.encoded_prefix) == true:
		thisNodePath := compact_decode(thisNode.flag_value.encoded_prefix)
		currentHexPath := append(currentHexPath, thisNodePath...) //int should be treated as part of ascii path
		mpt.GetAllKeyValuePairsHelper(mptKeyValuePairs, mpt.db[thisNode.flag_value.value], currentHexPath)

	case thisNode.node_type == 2 && is_ext_node(thisNode.flag_value.encoded_prefix) == false:
		thisNodePath := compact_decode(thisNode.flag_value.encoded_prefix)
		currentHexPath := append(currentHexPath, thisNodePath...)
		key := HexArraytoString(currentHexPath)
		mptKeyValuePairs[key] = thisNode.flag_value.value
	default:
		return nil, errors.New("Error in contructing key Value map from MPT")

	}

	return mptKeyValuePairs, nil

}

func HexArraytoString(hexArray []uint8) string {
	asciiPath := []uint8{}
	for i := 0; i < len(hexArray)-1; i = i + 2 {
		asciiPath = append(asciiPath, 16*hexArray[i]+hexArray[i+1])
	}
	return string(asciiPath)
}
