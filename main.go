package main

import "fmt"

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

func main() {
	fmt.Println("start ..")
	var s string
	s = SHARD.String()
	fmt.Println(s)
}
