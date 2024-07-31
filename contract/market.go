// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MarketMetaData contains all meta data concerning the Market contract.
var MarketMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"beforeLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"uploadSectors\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paddingSectors\",\"type\":\"uint256\"}],\"name\":\"chargeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"flow\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pricePerSector\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pricePerSector_\",\"type\":\"uint256\"}],\"name\":\"setPricePerSector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// MarketABI is the input ABI used to generate the binding from.
// Deprecated: Use MarketMetaData.ABI instead.
var MarketABI = MarketMetaData.ABI

// Market is an auto generated Go binding around an Ethereum contract.
type Market struct {
	MarketCaller     // Read-only binding to the contract
	MarketTransactor // Write-only binding to the contract
	MarketFilterer   // Log filterer for contract events
}

// MarketCaller is an auto generated read-only Go binding around an Ethereum contract.
type MarketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MarketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MarketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MarketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MarketSession struct {
	Contract     *Market           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MarketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MarketCallerSession struct {
	Contract *MarketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MarketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MarketTransactorSession struct {
	Contract     *MarketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MarketRaw is an auto generated low-level Go binding around an Ethereum contract.
type MarketRaw struct {
	Contract *Market // Generic contract binding to access the raw methods on
}

// MarketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MarketCallerRaw struct {
	Contract *MarketCaller // Generic read-only contract binding to access the raw methods on
}

// MarketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MarketTransactorRaw struct {
	Contract *MarketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMarket creates a new instance of Market, bound to a specific deployed contract.
func NewMarket(address common.Address, backend bind.ContractBackend) (*Market, error) {
	contract, err := bindMarket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Market{MarketCaller: MarketCaller{contract: contract}, MarketTransactor: MarketTransactor{contract: contract}, MarketFilterer: MarketFilterer{contract: contract}}, nil
}

// NewMarketCaller creates a new read-only instance of Market, bound to a specific deployed contract.
func NewMarketCaller(address common.Address, caller bind.ContractCaller) (*MarketCaller, error) {
	contract, err := bindMarket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MarketCaller{contract: contract}, nil
}

// NewMarketTransactor creates a new write-only instance of Market, bound to a specific deployed contract.
func NewMarketTransactor(address common.Address, transactor bind.ContractTransactor) (*MarketTransactor, error) {
	contract, err := bindMarket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MarketTransactor{contract: contract}, nil
}

// NewMarketFilterer creates a new log filterer instance of Market, bound to a specific deployed contract.
func NewMarketFilterer(address common.Address, filterer bind.ContractFilterer) (*MarketFilterer, error) {
	contract, err := bindMarket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MarketFilterer{contract: contract}, nil
}

// bindMarket binds a generic wrapper to an already deployed contract.
func bindMarket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MarketABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Market *MarketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Market.Contract.MarketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Market *MarketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Market.Contract.MarketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Market *MarketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Market.Contract.MarketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Market *MarketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Market.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Market *MarketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Market.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Market *MarketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Market.Contract.contract.Transact(opts, method, params...)
}

// Flow is a free data retrieval call binding the contract method 0x343aad82.
//
// Solidity: function flow() view returns(address)
func (_Market *MarketCaller) Flow(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Market.contract.Call(opts, &out, "flow")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Flow is a free data retrieval call binding the contract method 0x343aad82.
//
// Solidity: function flow() view returns(address)
func (_Market *MarketSession) Flow() (common.Address, error) {
	return _Market.Contract.Flow(&_Market.CallOpts)
}

// Flow is a free data retrieval call binding the contract method 0x343aad82.
//
// Solidity: function flow() view returns(address)
func (_Market *MarketCallerSession) Flow() (common.Address, error) {
	return _Market.Contract.Flow(&_Market.CallOpts)
}

// PricePerSector is a free data retrieval call binding the contract method 0x61ec5082.
//
// Solidity: function pricePerSector() view returns(uint256)
func (_Market *MarketCaller) PricePerSector(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Market.contract.Call(opts, &out, "pricePerSector")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PricePerSector is a free data retrieval call binding the contract method 0x61ec5082.
//
// Solidity: function pricePerSector() view returns(uint256)
func (_Market *MarketSession) PricePerSector() (*big.Int, error) {
	return _Market.Contract.PricePerSector(&_Market.CallOpts)
}

// PricePerSector is a free data retrieval call binding the contract method 0x61ec5082.
//
// Solidity: function pricePerSector() view returns(uint256)
func (_Market *MarketCallerSession) PricePerSector() (*big.Int, error) {
	return _Market.Contract.PricePerSector(&_Market.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(address)
func (_Market *MarketCaller) Reward(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Market.contract.Call(opts, &out, "reward")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(address)
func (_Market *MarketSession) Reward() (common.Address, error) {
	return _Market.Contract.Reward(&_Market.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(address)
func (_Market *MarketCallerSession) Reward() (common.Address, error) {
	return _Market.Contract.Reward(&_Market.CallOpts)
}

// ChargeFee is a paid mutator transaction binding the contract method 0xda6eb36a.
//
// Solidity: function chargeFee(uint256 beforeLength, uint256 uploadSectors, uint256 paddingSectors) returns()
func (_Market *MarketTransactor) ChargeFee(opts *bind.TransactOpts, beforeLength *big.Int, uploadSectors *big.Int, paddingSectors *big.Int) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "chargeFee", beforeLength, uploadSectors, paddingSectors)
}

// ChargeFee is a paid mutator transaction binding the contract method 0xda6eb36a.
//
// Solidity: function chargeFee(uint256 beforeLength, uint256 uploadSectors, uint256 paddingSectors) returns()
func (_Market *MarketSession) ChargeFee(beforeLength *big.Int, uploadSectors *big.Int, paddingSectors *big.Int) (*types.Transaction, error) {
	return _Market.Contract.ChargeFee(&_Market.TransactOpts, beforeLength, uploadSectors, paddingSectors)
}

// ChargeFee is a paid mutator transaction binding the contract method 0xda6eb36a.
//
// Solidity: function chargeFee(uint256 beforeLength, uint256 uploadSectors, uint256 paddingSectors) returns()
func (_Market *MarketTransactorSession) ChargeFee(beforeLength *big.Int, uploadSectors *big.Int, paddingSectors *big.Int) (*types.Transaction, error) {
	return _Market.Contract.ChargeFee(&_Market.TransactOpts, beforeLength, uploadSectors, paddingSectors)
}

// SetPricePerSector is a paid mutator transaction binding the contract method 0x14aa90a1.
//
// Solidity: function setPricePerSector(uint256 pricePerSector_) returns()
func (_Market *MarketTransactor) SetPricePerSector(opts *bind.TransactOpts, pricePerSector_ *big.Int) (*types.Transaction, error) {
	return _Market.contract.Transact(opts, "setPricePerSector", pricePerSector_)
}

// SetPricePerSector is a paid mutator transaction binding the contract method 0x14aa90a1.
//
// Solidity: function setPricePerSector(uint256 pricePerSector_) returns()
func (_Market *MarketSession) SetPricePerSector(pricePerSector_ *big.Int) (*types.Transaction, error) {
	return _Market.Contract.SetPricePerSector(&_Market.TransactOpts, pricePerSector_)
}

// SetPricePerSector is a paid mutator transaction binding the contract method 0x14aa90a1.
//
// Solidity: function setPricePerSector(uint256 pricePerSector_) returns()
func (_Market *MarketTransactorSession) SetPricePerSector(pricePerSector_ *big.Int) (*types.Transaction, error) {
	return _Market.Contract.SetPricePerSector(&_Market.TransactOpts, pricePerSector_)
}
