package cmd

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
	"github.com/zero-gravity-labs/zerog-storage-client/transfer"
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
		disperse bool
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
	uploadCmd.MarkFlagRequired("file")
	uploadCmd.Flags().StringVar(&uploadArgs.tags, "tags", "0x", "Tags of the file")

	uploadCmd.Flags().StringVar(&uploadArgs.url, "url", "", "Fullnode URL to interact with ZeroGStorage smart contract")
	uploadCmd.MarkFlagRequired("url")
	uploadCmd.Flags().StringVar(&uploadArgs.contract, "contract", "", "ZeroGStorage smart contract to interact with")
	uploadCmd.MarkFlagRequired("contract")
	uploadCmd.Flags().StringVar(&uploadArgs.key, "key", "", "Private key to interact with smart contract")
	uploadCmd.MarkFlagRequired("key")

	uploadCmd.Flags().StringSliceVar(&uploadArgs.node, "node", []string{}, "ZeroGStorage storage node URL")
	uploadCmd.MarkFlagRequired("node")

	uploadCmd.Flags().BoolVar(&uploadArgs.force, "force", false, "Force to upload file even already exists")
	uploadCmd.Flags().BoolVar(&uploadArgs.disperse, "disperse", false, "Disperse file amoung nodes")
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

	uploader := transfer.NewUploader(flow, clients)
	opt := transfer.UploadOption{
		Tags:     hexutil.MustDecode(uploadArgs.tags),
		Force:    uploadArgs.force,
		Disperse: uploadArgs.disperse,
		TaskSize: uploadArgs.taskSize,
	}

	file, err := core.Open(uploadArgs.file)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open file")
	}
	defer file.Close()

	if err := uploader.Upload(file, opt); err != nil {
		logrus.WithError(err).Fatal("Failed to upload file")
	}
}
