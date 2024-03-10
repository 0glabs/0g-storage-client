package cmd

import (
	"fmt"
	"os"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logLevel       string
	logColorForced bool

	rootCmd = &cobra.Command{
		Use:   "0g-storage-client",
		Short: "ZeroGStorage client to interact with ZeroGStorage network",
		PersistentPreRun: func(*cobra.Command, []string) {
			initLog()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logrus.InfoLevel.String(), "Log level")
	rootCmd.PersistentFlags().BoolVar(&logColorForced, "log-force-color", false, "Force to output colorful logs")
	rootCmd.PersistentFlags().Uint64Var(&blockchain.CustomGasPrice, "gas-price", 0, "Custom gas price to send transaction")
	rootCmd.PersistentFlags().Uint64Var(&blockchain.CustomGasLimit, "gas-limit", 0, "Custom gas limit to send transaction")
	rootCmd.PersistentFlags().BoolVar(&blockchain.Web3LogEnabled, "web3-log-enabled", false, "Enable log for web3 RPC")
}

func initLog() {
	if logColorForced {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).WithField("level", logLevel).Fatal("Failed to parse log level")
	}

	logrus.SetLevel(level)
}

// Execute is the command line entrypoint.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
