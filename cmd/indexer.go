package cmd

import (
	"time"

	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	indexerArgs struct {
		endpoint  string
		nodes     indexer.NodeManagerConfig
		locations indexer.IPLocationConfig
	}

	indexerCmd = &cobra.Command{
		Use:   "indexer",
		Short: "Start indexer service",
		Run:   startIndexer,
	}
)

func init() {
	indexerCmd.Flags().StringVar(&indexerArgs.endpoint, "endpoint", ":12345", "Indexer RPC endpoint")

	indexerCmd.Flags().StringSliceVar(&indexerArgs.nodes.TrustedNodes, "trusted", nil, "Trusted storage node URLs that separated by comma")
	indexerCmd.Flags().StringVar(&indexerArgs.nodes.DiscoveryNode, "node", "", "Storage node to discover peers in P2P network")
	indexerCmd.Flags().DurationVar(&indexerArgs.nodes.DiscoveryInterval, "discover-interval", 10*time.Minute, "Interval to discover peers in network")
	indexerCmd.Flags().DurationVar(&indexerArgs.nodes.UpdateInterval, "update-interval", 10*time.Minute, "Interval to update shard config of discovered peers")

	indexerCmd.Flags().StringVar(&indexerArgs.locations.CacheFile, "ip-location-cache-file", ".ip-location-cache.json", "File name to cache IP locations")
	indexerCmd.Flags().DurationVar(&indexerArgs.locations.CacheWriteInterval, "ip-location-cache-interval", 10*time.Minute, "Interval to write ip locations to cache file")
	indexerCmd.Flags().StringVar(&indexerArgs.locations.AccessToken, "ip-location-token", "", "Access token to retrieve IP location from ipinfo.io")

	indexerCmd.MarkFlagsOneRequired("trusted", "node")

	rootCmd.AddCommand(indexerCmd)
}

func startIndexer(*cobra.Command, []string) {
	indexer.InitDefaultIPLocationManager(indexerArgs.locations)

	closable, err := indexer.InitDefaultNodeManager(indexerArgs.nodes)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize the default node manager")
	}
	defer closable()

	api := indexer.NewIndexerApi()

	logrus.WithFields(logrus.Fields{
		"trusted":  len(indexerArgs.nodes.TrustedNodes),
		"discover": len(indexerArgs.nodes.DiscoveryNode) > 0,
	}).Info("Starting indexer service ...")

	util.MustServeRPC(indexerArgs.endpoint, map[string]interface{}{
		api.Namespace: api,
	})
}
