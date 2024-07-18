package cmd

import (
	"time"

	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	trustedNodes []string
	storageNode  string
	endpoint     string

	discoverInterval time.Duration
	updateInterval   time.Duration

	ipLocationCacheFile       string
	ipLocationPersistInterval time.Duration

	indexerCmd = &cobra.Command{
		Use:   "indexer",
		Short: "Start indexer service",
		Run:   startIndexer,
	}
)

func init() {
	indexerCmd.Flags().StringSliceVar(&trustedNodes, "trusted", nil, "Trusted storage node URLs that separated by comma")
	indexerCmd.Flags().StringVar(&storageNode, "node", "", "Storage node to discover peers in P2P network")
	indexerCmd.Flags().StringVar(&endpoint, "endpoint", ":12345", "Indexer RPC endpoint")
	indexerCmd.Flags().DurationVar(&discoverInterval, "discover-interval", 10*time.Minute, "Interval to discover peers in network")
	indexerCmd.Flags().DurationVar(&updateInterval, "update-interval", 10*time.Minute, "Interval to update shard config of discovered peers")
	indexerCmd.Flags().StringVar(&ipLocationCacheFile, "ip-location-cache-file", ".ip-location-cache.json", "File name to cache ip locations")
	indexerCmd.Flags().DurationVar(&ipLocationPersistInterval, "ip-location-cache-interval", 10*time.Minute, "Interval to write ip locations to cache file")
	indexerCmd.Flags().StringVar(&indexer.IPLocationToken, "ip-location-token", "", "Access token to retrieve IP location from ipinfo.io")

	rootCmd.AddCommand(indexerCmd)
}

func startIndexer(*cobra.Command, []string) {
	if len(trustedNodes) == 0 && len(storageNode) == 0 {
		logrus.Fatal("Neither 'trusted' nor 'node' specified")
	}

	// initialize ip location cache at first
	indexer.StartIPLocationCache(ipLocationCacheFile, ipLocationPersistInterval)

	var manager indexer.NodeManager

	// add trusted nodes
	if err := manager.AddTrustedNodes(trustedNodes...); err != nil {
		logrus.WithError(err).Fatal("Failed to add trusted nodes")
	}
	defer manager.Close()

	// discover and update periodically
	if len(storageNode) > 0 {
		adminClient, err := node.NewAdminClient(storageNode, indexer.ZgsClientOpt)
		if err != nil {
			logrus.WithError(err).WithField("url", storageNode).Fatal("Failed to create admin client")
		}

		go manager.Discover(adminClient, discoverInterval)
		go manager.Update(updateInterval)
	}

	api := indexer.NewIndexerApi(&manager)

	logrus.WithFields(logrus.Fields{
		"trusted":  len(trustedNodes),
		"discover": len(storageNode) > 0,
	}).Info("Starting indexer service ...")

	util.MustServeRPC(endpoint, map[string]interface{}{
		api.Namespace: api,
	})
}
