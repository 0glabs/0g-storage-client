package cmd

import (
	"context"
	"math/big"
	"os"
	"path/filepath"
	"time"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/0glabs/0g-storage-client/transfer/dir"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/web3go"
	"github.com/pkg/errors"
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

		expectedReplica     uint
		maxBatchUploadFiles uint

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
	uploadCmd.Flags().UintVar(&uploadArgs.maxBatchUploadFiles, "max-batch-upload-files", 16, "maximum number of files to upload per batch")

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

	info, err := os.Stat(uploadArgs.file)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get file info")
	}

	uploader, closer, err := newUploader(ctx, w3client, opt)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize uploader")
	}
	defer closer()

	if info.IsDir() {
		err = uploadDir(ctx, uploader, opt)
	} else {
		err = uploadFile(ctx, uploader, opt)
	}

	if err != nil {
		logrus.WithError(err).Fatal("Failed to upload file")
	}
}

func uploadDir(ctx context.Context, uploader *transfer.Uploader, opt transfer.UploadOption) error {
	root, err := dir.BuildFileTree(uploadArgs.file)
	if err != nil {
		return errors.WithMessage(err, "failed to build file tree")
	}

	tdata, err := root.MarshalBinary()
	if err != nil {
		return errors.WithMessage(err, "failed to encode file tree")
	}

	iterdata, err := core.NewDataInMemory(tdata)
	if err != nil {
		return errors.WithMessage(err, "failed to create data in memory")
	}

	merkleTree, err := core.MerkleTree(iterdata)
	if err != nil {
		return errors.WithMessage(err, "failed to create merkle tree")
	}

	var uploadFilePaths []string
	if err = root.Traverse(func(fn *dir.FsNode, path string) error {
		if fn.Type == dir.FileTypeFile {
			uploadFilePaths = append(uploadFilePaths, path)
		}

		if uint(len(uploadFilePaths)) < uploadArgs.maxBatchUploadFiles {
			return nil
		}

		if err := batchUploadFiles(ctx, uploader, uploadFilePaths, opt); err != nil {
			return err
		}

		uploadFilePaths = uploadFilePaths[:0]
		return nil
	}); err == nil {
		err = batchUploadFiles(ctx, uploader, uploadFilePaths, opt, iterdata)
	}

	if err != nil {
		return errors.WithMessage(err, "failed to batch upload files")
	}

	logrus.WithField("merkleRoot", merkleTree.Root()).Info("Directory upload finished")
	return nil
}

func batchUploadFiles(ctx context.Context, up *transfer.Uploader, filePaths []string, opt transfer.UploadOption, datas ...core.IterableData) error {
	for _, path := range filePaths {
		file, err := core.Open(filepath.Join(uploadArgs.file, path))
		if err != nil {
			return errors.WithMessagef(err, "failed to open file %s", path)
		}
		defer file.Close()

		datas = append(datas, file)
	}

	var opts []transfer.UploadOption
	for range datas {
		opts = append(opts, opt)
	}

	txHash, dataRoots, err := up.BatchUpload(ctx, datas, true, opts)
	if err != nil {
		return errors.WithMessage(err, "failed to batch upload file nodes")
	}

	logrus.WithFields(logrus.Fields{
		"txHash":    txHash,
		"dataRoots": dataRoots,
		"filePaths": filePaths,
	}).Info("Batch upload finished")
	return nil
}

func uploadFile(ctx context.Context, up *transfer.Uploader, opt transfer.UploadOption) error {
	file, err := core.Open(uploadArgs.file)
	if err != nil {
		return errors.WithMessage(err, "failed to open file")
	}
	defer file.Close()

	if err := up.Upload(ctx, file, opt); err != nil {
		return errors.WithMessage(err, "failed to upload file")
	}

	return nil
}

func newUploader(ctx context.Context, w3client *web3go.Client, opt transfer.UploadOption) (*transfer.Uploader, func(), error) {
	if uploadArgs.indexer != "" {
		indexerClient, err := indexer.NewClient(uploadArgs.indexer, indexer.IndexerClientOption{
			ProviderOption: providerOption,
			LogOption:      zg_common.LogOption{Logger: logrus.StandardLogger()},
		})
		if err != nil {
			return nil, nil, errors.WithMessage(err, "failed to initialize indexer client")
		}

		up, err := indexerClient.NewUploaderFromIndexerNodes(ctx, w3client, opt.ExpectedReplica, nil)
		if err != nil {
			return nil, nil, errors.WithMessage(err, "failed to initialize uploader")
		}

		return up, func() { indexerClient.Close() }, nil
	}

	clients := node.MustNewZgsClients(uploadArgs.node, providerOption)
	closer := func() {
		for _, client := range clients {
			client.Close()
		}
	}

	up, err := transfer.NewUploader(ctx, w3client, clients, zg_common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		closer()
		return nil, nil, err
	}

	return up, closer, nil
}
