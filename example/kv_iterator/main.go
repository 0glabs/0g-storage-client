package main

import (
	"context"
	"fmt"

	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const KvClientAddr = "http://127.0.0.1:6789"

func main() {
	ctx := context.Background()
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
	err = iter.SeekToFirst(ctx)
	if err != nil {
		fmt.Printf("SeekToFirst error: %s", err)
	}
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		err = iter.Next(ctx)
		if err != nil {
			fmt.Printf("Next error: %s", err)
		}
	}

	fmt.Println("end to begin:")
	err = iter.SeekToLast(ctx)
	if err != nil {
		fmt.Printf("SeekToLast error: %s", err)
	}
	for iter.Valid() {
		pair := iter.KeyValue()
		fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
		err = iter.Prev(ctx)
		if err != nil {
			fmt.Printf("Prev error: %s", err)
		}
	}

	fmt.Printf("seek before %v\n", string(key1))
	err = iter.SeekBefore(ctx, key1)
	if err != nil {
		fmt.Printf("SeekBefore error: %s", err)
	}
	pair := iter.KeyValue()
	fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))

	fmt.Printf("seek after %v\n", string(key0))
	err = iter.SeekAfter(ctx, key0)
	if err != nil {
		fmt.Printf("SeekAfter error: %s", err)
	}
	pair = iter.KeyValue()
	fmt.Printf("%v: %v\n", string(pair.Key), string(pair.Data))
}
