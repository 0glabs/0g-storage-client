package cmd

import (
	"context"
	"time"

	"github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	downloadArgs struct {
		file string

		indexer string
		nodes   []string

		root  string
		proof bool

		timeout time.Duration
	}

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download file from ZeroGStorage network",
		Run:   download,
	}
)

func init() {
	downloadCmd.Flags().StringVar(&downloadArgs.file, "file", "", "File name to download")
	downloadCmd.MarkFlagRequired("file")

	downloadCmd.Flags().StringSliceVar(&downloadArgs.nodes, "node", []string{}, "ZeroGStorage storage node URL. Multiple nodes could be specified and separated by comma, e.g. url1,url2,url3")
	downloadCmd.Flags().StringVar(&downloadArgs.indexer, "indexer", "", "ZeroGStorage indexer URL")
	downloadCmd.MarkFlagsOneRequired("indexer", "node")

	downloadCmd.Flags().StringVar(&downloadArgs.root, "root", "", "Merkle root to download file")
	downloadCmd.MarkFlagRequired("root")
	downloadCmd.Flags().BoolVar(&downloadArgs.proof, "proof", false, "Whether to download with merkle proof for validation")

	downloadCmd.Flags().DurationVar(&downloadArgs.timeout, "timeout", 0, "cli task timeout, 0 for no timeout")

	rootCmd.AddCommand(downloadCmd)
}

func download(*cobra.Command, []string) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if downloadArgs.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, downloadArgs.timeout)
		defer cancel()
	}

	if downloadArgs.indexer != "" {
		indexerClient, err := indexer.NewClient(downloadArgs.indexer, indexer.IndexerClientOption{LogOption: common.LogOption{Logger: logrus.StandardLogger()}})
		if err != nil {
			logrus.WithError(err).Fatal("Failed to initialize indexer client")
		}
		if err := indexerClient.Download(ctx, downloadArgs.root, downloadArgs.file, downloadArgs.proof); err != nil {
			logrus.WithError(err).Fatal("Failed to download file from indexer")
		}
		return
	}

	nodes := node.MustNewZgsClients(downloadArgs.nodes)

	downloader, err := transfer.NewDownloader(nodes, common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize downloader")
	}

	if err := downloader.Download(ctx, downloadArgs.root, downloadArgs.file, downloadArgs.proof); err != nil {
		logrus.WithError(err).Fatal("Failed to download file")
	}
}
