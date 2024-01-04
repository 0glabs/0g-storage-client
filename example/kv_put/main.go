package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/kv"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

const ZgsClientAddr = "http://127.0.0.1:5678"
const BlockchainClientAddr = ""
const PrivKey = ""
const FlowContractAddr = ""

func main() {
	zgsClient, err := node.NewClient(ZgsClientAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	blockchainClient := blockchain.MustNewWeb3(BlockchainClientAddr, PrivKey)
	defer blockchainClient.Close()
	blockchain.CustomGasLimit = 1000000
	zgs, err := contract.NewFlowContract(common.HexToAddress(FlowContractAddr), blockchainClient)
	if err != nil {
		fmt.Println(err)
		return
	}
	kvClient := kv.NewClient(zgsClient, zgs)
	batcher := kvClient.Batcher()
	batcher.Set(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd"),
		[]byte("TESTKEY0"),
		[]byte{69, 70, 71, 72, 73},
	)
	batcher.Set(common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000f2bd"),
		[]byte("TESTKEY1"),
		[]byte{74, 75, 76, 77, 78},
	)
	err = batcher.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}
}
