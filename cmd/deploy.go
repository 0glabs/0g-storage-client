package cmd

import (
	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	deployArgs struct {
		url            string
		key            string
		bytecodeOrFile string
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy ZeroGStorage contract to specified blockchain",
		Run:   deploy,
	}
)

func init() {
	deployCmd.Flags().StringVar(&deployArgs.url, "url", "", "Fullnode URL to interact with blockchain")
	deployCmd.MarkFlagRequired("url")
	deployCmd.Flags().StringVar(&deployArgs.key, "key", "", "Private key to create smart contract")
	deployCmd.MarkFlagRequired("key")
	deployCmd.Flags().StringVar(&deployArgs.bytecodeOrFile, "bytecode", "", "ZeroGStorage smart contract bytecode")
	deployCmd.MarkFlagRequired("bytecode")

	rootCmd.AddCommand(deployCmd)
}

func deploy(*cobra.Command, []string) {
	client := blockchain.MustNewWeb3(deployArgs.url, deployArgs.key)

	contract, err := blockchain.Deploy(client, deployArgs.bytecodeOrFile)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to deploy smart contract")
	}

	logrus.WithField("contract", contract).Info("Smart contract deployed")
}
