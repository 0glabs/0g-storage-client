package cmd

import (
	"context"
	"fmt"
	"math"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	kvWriteArgs struct {
		streamId string
		keys     []string
		values   []string
		version  uint64

		url      string
		contract string
		key      string

		node    []string
		indexer string

		expectedReplica uint

		skipTx           bool
		finalityRequired bool
		taskSize         uint
	}

	kvWriteCmd = &cobra.Command{
		Use:   "kv-write",
		Short: "write to kv streams",
		Run:   kvWrite,
	}
)

func init() {
	kvWriteCmd.Flags().StringVar(&kvWriteArgs.streamId, "stream-id", "0x", "stream to read/write")
	kvWriteCmd.MarkFlagRequired("stream-id")

	kvWriteCmd.Flags().StringSliceVar(&kvWriteArgs.keys, "stream-keys", []string{}, "kv keys")
	kvWriteCmd.MarkFlagRequired("kv-keys")
	kvWriteCmd.Flags().StringSliceVar(&kvWriteArgs.values, "stream-values", []string{}, "kv values")
	kvWriteCmd.MarkFlagRequired("kv-values")

	kvWriteCmd.Flags().Uint64Var(&kvWriteArgs.version, "version", math.MaxUint64, "key version")

	kvWriteCmd.Flags().StringVar(&kvWriteArgs.url, "url", "", "Fullnode URL to interact with ZeroGStorage smart contract")
	kvWriteCmd.MarkFlagRequired("url")
	kvWriteCmd.Flags().StringVar(&kvWriteArgs.contract, "contract", "", "ZeroGStorage smart contract to interact with")
	kvWriteCmd.MarkFlagRequired("contract")
	kvWriteCmd.Flags().StringVar(&kvWriteArgs.key, "key", "", "Private key to interact with smart contract")
	kvWriteCmd.MarkFlagRequired("key")

	kvWriteCmd.Flags().StringSliceVar(&kvWriteArgs.node, "node", []string{}, "ZeroGStorage storage node URL")
	kvWriteCmd.Flags().StringVar(&kvWriteArgs.indexer, "indexer", "", "ZeroGStorage indexer URL")

	kvWriteCmd.Flags().UintVar(&kvWriteArgs.expectedReplica, "expected-replica", 1, "expected number of replications to kvWrite")

	// note: for KV operations, skip-tx should by default to be false
	kvWriteCmd.Flags().BoolVar(&kvWriteArgs.skipTx, "skip-tx", false, "Skip sending the transaction on chain if already exists")
	kvWriteCmd.Flags().BoolVar(&kvWriteArgs.finalityRequired, "finality-required", false, "Wait for file finality on nodes to kvWrite")
	kvWriteCmd.Flags().UintVar(&kvWriteArgs.taskSize, "task-size", 10, "Number of segments to kvWrite in single rpc request")

	rootCmd.AddCommand(kvWriteCmd)
}

func kvWrite(*cobra.Command, []string) {
	w3client := blockchain.MustNewWeb3(kvWriteArgs.url, kvWriteArgs.key)
	defer w3client.Close()
	contractAddr := common.HexToAddress(kvWriteArgs.contract)
	flow, err := contract.NewFlowContract(contractAddr, w3client)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create flow contract")
	}

	opt := transfer.UploadOption{
		FinalityRequired: kvWriteArgs.finalityRequired,
		TaskSize:         kvWriteArgs.taskSize,
		ExpectedReplica:  kvWriteArgs.expectedReplica,
		SkipTx:           kvWriteArgs.skipTx,
	}

	var clients []*node.ZgsClient
	if kvWriteArgs.indexer != "" {
		indexerClient, err := indexer.NewClient(kvWriteArgs.indexer, indexer.IndexerClientOption{LogOption: zg_common.LogOption{Logger: logrus.StandardLogger()}})
		if err != nil {
			logrus.WithError(err).Fatal("Failed to initialize indexer client")
		}
		if clients, err = indexerClient.SelectNodes(context.Background(), max(1, opt.ExpectedReplica)); err != nil {
			logrus.WithError(err).Fatal("failed to select nodes from indexer")
		}
	}
	if len(clients) == 0 {
		if len(kvWriteArgs.node) == 0 {
			logrus.Fatal("At least one of --node and --indexer should not be empty")
		}
		clients = node.MustNewZgsClients(kvWriteArgs.node)
		for _, client := range clients {
			defer client.Close()
		}
	}

	batcher := kv.NewBatcher(kvWriteArgs.version, clients, flow, zg_common.LogOption{Logger: logrus.StandardLogger()})
	if len(kvWriteArgs.keys) != len(kvWriteArgs.values) {
		logrus.Fatal("keys and values length mismatch")
	}
	if len(kvWriteArgs.keys) == 0 {
		logrus.Fatal("no keys to write")
	}
	streamId := common.HexToHash(kvWriteArgs.streamId)

	for i := range kvWriteArgs.keys {
		batcher.Set(streamId,
			[]byte(kvWriteArgs.keys[i]),
			[]byte(kvWriteArgs.values[i]),
		)
	}

	err = batcher.Exec(context.Background(), opt)
	if err != nil {
		fmt.Println(err)
		return
	}
}
