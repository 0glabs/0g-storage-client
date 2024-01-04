package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zero-gravity-labs/zerog-storage-client/gateway"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

var (
	gatewayArgs struct {
		nodes []string
	}

	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Start gateway service",
		Run:   startGateway,
	}
)

func init() {
	gatewayCmd.Flags().StringSliceVar(&gatewayArgs.nodes, "nodes", []string{
		"http://127.0.0.1:5678",
		"http://127.0.0.1:5679",
		"http://127.0.0.1:5680",
	}, "Storage node list separated by comma")
	gatewayCmd.Flags().StringVar(&gateway.LocalFileRepo, "repo", "", "Local file repository")

	rootCmd.AddCommand(gatewayCmd)
}

func startGateway(*cobra.Command, []string) {
	nodes := node.MustNewClients(gatewayArgs.nodes)
	gateway.MustServeLocal(nodes)
}
