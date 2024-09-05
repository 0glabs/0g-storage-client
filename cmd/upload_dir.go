package cmd

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	uploadDirArgs uploadArgument

	uploadDirCmd = &cobra.Command{
		Use:   "upload-dir",
		Short: "Upload directory to ZeroGStorage network",
		Run:   uploadDir,
	}
)

func init() {
	bindUploadFlags(uploadDirCmd, &uploadDirArgs, false)
	rootCmd.AddCommand(uploadDirCmd)
}

func uploadDir(*cobra.Command, []string) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if uploadDirArgs.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, uploadDirArgs.timeout)
		defer cancel()
	}

	w3client := blockchain.MustNewWeb3(uploadDirArgs.url, uploadDirArgs.key, providerOption)
	defer w3client.Close()

	opt := transfer.UploadOption{
		Tags:             hexutil.MustDecode(uploadDirArgs.tags),
		FinalityRequired: uploadDirArgs.finalityRequired,
		TaskSize:         uploadDirArgs.taskSize,
		ExpectedReplica:  uploadDirArgs.expectedReplica,
		SkipTx:           uploadDirArgs.skipTx,
	}

	uploader, closer, err := newUploader(ctx, uploadDirArgs, w3client, opt)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize uploader")
	}
	defer closer()

	txnHash, rootHash, err := uploader.UploadDir(ctx, uploadDirArgs.file, opt)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to upload directory")
	}

	logrus.WithFields(logrus.Fields{
		"txnHash":  txnHash,
		"rootHash": rootHash,
	}).Info("Directory uploaded done")
}
