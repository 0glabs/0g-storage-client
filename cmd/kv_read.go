package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/0glabs/0g-storage-client/kv"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	kvReadArgs struct {
		streamId string
		keys     []string
		version  uint64

		node string
	}

	kvReadCmd = &cobra.Command{
		Use:   "kv-read",
		Short: "read kv streams",
		Run:   kvRead,
	}
)

func init() {
	kvReadCmd.Flags().StringVar(&kvReadArgs.streamId, "stream-id", "0x", "stream to read/write")
	kvReadCmd.MarkFlagRequired("stream-id")

	kvReadCmd.Flags().StringSliceVar(&kvReadArgs.keys, "stream-keys", []string{}, "kv keys")
	kvReadCmd.MarkFlagRequired("kv-keys")

	kvReadCmd.Flags().Uint64Var(&kvReadArgs.version, "version", math.MaxUint64, "key version")

	kvReadCmd.Flags().StringVar(&kvReadArgs.node, "node", "", "kv node url")
	kvReadCmd.MarkFlagRequired("node")

	rootCmd.AddCommand(kvReadCmd)
}

func kvRead(*cobra.Command, []string) {
	ctx := context.Background()
	client, err := node.NewClient(kvReadArgs.node)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize kv client")
	}
	kvClient := kv.NewClient(client, nil)
	streamId := common.HexToHash(kvReadArgs.streamId)

	m := make(map[string]string)
	for _, key := range kvReadArgs.keys {
		val, err := kvClient.GetValue(ctx, streamId, []byte(key))
		if err != nil {
			logrus.WithError(err).Fatalf("failed to read key %v", key)
		}
		m[key] = string(val.Data)
	}
	bs, _ := json.Marshal(m)
	fmt.Println(string(bs))
}
