package cmd

import (
	"context"
	"time"

	"github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type downloadArgument struct {
	file string

	indexer string
	nodes   []string

	root  string
	proof bool

	timeout time.Duration
}

func bindDownloadFlags(cmd *cobra.Command, args *downloadArgument) {
	cmd.Flags().StringVar(&args.file, "file", "", "File name to download")
	cmd.MarkFlagRequired("file")

	cmd.Flags().StringSliceVar(&args.nodes, "node", []string{}, "ZeroGStorage storage node URL. Multiple nodes could be specified and separated by comma, e.g. url1,url2,url3")
	cmd.Flags().StringVar(&args.indexer, "indexer", "", "ZeroGStorage indexer URL")
	cmd.MarkFlagsOneRequired("indexer", "node")

	cmd.Flags().StringVar(&args.root, "root", "", "Merkle root to download file")
	cmd.MarkFlagRequired("root")
	cmd.Flags().BoolVar(&args.proof, "proof", false, "Whether to download with merkle proof for validation")

	cmd.Flags().DurationVar(&args.timeout, "timeout", 0, "cli task timeout, 0 for no timeout")
}

var (
	downloadArgs downloadArgument

	downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "Download file from ZeroGStorage network",
		Run:   download,
	}
)

func init() {
	bindDownloadFlags(downloadCmd, &downloadArgs)

	rootCmd.AddCommand(downloadCmd)
}

func download(*cobra.Command, []string) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if downloadArgs.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, downloadArgs.timeout)
		defer cancel()
	}

	downloader, closer, err := newDownloader(downloadArgs)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize downloader")
	}
	defer closer()

	if err := downloader.Download(ctx, downloadArgs.root, downloadArgs.file, downloadArgs.proof); err != nil {
		logrus.WithError(err).Fatal("Failed to download file")
	}
}

func newDownloader(args downloadArgument) (transfer.IDownloader, func(), error) {
	if args.indexer != "" {
		indexerClient, err := indexer.NewClient(args.indexer, indexer.IndexerClientOption{
			ProviderOption: providerOption,
			LogOption:      common.LogOption{Logger: logrus.StandardLogger()},
		})
		if err != nil {
			return nil, nil, errors.WithMessage(err, "failed to initialize indexer client")
		}

		return indexerClient, indexerClient.Close, nil
	}

	clients := node.MustNewZgsClients(args.nodes, providerOption)
	closer := func() {
		for _, client := range clients {
			client.Close()
		}
	}

	downloader, err := transfer.NewDownloader(clients, common.LogOption{Logger: logrus.StandardLogger()})
	if err != nil {
		closer()
		return nil, nil, err
	}

	return downloader, closer, nil
}
