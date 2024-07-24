package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	fmt.Printf("key: %v\n", args[0])
	fmt.Printf("contract: %v\n", args[1])
	fmt.Printf("chain node: %v\n", args[2])
	fmt.Printf("zgs node: %v\n", args[3])
	fmt.Printf("indexer: %v\n", args[4])
}
