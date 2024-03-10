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
	key := []byte("TESTKEY0")
	key1 := []byte("TESTKEY1")
	key2 := []byte("TESTKEY2")
	account := ethCommon.HexToAddress("0x578dd2bfc41bb66e9f0ae0802c613996440c9597")

	kvClient := kv.NewClient(client, nil)
	val, _ := kvClient.GetValue(streamId, key)
	fmt.Println(string(val.Data))
	val, _ = kvClient.GetValue(streamId, key1)
	fmt.Println(string(val.Data))
	val, _ = kvClient.GetValue(streamId, key2)
	fmt.Println(val)
	fmt.Println(kvClient.GetTransactionResult(2))
	fmt.Println(kvClient.GetHoldingStreamIds())
	fmt.Println(kvClient.HasWritePermission(account, streamId, key))
	fmt.Println(kvClient.IsAdmin(account, streamId))
	fmt.Println(kvClient.IsSpecialKey(streamId, key))
	fmt.Println(kvClient.IsWriterOfKey(account, streamId, key))
	fmt.Println(kvClient.IsWriterOfStream(account, streamId))
}
