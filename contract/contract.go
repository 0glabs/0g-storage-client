package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/types"
	"github.com/sirupsen/logrus"
)

type FlowContract struct {
	*blockchain.Contract
	*Flow
	clientWithSigner *web3go.Client
}

type TxRetryOption struct {
	Timeout          time.Duration
	MaxNonGasRetries int
	MaxGasPrice      *big.Int
}

var SpecifiedBlockError = "Specified block header does not exist"
var DefaultTimeout = 15 * time.Second
var DefaultMaxNonGasRetries = 20

func IsRetriableSubmitLogEntryError(msg string) bool {
	return strings.Contains(msg, SpecifiedBlockError) || strings.Contains(msg, "mempool") || strings.Contains(msg, "timeout")
}

func NewFlowContract(flowAddress common.Address, clientWithSigner *web3go.Client) (*FlowContract, error) {
	backend, signer := clientWithSigner.ToClientForContract()

	contract, err := blockchain.NewContract(clientWithSigner, signer)
	if err != nil {
		return nil, err
	}

	flow, err := NewFlow(flowAddress, backend)
	if err != nil {
		return nil, err
	}

	return &FlowContract{contract, flow, clientWithSigner}, nil
}

func (f *FlowContract) GetGasPrice() (*big.Int, error) {
	gasPrice, err := f.clientWithSigner.Eth.GasPrice()
	if err != nil {
		return nil, err
	}

	return gasPrice, nil
}

func (f *FlowContract) GetMarketContract(ctx context.Context) (*Market, error) {
	marketAddr, err := f.Market(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, err
	}

	backend, _ := f.clientWithSigner.ToClientForContract()

	market, err := NewMarket(marketAddr, backend)
	if err != nil {
		return nil, err
	}

	return market, nil
}

func (submission Submission) String() string {
	var heights []uint64
	for _, v := range submission.Nodes {
		heights = append(heights, v.Height.Uint64())
	}

	return fmt.Sprintf("{ Size: %v, Heights: %v }", submission.Length, heights)
}

func (submission Submission) Root() common.Hash {
	numNodes := len(submission.Nodes)

	// should be never occur
	if numNodes == 0 {
		return common.Hash{}
	}

	// calculate root in reverse order
	root := submission.Nodes[numNodes-1].Root
	for i := 1; i < numNodes; i++ {
		left := submission.Nodes[numNodes-1-i]
		root = crypto.Keccak256Hash(left.Root[:], root[:])
	}

	return root
}

func TransactWithGasAdjustment(
	contract *FlowContract,
	method string,
	opts *bind.TransactOpts,
	retryOpts *TxRetryOption,
	params ...interface{},
) (*types.Receipt, error) {
	// Set timeout and max non-gas retries from retryOpts if provided.
	if retryOpts == nil {
		retryOpts = &TxRetryOption{
			Timeout:          DefaultTimeout,
			MaxNonGasRetries: DefaultMaxNonGasRetries,
		}
	}

	if retryOpts.MaxNonGasRetries == 0 {
		retryOpts.MaxNonGasRetries = DefaultMaxNonGasRetries
	}

	if retryOpts.Timeout == 0 {
		retryOpts.Timeout = DefaultTimeout
	}

	logrus.WithField("timeout", retryOpts.Timeout).WithField("maxNonGasRetries", retryOpts.MaxNonGasRetries).Debug("Set retry options")

	if opts.GasPrice == nil {
		// Get the current gas price if not set.
		gasPrice, err := contract.GetGasPrice()
		if err != nil {
			return nil, fmt.Errorf("failed to get gas price: %w", err)
		}
		opts.GasPrice = gasPrice
		logrus.WithField("gasPrice", opts.GasPrice).Debug("Receive current gas price from chain node")
	}

	logrus.WithField("gasPrice", opts.GasPrice).Info("Set gas price")

	nRetries := 0
	for {
		// Create a fresh context per iteration.
		ctx, cancel := context.WithTimeout(context.Background(), retryOpts.Timeout)
		opts.Context = ctx
		tx, err := contract.FlowTransactor.contract.Transact(opts, method, params...)
		
		var receipt *types.Receipt
		if err == nil {
			// Wait for successful execution
			receipt, err = contract.WaitForReceipt(ctx, tx.Hash(), true, blockchain.RetryOption{NRetries: retryOpts.MaxNonGasRetries})
			if err == nil {
				cancel() // cancel this iteration's context
				return receipt, nil
			}
		}
		cancel() // cancel this iteration's context

		logrus.WithError(err).Error("Failed to send transaction")

		errStr := strings.ToLower(err.Error())

		if !IsRetriableSubmitLogEntryError(errStr) {
			return nil, fmt.Errorf("failed to send transaction: %w", err)
		}

		if strings.Contains(errStr, "mempool") || strings.Contains(errStr, "timeout") {
			if retryOpts.MaxGasPrice == nil {
				return nil, fmt.Errorf("mempool full and no max gas price is set, failed to send transaction: %w", err)
			} else {
				newGasPrice := new(big.Int).Mul(opts.GasPrice, big.NewInt(15))
				newGasPrice.Div(newGasPrice, big.NewInt(10))
				if newGasPrice.Cmp(retryOpts.MaxGasPrice) > 0 {
					opts.GasPrice = new(big.Int).Set(retryOpts.MaxGasPrice)
				} else {
					opts.GasPrice = newGasPrice
				}
				logrus.WithError(err).Infof("Increasing gas price to %v due to mempool/timeout error", opts.GasPrice)
			}
		} else {
			nRetries++
			if nRetries >= retryOpts.MaxNonGasRetries {
				return nil, fmt.Errorf("failed to send transaction after %d retries: %w", nRetries, err)
			}
			logrus.WithError(err).Infof("Retrying with same gas price %v, attempt %d", opts.GasPrice, nRetries)
		}

		time.Sleep(10 * time.Second)
	}

}

func (submission Submission) Fee(pricePerSector *big.Int) *big.Int {
	var sectors int64
	for _, node := range submission.Nodes {
		sectors += 1 << node.Height.Int64()
	}

	return big.NewInt(0).Mul(big.NewInt(sectors), pricePerSector)
}
