package cmd

import (
	"context"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	uploadArgs struct {
		file string
		tags string

		url      string
		contract string
		key      string

		node []string

		force    bool
		taskSize uint
	}

	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload file to ZeroGStorage network",
		Run:   upload,
	}
)

func init() {
	uploadCmd.Flags().StringVar(&uploadArgs.file, "file", "", "File name to upload")
	err := uploadCmd.MarkFlagRequired("file")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: file")
	}
	uploadCmd.Flags().StringVar(&uploadArgs.tags, "tags", "0x", "Tags of the file")

	uploadCmd.Flags().StringVar(&uploadArgs.url, "url", "", "Fullnode URL to interact with ZeroGStorage smart contract")
	err = uploadCmd.MarkFlagRequired("url")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: url")
	}
	uploadCmd.Flags().StringVar(&uploadArgs.contract, "contract", "", "ZeroGStorage smart contract to interact with")
	err = uploadCmd.MarkFlagRequired("contract")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: contract")
	}
	uploadCmd.Flags().StringVar(&uploadArgs.key, "key", "", "Private key to interact with smart contract")
	err = uploadCmd.MarkFlagRequired("key")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: key")
	}

	uploadCmd.Flags().StringSliceVar(&uploadArgs.node, "node", []string{}, "ZeroGStorage storage node URL")
	err = uploadCmd.MarkFlagRequired("node")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: node")
	}

	uploadCmd.Flags().BoolVar(&uploadArgs.force, "force", false, "Force to upload file even already exists")
	uploadCmd.Flags().UintVar(&uploadArgs.taskSize, "task-size", 10, "Number of segments to upload in single rpc request")

	rootCmd.AddCommand(uploadCmd)
}

func upload(*cobra.Command, []string) {
	w3client := blockchain.MustNewWeb3(uploadArgs.url, uploadArgs.key)
	defer w3client.Close()
	contractAddr := common.HexToAddress(uploadArgs.contract)
	flow, err := contract.NewFlowContract(contractAddr, w3client)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create flow contract")
	}
	clients := node.MustNewClients(uploadArgs.node)
	for _, client := range clients {
		defer client.Close()
	}

	uploader, err := transfer.NewUploader(flow, clients, zg_common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize uploader")
	}
	opt := transfer.UploadOption{
		Tags:     hexutil.MustDecode(uploadArgs.tags),
		Force:    uploadArgs.force,
		TaskSize: uploadArgs.taskSize,
	}

	file, err := core.Open(uploadArgs.file)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open file")
	}
	defer file.Close()

	if err := uploader.Upload(context.Background(), file, opt); err != nil {
		logrus.WithError(err).Fatal("Failed to upload file")
	}
}
