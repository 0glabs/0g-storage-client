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
	opt := transfer.UploadOption{
		Tags:             hexutil.MustDecode(uploadArgs.tags),
		FinalityRequired: uploadArgs.finalityRequired,
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

// uploadDir handles the directory upload process by traversing the directory, batching the file uploads, and then uploading the directory metadata.
func uploadDir(ctx context.Context, uploader *transfer.Uploader, opt transfer.UploadOption) error {
	// Ensure nonce and fee options are not set for directory uploads since we might use batch.
	if opt.Nonce != nil || opt.Fee != nil {
		return errors.New("nonce and fee options are not supported for directory uploading")
	}

	// Build the file tree representation of the directory.
	root, err := dir.BuildFileTree(uploadArgs.file)
	if err != nil {
		return errors.WithMessage(err, "failed to build file tree")
	}

	tdata, err := root.MarshalBinary()
	if err != nil {
		return errors.WithMessage(err, "failed to encode file tree")
	}

	// Create an in-memory data object from the encoded file tree.
	iterdata, err := core.NewDataInMemory(tdata)
	if err != nil {
		return errors.WithMessage(err, "failed to create data in memory")
	}

	// Generate the Merkle tree from the in-memory data.
	merkleTree, err := core.MerkleTree(iterdata)
	if err != nil {
		return errors.WithMessage(err, "failed to create merkle tree")
	}

	// Check if the directory already exists in the storage network.
	if found, err := uploader.FileExists(ctx, merkleTree.Root()); err == nil && found {
		logrus.Info("This folder already exists in the storage network")
		return nil
	}

	// Traverse the file tree and batch upload files.
	var uploadFilePaths []string
	if err = root.Traverse(func(fn *dir.FsNode, path string) error {
		if fn.Type != dir.FileTypeFile {
			return nil
		}

		uploadFilePaths = append(uploadFilePaths, path)
		// If the batch size limit is reached, upload the batch.
		if len(uploadFilePaths) >= int(uploadArgs.maxBatchUploadFiles) {
			if err := batchUploadFiles(ctx, uploader, uploadFilePaths, opt); err != nil {
				return err
			}
			// Clear the slice after uploading the batch.
			uploadFilePaths = uploadFilePaths[:0]
		}

		return nil
	}); err != nil {
		return err
	}

	// Upload any remaining files in the last batch.
	if len(uploadFilePaths) > 0 {
		if err = batchUploadFiles(ctx, uploader, uploadFilePaths, opt); err != nil {
			return err
		}
	}

	// Finally, upload the directory metadata.
	if err := uploader.Upload(ctx, iterdata, opt); err != nil {
		return err
	}

	logrus.WithField("root", merkleTree.Root()).Info("Directory upload finished")
	return nil
}

// batchUploadFiles handles the batch upload of files to the storage network.
func batchUploadFiles(ctx context.Context, up *transfer.Uploader, filePaths []string, opt transfer.UploadOption) error {
	var datas []core.IterableData
	for _, path := range filePaths {
		file, err := core.Open(filepath.Join(uploadArgs.file, path))
		if err != nil {
			return errors.WithMessagef(err, "failed to open file %s", path)
		}
		defer file.Close()

		datas = append(datas, file)
	}

	batchOpt := transfer.BatchUploadOption{}
	for range datas {
		batchOpt.DataOptions = append(batchOpt.DataOptions, opt)
	}

	txHash, dataRoots, err := up.BatchUpload(ctx, datas, true, batchOpt)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"txHash":    txHash,
		"dataRoots": dataRoots,
		"filePaths": filePaths,
	}).Info("Batch upload completed")
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
			return nil, nil, err
		}

		return up, indexerClient.Close, nil
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
