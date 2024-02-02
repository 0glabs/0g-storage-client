package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var CustomGasPrice uint64
var CustomGasLimit uint64

func Deploy(clientWithSigner *web3go.Client, dataOrFile string) (common.Address, error) {
	signer, err := defaultSigner(clientWithSigner)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Failed to detect account")
	}
	from := signer.Address()

	bytecode, err := parseBytecode(dataOrFile)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Failed to parse bytecode")
	}

	var gasPrice *hexutil.Big
	if CustomGasPrice > 0 {
		gasPrice = (*hexutil.Big)(new(big.Int).SetUint64(CustomGasPrice))
	}

	var gasLimit *hexutil.Uint64
	if CustomGasLimit > 0 {
		gasLimit = (*hexutil.Uint64)(&CustomGasLimit)
	}

	txHash, err := clientWithSigner.Eth.SendTransactionByArgs(types.TransactionArgs{
		From:     &from,
		Data:     &bytecode,
		GasPrice: gasPrice,
		Gas:      gasLimit,
	})
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Failed to send transaction")
	}

	logrus.WithField("hash", txHash).Info("Transaction sent to blockchain")

	receipt, err := WaitForReceipt(clientWithSigner, txHash, true)
	if err != nil {
		return common.Address{}, errors.WithMessage(err, "Failed to wait for receipt")
	}

	return *receipt.ContractAddress, nil
}

func parseBytecode(dataOrFile string) (hexutil.Bytes, error) {
	if strings.HasPrefix(dataOrFile, "0x") {
		return hexutil.Decode(dataOrFile)
	}

	content, err := ioutil.ReadFile(dataOrFile)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read file")
	}

	var data map[string]interface{}
	if err = json.Unmarshal(content, &data); err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal JSON")
	}

	bytecode, ok := data["bytecode"]
	if !ok {
		return nil, errors.New("bytecode field not found in JSON file")
	}

	if bytecodeStr, ok := bytecode.(string); ok {
		return hexutil.Decode(bytecodeStr)
	}

	bytecodeObj, ok := bytecode.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid type for bytecode field")
	}

	bytecode, ok = bytecodeObj["object"]
	if !ok {
		return nil, errors.New("bytecode.object field not found in JSON file")
	}

	if bytecodeStr, ok := bytecode.(string); ok {
		return hexutil.Decode(bytecodeStr)
	}

	return nil, errors.New("invalid type for bytecode field")
}

type Contract struct {
	client  *web3go.Client
	account common.Address // account to send transaction
	signer  bind.SignerFn
}

func NewContract(clientWithSigner *web3go.Client, signerFn bind.SignerFn) (*Contract, error) {
	signer, err := defaultSigner(clientWithSigner)
	if err != nil {
		return nil, err
	}

	return &Contract{
		client:  clientWithSigner,
		account: signer.Address(),
		signer:  signerFn,
	}, nil
}

func (c *Contract) CreateTransactOpts() (*bind.TransactOpts, error) {
	var gasPrice *big.Int
	if CustomGasPrice > 0 {
		gasPrice = new(big.Int).SetUint64(CustomGasPrice)
	}

	return &bind.TransactOpts{
		From:     c.account,
		GasPrice: gasPrice,
		GasLimit: CustomGasLimit,
		Signer:   c.signer,
	}, nil
}

func (c *Contract) WaitForReceipt(txHash common.Hash, successRequired bool, opts ...RetryOption) (*types.Receipt, error) {
	return WaitForReceipt(c.client, txHash, successRequired, opts...)
}
