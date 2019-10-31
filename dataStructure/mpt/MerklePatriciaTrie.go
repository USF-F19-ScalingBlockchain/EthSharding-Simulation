package mpt

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
)

// MerklePatriciaTrie contains vars db and Root
// db is a map where key is id of the node and Value is the node itself
// Root is the Node id for the root of the trie
type MerklePatriciaTrie struct {
	db   map[string]Node
	Root string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

// initializes the MerklePatriciaTrie
func (mpt *MerklePatriciaTrie) Initial() {

	mpt.db = make(map[string]Node)
	mpt.Root = ""
}

//
func compact_encode(hex_array []uint8) []uint8 {
	var term int
	if hex_array[len(hex_array)-1] == 16 { //checking if last element in array is 16
		term = 1 // if last element is 16, term = 1 and remove the last element from hex_array
		hex_array = hex_array[:len(hex_array)-1]
	}

	oddlen := len(hex_array) % 2 //checking if length is odd (oddlen = 1) otherwise (oddlen = 0)

	flags := []uint8{(uint8(2*term + oddlen))} //calculating flags value

	//changing hex_array based on value of oddlen
	if oddlen == 1 { //odd       // prefix -> flags so... ( either, 1 - Odd length extension
		hex_array = append(flags, hex_array...) //      3 - Odd length leaf        )
	} else { //when oddlen = 0 // even
		flags = append(flags, uint8(0))
		hex_array = append(flags, hex_array...) //prefix -> flags + 0 so...  (either, 00 -  Even Length Extention)
	} // or 20 - Even length Leaf

	encoded_arr := []uint8{} //array to return
	for i := 0; i < len(hex_array); i += 2 {
		encoded_arr = append(encoded_arr, 16*hex_array[i]+hex_array[i+1])
	}

	return encoded_arr
} // closing compact_encode

//
func compact_decode(encoded_arr []uint8) []uint8 {
	decoded_arr := AsciiArrayToHexArray(encoded_arr) //converting ascii array to hex array

	prefix := decoded_arr[0]
	switch prefix {
	case 0:
		decoded_arr = decoded_arr[2:]
	case 1:
		decoded_arr = decoded_arr[1:]
	case 2:
		decoded_arr = decoded_arr[2:]
	case 3:
		decoded_arr = decoded_arr[1:]
	}
	return decoded_arr
} //closing compact_decode

// AsciiArrayToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func AsciiArrayToHexArray(encoded_arr []uint8) []uint8 {
	hexArray := []uint8{}
	for _, element := range encoded_arr {
		div, mod := element/16, element%16
		hexArray = append(hexArray, div)
		hexArray = append(hexArray, mod)
	}
	return hexArray
}

////
//func test_compact_encode() {
//	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
//	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
//	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
//	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
//}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

// StringToHexArray takes in string as input and returns a splitted-hex(0-15) array : 1st
func StringToHexArray(s string) []uint8 {
	uint8Array := []uint8(s)
	hexArray := AsciiArrayToHexArray(uint8Array)
	return hexArray
}

//
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
		str = node.flag_value.value + string(node.flag_value.encoded_prefix)
	}
	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

//
func HexArraytoString(hexArray []uint8) string {
	asciiPath := []uint8{}
	for i := 0; i < len(hexArray)-1; i = i + 2 {
		asciiPath = append(asciiPath, 16*hexArray[i]+hexArray[i+1])
	}
	return string(asciiPath)
}

func Show() {
	fmt.Println("Shows MPT data")
}

//testing compact_decode and compact_encode
func test_compact_decode_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func TestCompact() {
	test_compact_decode_encode()
}
