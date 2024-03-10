package main

import (
	"fmt"

	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const KvClientAddr = "http://127.0.0.1:6789"

func main() {
	client, err := node.NewClient(KvClientAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	streamId := ethCommon.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd")
	key0 := []byte("TESTKEY")
	key1 := []byte("TESTKEY2")

	kvClient := kv.NewClient(client, nil)
	iter := kvClient.NewIterator(streamId)

	fmt.Println("begin to end:")
	iter.SeekToFirst()
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Next()
	}

	fmt.Println("end to begin:")
	iter.SeekToLast()
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		iter.Prev()
	}

	fmt.Printf("seek before %v\n", string(key1))
	iter.SeekBefore(key1)
	pair := iter.KeyValue()
	fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))

	fmt.Printf("seek after %v\n", string(key0))
	iter.SeekAfter(key0)
	pair = iter.KeyValue()
	fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
}
