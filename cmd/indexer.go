package cmd

import (
	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	nodes    []string
	endpoint string

	indexderCmd = &cobra.Command{
		Use:   "indexer",
		Short: "Start indexer service",
		Run:   startIndexer,
	}
)

func init() {
	indexderCmd.Flags().StringSliceVar(&nodes, "nodes", nil, "Storage node URLs that separated by comma")
	err := indexderCmd.MarkFlagRequired("nodes")
	if err != nil {
		logrus.WithError(err).Error("Failed to MarkFlagRequired flag: nodes")
	}
	indexderCmd.Flags().StringVar(&endpoint, "endpoint", ":12345", "Indexer RPC endpoint")

	rootCmd.AddCommand(indexderCmd)
}

func startIndexer(*cobra.Command, []string) {
	var clients []*node.Client

	for _, v := range nodes {
		client, err := node.NewClient(v)
		if err != nil {
			logrus.WithError(err).WithField("node", v).Fatal("Failed to dail storage node")
		}

		clients = append(clients, client)
	}

	defer func() {
		for _, v := range clients {
			v.Close()
		}
	}()

	api := indexer.NewIndexerApi(clients)

	util.MustServeRPC(endpoint, map[string]interface{}{
		api.Namespace: api,
	})
}
