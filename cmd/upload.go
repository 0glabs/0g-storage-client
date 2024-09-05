package cmd

import (
	"context"
	"math/big"
	"time"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	uploadArgs struct {
		file string
		tags string

		url string
		key string

		node    []string
		indexer string

		expectedReplica uint

		skipTx           bool
		finalityRequired bool
		taskSize         uint

		fee   float64
		nonce uint

		timeout time.Duration
	}

	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload file to ZeroGStorage network",
		Run:   upload,
	}
)

func init() {
	uploadCmd.Flags().StringVar(&uploadArgs.file, "file", "", "File name to upload")
	uploadCmd.MarkFlagRequired("file")
	uploadCmd.Flags().StringVar(&uploadArgs.tags, "tags", "0x", "Tags of the file")

	uploadCmd.Flags().StringVar(&uploadArgs.url, "url", "", "Fullnode URL to interact with ZeroGStorage smart contract")
	uploadCmd.MarkFlagRequired("url")
	uploadCmd.Flags().StringVar(&uploadArgs.key, "key", "", "Private key to interact with smart contract")
	uploadCmd.MarkFlagRequired("key")

	uploadCmd.Flags().StringSliceVar(&uploadArgs.node, "node", []string{}, "ZeroGStorage storage node URL")
	uploadCmd.Flags().StringVar(&uploadArgs.indexer, "indexer", "", "ZeroGStorage indexer URL")
	uploadCmd.MarkFlagsOneRequired("indexer", "node")
	uploadCmd.MarkFlagsMutuallyExclusive("indexer", "node")

	uploadCmd.Flags().UintVar(&uploadArgs.expectedReplica, "expected-replica", 1, "expected number of replications to upload")

	uploadCmd.Flags().BoolVar(&uploadArgs.skipTx, "skip-tx", true, "Skip sending the transaction on chain if already exists")
	uploadCmd.Flags().BoolVar(&uploadArgs.finalityRequired, "finality-required", false, "Wait for file finality on nodes to upload")
	uploadCmd.Flags().UintVar(&uploadArgs.taskSize, "task-size", 10, "Number of segments to upload in single rpc request")

	uploadCmd.Flags().DurationVar(&uploadArgs.timeout, "timeout", 0, "cli task timeout, 0 for no timeout")

	uploadCmd.Flags().Float64Var(&uploadArgs.fee, "fee", 0, "fee paid in a0gi")
	uploadCmd.Flags().UintVar(&uploadArgs.nonce, "nonce", 0, "nonce of upload transaction")

	rootCmd.AddCommand(uploadCmd)
}

func upload(*cobra.Command, []string) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if uploadArgs.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, uploadArgs.timeout)
		defer cancel()
	}

	w3client := blockchain.MustNewWeb3(uploadArgs.url, uploadArgs.key, providerOption)
	defer w3client.Close()

	var fee *big.Int
	if uploadArgs.fee > 0 {
		feeInA0GI := big.NewFloat(uploadArgs.fee)
		fee, _ = feeInA0GI.Mul(feeInA0GI, big.NewFloat(1e18)).Int(nil)
	}
	var nonce *big.Int
	if uploadArgs.nonce > 0 {
		nonce = big.NewInt(int64(uploadArgs.nonce))
	}
	finalityRequired := transfer.TransactionPacked
	if uploadArgs.finalityRequired {
		finalityRequired = transfer.FileFinalized
	}
	opt := transfer.UploadOption{
		Tags:             hexutil.MustDecode(uploadArgs.tags),
		FinalityRequired: finalityRequired,
		TaskSize:         uploadArgs.taskSize,
		ExpectedReplica:  uploadArgs.expectedReplica,
		SkipTx:           uploadArgs.skipTx,
		Fee:              fee,
		Nonce:            nonce,
	}

	file, err := core.Open(uploadArgs.file)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open file")
	}
	defer file.Close()

	if uploadArgs.indexer != "" {
		indexerClient, err := indexer.NewClient(uploadArgs.indexer, indexer.IndexerClientOption{
			ProviderOption: providerOption,
			LogOption:      zg_common.LogOption{Logger: logrus.StandardLogger()},
		})
		if err != nil {
			logrus.WithError(err).Fatal("Failed to initialize indexer client")
		}
		if _, err := indexerClient.Upload(ctx, w3client, file, opt); err != nil {
			logrus.WithError(err).Fatal("Failed to upload file")
		}
		return
	}

	clients := node.MustNewZgsClients(uploadArgs.node, providerOption)
	for _, client := range clients {
		defer client.Close()
	}

	uploader, err := transfer.NewUploader(ctx, w3client, clients, zg_common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize uploader")
	}

	if _, err := uploader.Upload(ctx, file, opt); err != nil {
		logrus.WithError(err).Fatal("Failed to upload file")
	}
}
