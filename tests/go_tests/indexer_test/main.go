package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func runTest() error {
	ctx := context.Background()
	// load args, indexer's trusted & discover node is node[0]
	args := os.Args[1:]
	key := args[0]
	contractAddr := eth_common.HexToAddress(args[1])
	chainUrl := args[2]
	zgsNodeUrls := strings.Split(args[3], ",")
	indexerUrl := args[4]

	w3client := blockchain.MustNewWeb3(chainUrl, key)
	defer w3client.Close()
	flow, err := contract.NewFlowContract(contractAddr, w3client)
	if err != nil {
		return fmt.Errorf("failed to create flow contract")
	}
	// upload by indexer
	data, err := core.NewDataInMemory([]byte("indexer_test_data"))
	if err != nil {
		return errors.WithMessage(err, "failed to initialize data")
	}
	indexerClient, err := indexer.NewClient(indexerUrl, indexer.IndexerClientOption{LogOption: common.LogOption{Logger: logrus.StandardLogger()}})
	if err != nil {
		return errors.WithMessage(err, "failed to initialize indexer client")
	}
	if _, err := indexerClient.Upload(ctx, flow, data, transfer.UploadOption{
		FinalityRequired: true,
	}); err != nil {
		return errors.WithMessage(err, "failed to upload file")
	}
	tree, err := core.MerkleTree(data)
	if err != nil {
		return errors.WithMessage(err, "failed to build merkle tree")
	}
	root := tree.Root()
	// check file location
	locations, err := indexerClient.GetFileLocations(ctx, root.Hex())
	if err != nil {
		return errors.WithMessage(err, "failed to get file locations")
	}
	if len(locations) != 1 {
		return fmt.Errorf("unexpected file location length: %v", len(locations))
	}
	if locations[0].URL != zgsNodeUrls[0] {
		return fmt.Errorf("unexpected file location: %v", locations[0].URL)
	}
	// upload data to second node
	data, err = core.NewDataInMemory([]byte("indexer_test_data_2"))
	if err != nil {
		return errors.WithMessage(err, "failed to initialize data")
	}
	clients := node.MustNewZgsClients(zgsNodeUrls[1:])
	for _, client := range clients {
		defer client.Close()
	}

	uploader, err := transfer.NewUploader(ctx, flow, clients, common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		return errors.WithMessage(err, "failed to initialize uploader")
	}

	if _, err := uploader.Upload(context.Background(), data, transfer.UploadOption{
		FinalityRequired: true,
	}); err != nil {
		return errors.WithMessage(err, "failed to upload file")
	}
	tree, err = core.MerkleTree(data)
	if err != nil {
		return errors.WithMessage(err, "failed to build merkle tree")
	}
	root = tree.Root()

	client0 := node.MustNewZgsClient(zgsNodeUrls[0])
	for {
		info, err := client0.GetFileInfo(ctx, root)
		if err != nil {
			return errors.WithMessage(err, "failed to get file info")
		}
		if info != nil {
			break
		}
		time.Sleep(time.Second)
	}
	// node list
	_, err = indexerClient.GetShardedNodes(ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to get sharded nodes")
	}
	return nil
}

func main() {
	if err := util.WaitUntil(runTest, time.Minute*3); err != nil {
		logrus.WithError(err).Fatalf("indexer test failed")
	}
}
