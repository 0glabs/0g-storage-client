package indexer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/node"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	defaultZgsClientOpt = providers.Option{
		RequestTimeout: 3 * time.Second,
	}

	defaultNodeManager = NodeManager{}
)

type NodeManagerConfig struct {
	TrustedNodes []string

	DiscoveryNode     string
	DiscoveryInterval time.Duration

	UpdateInterval time.Duration
}

// NodeManager manages trusted storage nodes and auto discover peers from network.
type NodeManager struct {
	config NodeManagerConfig

	trusted sync.Map // url -> *node.ZgsClient

	discoverNode *node.AdminClient
	discovered   sync.Map // url -> *shard.ShardedNode
}

// InitDefaultNodeManager initializes the default `NodeManager`.
func InitDefaultNodeManager(config NodeManagerConfig) (closable func(), err error) {
	defaultNodeManager.config = config

	if len(config.DiscoveryNode) > 0 {
		if defaultNodeManager.discoverNode, err = node.NewAdminClient(config.DiscoveryNode, defaultZgsClientOpt); err != nil {
			return nil, errors.WithMessage(err, "Failed to create admin client to discover peers")
		}
	}

	if err = defaultNodeManager.AddTrustedNodes(config.TrustedNodes...); err != nil {
		return nil, errors.WithMessage(err, "Failed to add trusted nodes")
	}

	if len(config.DiscoveryNode) > 0 {
		go util.ScheduleNow(defaultNodeManager.discover, config.DiscoveryInterval, "Failed to discover storage nodes")
		go util.Schedule(defaultNodeManager.update, config.UpdateInterval, "Failed to update shard configs")
	}

	return defaultNodeManager.close, nil
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

func parseIP(url string) string {
	url = strings.ToLower(url)

	index := strings.Index(url, ":")
	if index == -1 {
		return url
	}

	url = url[index+3:]

	index = strings.Index(url, ":")
	if index == -1 {
		return url
	}

	return url[:index]
}

// AddTrustedNodes add trusted storage nodes.
func (nm *NodeManager) AddTrustedNodes(nodes ...string) error {
	for _, v := range nodes {
		ip := parseIP(v)
		if _, err := defaultIPLocationManager.Query(ip); err != nil {
			logrus.WithError(err).WithField("ip", ip).Warn("Failed to query IP location")
		}

		client, err := node.NewZgsClient(v)
		if err != nil {
			return errors.WithMessagef(err, "Failed to create zgs client, url = %v", v)
		}

		nm.trusted.LoadOrStore(v, client)
	}

	return nil
}

func (nm *NodeManager) close() {
	nm.trusted.Range(func(key, value any) bool {
		value.(*node.ZgsClient).Close()
		return true
	})

	if nm.discoverNode != nil {
		nm.discoverNode.Close()
	}
}

// discover discovers peers from storage node.
func (nm *NodeManager) discover() error {
	start := time.Now()
	peers, err := nm.discoverNode.GetPeers(context.Background())
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
		node, err := nm.updateNode(url)
		if err != nil {
			logrus.WithError(err).WithField("url", url).Debug("Failed to update shard config")
		} else {
			logrus.WithFields(logrus.Fields{
				"url":     url,
				"shard":   node.Config,
				"latency": node.Latency,
			}).Debug("New peer discovered")
		}

		numNew++
	}

	if numNew > 0 {
		logrus.WithField("count", numNew).Info("New peers discovered")
	}

	return nil
}

// updateNode updates the shard config of specified storage node by `url`.
func (nm *NodeManager) updateNode(url string) (*shard.ShardedNode, error) {
	// query ip location at first
	ip := parseIP(url)
	if _, err := defaultIPLocationManager.Query(ip); err != nil {
		logrus.WithError(err).WithField("ip", ip).Warn("Failed to query IP location")
	}

	zgsClient, err := node.NewZgsClient(url, defaultZgsClientOpt)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create zgs client")
	}
	defer zgsClient.Close()

	start := time.Now()

	config, err := zgsClient.GetShardConfig(context.Background())
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to retrieve shard config from storage node")
	}

	if !config.IsValid() {
		return nil, errors.Errorf("Invalid shard config retrieved %v", config)
	}

	node := &shard.ShardedNode{
		URL:     url,
		Config:  config,
		Latency: time.Since(start).Milliseconds(),
		Since:   time.Now().Unix(),
	}

	nm.discovered.Store(url, node)

	return node, nil
}

// update updates shard configs of all storage nodes.
func (nm *NodeManager) update() error {
	var urls []string
	nm.discovered.Range(func(key, value any) bool {
		urls = append(urls, key.(string))
		return true
	})

	if len(urls) == 0 {
		return nil
	}

	logrus.WithField("nodes", len(urls)).Info("Begin to update shard config")

	start := time.Now()

	for _, v := range urls {
		if _, err := nm.updateNode(v); err != nil {
			logrus.WithError(err).WithField("url", v).Debug("Failed to update shard config")
			nm.discovered.Delete(v)
		}
	}

	logrus.WithFields(logrus.Fields{
		"nodes":   len(urls),
		"elapsed": time.Since(start),
	}).Info("Completed to update shard config")

	return nil
}
