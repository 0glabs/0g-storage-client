package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/node"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ZgsClientOpt = providers.Option{
	RequestTimeout: 3 * time.Second,
}

// NodeManager manages trusted storage nodes and auto discover peers from network.
type NodeManager struct {
	trusted    sync.Map // url -> *node.ZgsClient
	discovered sync.Map // url -> *shard.ShardedNode
}

// Trusted returns trusted sharded nodes.
func (nm *NodeManager) Trusted() ([]*shard.ShardedNode, error) {
	var clients []*node.ZgsClient

	nm.trusted.Range(func(key, value any) bool {
		clients = append(clients, value.(*node.ZgsClient))
		return true
	})

	var nodes []*shard.ShardedNode

	for _, v := range clients {
		start := time.Now()
		config, err := v.GetShardConfig(context.Background())
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to retrieve shard config from trusted storage node %v", v.URL())
		}

		if !config.IsValid() {
			return nil, errors.Errorf("Invalid shard config retrieved from trusted storage node %v", v.URL())
		}

		nodes = append(nodes, &shard.ShardedNode{
			URL:     v.URL(),
			Config:  config,
			Latency: time.Since(start).Milliseconds(),
		})
	}

	return nodes, nil
}

// Discovered returns discovered sharded nodes.
func (nm *NodeManager) Discovered() []*shard.ShardedNode {
	var nodes []*shard.ShardedNode

	nm.discovered.Range(func(key, value any) bool {
		nodes = append(nodes, value.(*shard.ShardedNode))
		return true
	})

	return nodes
}

// AddTrustedNodes add trusted storage nodes.
func (nm *NodeManager) AddTrustedNodes(nodes ...string) error {
	for _, v := range nodes {
		client, err := node.NewZgsClient(v)
		if err != nil {
			return errors.WithMessagef(err, "Failed to create zgs client, url = %v", v)
		}

		nm.trusted.LoadOrStore(v, client)
	}

	return nil
}

func (nm *NodeManager) Close() {
	nm.trusted.Range(func(key, value any) bool {
		value.(*node.ZgsClient).Close()
		return true
	})
}

// Discover discover peers from storage node periodically.
func (nm *NodeManager) Discover(adminClient *node.AdminClient, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// discover once during startup
	if err := nm.discoverOnce(adminClient); err != nil {
		logrus.WithError(err).Warn("Failed to discover storage nodes during startup")
	}

	for range ticker.C {
		if err := nm.discoverOnce(adminClient); err != nil {
			logrus.WithError(err).Warn("Failed to discover storage nodes")
		}
	}
}

func (nm *NodeManager) discoverOnce(adminClient *node.AdminClient) error {
	start := time.Now()
	peers, err := adminClient.GetPeers(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to retrieve peers from storage node")
	}

	logrus.WithFields(logrus.Fields{
		"peers":   len(peers),
		"elapsed": time.Since(start),
	}).Debug("Succeeded to retrieve peers from storage node")

	var numNew int

	for _, v := range peers {
		// public ip address required
		if len(v.SeenIps) == 0 {
			continue
		}

		url := fmt.Sprintf("http://%v:5678", v.SeenIps[0])

		// ignore trusted node
		if _, ok := nm.trusted.Load(url); ok {
			continue
		}

		// discovered already
		if _, ok := nm.discovered.Load(url); ok {
			continue
		}

		// update shard config
		if err = nm.updateNode(url); err != nil {
			logrus.WithError(err).WithField("url", url).Debug("Failed to update shard config")
		}

		logrus.WithField("url", url).Debug("New peer discovered")

		numNew++
	}

	if numNew > 0 {
		logrus.WithField("count", numNew).Info("New peers discovered")
	}

	return nil
}

func (nm *NodeManager) updateNode(url string) error {
	zgsClient, err := node.NewZgsClient(url, ZgsClientOpt)
	if err != nil {
		return errors.WithMessage(err, "Failed to create zgs client")
	}

	start := time.Now()

	config, err := zgsClient.GetShardConfig(context.Background())
	if err != nil {
		return errors.WithMessage(err, "Failed to retrieve shard config from storage node")
	}

	if !config.IsValid() {
		return errors.Errorf("Invalid shard config retrieved %v", config)
	}

	nm.discovered.Store(url, &shard.ShardedNode{
		URL:     url,
		Config:  config,
		Latency: time.Since(start).Milliseconds(),
	})

	return nil
}

// Update update shard config of discovered peers.
func (nm *NodeManager) Update(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		nm.updateOnce()
	}
}

func (nm *NodeManager) updateOnce() {
	var urls []string
	nm.discovered.Range(func(key, value any) bool {
		urls = append(urls, key.(string))
		return true
	})

	if len(urls) == 0 {
		return
	}

	logrus.WithField("nodes", len(urls)).Info("Begin to update shard config")

	start := time.Now()

	for _, v := range urls {
		if err := nm.updateNode(v); err != nil {
			logrus.WithError(err).WithField("url", v).Debug("Failed to update shard config")
			nm.discovered.Delete(v)
		}
	}

	logrus.WithFields(logrus.Fields{
		"nodes":   len(urls),
		"elapsed": time.Since(start),
	}).Info("Completed to update shard config")
}
