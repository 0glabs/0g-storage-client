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

// EpochRange is an auto generated low-level Go binding around an user-defined struct.
type EpochRange struct {
	Start *big.Int
	End   *big.Int
}

// EpochRangeWithContextDigest is an auto generated low-level Go binding around an user-defined struct.
type EpochRangeWithContextDigest struct {
	Start  *big.Int
	End    *big.Int
	Digest [32]byte
}

// MineContext is an auto generated low-level Go binding around an user-defined struct.
type MineContext struct {
	Epoch       *big.Int
	MineStart   *big.Int
	FlowRoot    [32]byte
	FlowLength  *big.Int
	BlockDigest [32]byte
	Digest      [32]byte
}

// Submission is an auto generated low-level Go binding around an user-defined struct.
type Submission struct {
	Length *big.Int
	Tags   []byte
	Nodes  []SubmissionNode
}

// SubmissionNode is an auto generated low-level Go binding around an user-defined struct.
type SubmissionNode struct {
	Root   [32]byte
	Height *big.Int
}

// FlowMetaData contains all meta data concerning the Flow contract.
var FlowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InvalidSubmission\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"startMerkleRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"submissionIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"flowLength\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"context\",\"type\":\"bytes32\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"identity\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"submissionIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startPos\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tags\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"internalType\":\"structSubmissionNode[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structSubmission\",\"name\":\"submission\",\"type\":\"tuple\"}],\"name\":\"Submit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tags\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"internalType\":\"structSubmissionNode[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"internalType\":\"structSubmission[]\",\"name\":\"submissions\",\"type\":\"tuple[]\"}],\"name\":\"batchSubmit\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"indexes\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[]\",\"name\":\"digests\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"startIndexes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"lengths\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"epochStartPosition\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"firstBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContext\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mineStart\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"flowRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"flowLength\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structMineContext\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"name\":\"getEpochRange\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"start\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"end\",\"type\":\"uint128\"}],\"internalType\":\"structEpochRange\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"makeContext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"makeContextWithResult\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mineStart\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"flowRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"flowLength\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structMineContext\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_length\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"alignExp\",\"type\":\"uint256\"}],\"name\":\"nextAlign\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_length\",\"type\":\"uint256\"}],\"name\":\"nextPow2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSubmissions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"targetPosition\",\"type\":\"uint128\"}],\"name\":\"queryContextAtPosition\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"start\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"end\",\"type\":\"uint128\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structEpochRangeWithContextDigest\",\"name\":\"range\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootHistory\",\"outputs\":[{\"internalType\":\"contractIDigestHistory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submissionIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tags\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"internalType\":\"structSubmissionNode[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"internalType\":\"structSubmission\",\"name\":\"submission\",\"type\":\"tuple\"}],\"name\":\"submit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstagedHeight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"zeros\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// FlowABI is the input ABI used to generate the binding from.
// Deprecated: Use FlowMetaData.ABI instead.
var FlowABI = FlowMetaData.ABI

// Flow is an auto generated Go binding around an Ethereum contract.
type Flow struct {
	FlowCaller     // Read-only binding to the contract
	FlowTransactor // Write-only binding to the contract
	FlowFilterer   // Log filterer for contract events
}

// FlowCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlowSession struct {
	Contract     *Flow             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlowCallerSession struct {
	Contract *FlowCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// FlowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlowTransactorSession struct {
	Contract     *FlowTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlowRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlowRaw struct {
	Contract *Flow // Generic contract binding to access the raw methods on
}

// FlowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlowCallerRaw struct {
	Contract *FlowCaller // Generic read-only contract binding to access the raw methods on
}

// FlowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlowTransactorRaw struct {
	Contract *FlowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlow creates a new instance of Flow, bound to a specific deployed contract.
func NewFlow(address common.Address, backend bind.ContractBackend) (*Flow, error) {
	contract, err := bindFlow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Flow{FlowCaller: FlowCaller{contract: contract}, FlowTransactor: FlowTransactor{contract: contract}, FlowFilterer: FlowFilterer{contract: contract}}, nil
}

// NewFlowCaller creates a new read-only instance of Flow, bound to a specific deployed contract.
func NewFlowCaller(address common.Address, caller bind.ContractCaller) (*FlowCaller, error) {
	contract, err := bindFlow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlowCaller{contract: contract}, nil
}

// NewFlowTransactor creates a new write-only instance of Flow, bound to a specific deployed contract.
func NewFlowTransactor(address common.Address, transactor bind.ContractTransactor) (*FlowTransactor, error) {
	contract, err := bindFlow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlowTransactor{contract: contract}, nil
}

// NewFlowFilterer creates a new log filterer instance of Flow, bound to a specific deployed contract.
func NewFlowFilterer(address common.Address, filterer bind.ContractFilterer) (*FlowFilterer, error) {
	contract, err := bindFlow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlowFilterer{contract: contract}, nil
}

// bindFlow binds a generic wrapper to an already deployed contract.
func bindFlow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FlowABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flow *FlowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flow.Contract.FlowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flow *FlowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flow.Contract.FlowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flow *FlowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flow.Contract.FlowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Flow *FlowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Flow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Flow *FlowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Flow *FlowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Flow.Contract.contract.Transact(opts, method, params...)
}

// CurrentLength is a free data retrieval call binding the contract method 0xa3d35f36.
//
// Solidity: function currentLength() view returns(uint256)
func (_Flow *FlowCaller) CurrentLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "currentLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentLength is a free data retrieval call binding the contract method 0xa3d35f36.
//
// Solidity: function currentLength() view returns(uint256)
func (_Flow *FlowSession) CurrentLength() (*big.Int, error) {
	return _Flow.Contract.CurrentLength(&_Flow.CallOpts)
}

// CurrentLength is a free data retrieval call binding the contract method 0xa3d35f36.
//
// Solidity: function currentLength() view returns(uint256)
func (_Flow *FlowCallerSession) CurrentLength() (*big.Int, error) {
	return _Flow.Contract.CurrentLength(&_Flow.CallOpts)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_Flow *FlowCaller) Epoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "epoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_Flow *FlowSession) Epoch() (*big.Int, error) {
	return _Flow.Contract.Epoch(&_Flow.CallOpts)
}

// Epoch is a free data retrieval call binding the contract method 0x900cf0cf.
//
// Solidity: function epoch() view returns(uint256)
func (_Flow *FlowCallerSession) Epoch() (*big.Int, error) {
	return _Flow.Contract.Epoch(&_Flow.CallOpts)
}

// EpochStartPosition is a free data retrieval call binding the contract method 0x93e405a0.
//
// Solidity: function epochStartPosition() view returns(uint256)
func (_Flow *FlowCaller) EpochStartPosition(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "epochStartPosition")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochStartPosition is a free data retrieval call binding the contract method 0x93e405a0.
//
// Solidity: function epochStartPosition() view returns(uint256)
func (_Flow *FlowSession) EpochStartPosition() (*big.Int, error) {
	return _Flow.Contract.EpochStartPosition(&_Flow.CallOpts)
}

// EpochStartPosition is a free data retrieval call binding the contract method 0x93e405a0.
//
// Solidity: function epochStartPosition() view returns(uint256)
func (_Flow *FlowCallerSession) EpochStartPosition() (*big.Int, error) {
	return _Flow.Contract.EpochStartPosition(&_Flow.CallOpts)
}

// FirstBlock is a free data retrieval call binding the contract method 0x231b0268.
//
// Solidity: function firstBlock() view returns(uint256)
func (_Flow *FlowCaller) FirstBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "firstBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FirstBlock is a free data retrieval call binding the contract method 0x231b0268.
//
// Solidity: function firstBlock() view returns(uint256)
func (_Flow *FlowSession) FirstBlock() (*big.Int, error) {
	return _Flow.Contract.FirstBlock(&_Flow.CallOpts)
}

// FirstBlock is a free data retrieval call binding the contract method 0x231b0268.
//
// Solidity: function firstBlock() view returns(uint256)
func (_Flow *FlowCallerSession) FirstBlock() (*big.Int, error) {
	return _Flow.Contract.FirstBlock(&_Flow.CallOpts)
}

// GetContext is a free data retrieval call binding the contract method 0x127f0f07.
//
// Solidity: function getContext() view returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowCaller) GetContext(opts *bind.CallOpts) (MineContext, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "getContext")

	if err != nil {
		return *new(MineContext), err
	}

	out0 := *abi.ConvertType(out[0], new(MineContext)).(*MineContext)

	return out0, err

}

// GetContext is a free data retrieval call binding the contract method 0x127f0f07.
//
// Solidity: function getContext() view returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowSession) GetContext() (MineContext, error) {
	return _Flow.Contract.GetContext(&_Flow.CallOpts)
}

// GetContext is a free data retrieval call binding the contract method 0x127f0f07.
//
// Solidity: function getContext() view returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowCallerSession) GetContext() (MineContext, error) {
	return _Flow.Contract.GetContext(&_Flow.CallOpts)
}

// GetEpochRange is a free data retrieval call binding the contract method 0x7d590708.
//
// Solidity: function getEpochRange(bytes32 digest) view returns((uint128,uint128))
func (_Flow *FlowCaller) GetEpochRange(opts *bind.CallOpts, digest [32]byte) (EpochRange, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "getEpochRange", digest)

	if err != nil {
		return *new(EpochRange), err
	}

	out0 := *abi.ConvertType(out[0], new(EpochRange)).(*EpochRange)

	return out0, err

}

// GetEpochRange is a free data retrieval call binding the contract method 0x7d590708.
//
// Solidity: function getEpochRange(bytes32 digest) view returns((uint128,uint128))
func (_Flow *FlowSession) GetEpochRange(digest [32]byte) (EpochRange, error) {
	return _Flow.Contract.GetEpochRange(&_Flow.CallOpts, digest)
}

// GetEpochRange is a free data retrieval call binding the contract method 0x7d590708.
//
// Solidity: function getEpochRange(bytes32 digest) view returns((uint128,uint128))
func (_Flow *FlowCallerSession) GetEpochRange(digest [32]byte) (EpochRange, error) {
	return _Flow.Contract.GetEpochRange(&_Flow.CallOpts, digest)
}

// NextAlign is a free data retrieval call binding the contract method 0x555430a1.
//
// Solidity: function nextAlign(uint256 _length, uint256 alignExp) pure returns(uint256)
func (_Flow *FlowCaller) NextAlign(opts *bind.CallOpts, _length *big.Int, alignExp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "nextAlign", _length, alignExp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextAlign is a free data retrieval call binding the contract method 0x555430a1.
//
// Solidity: function nextAlign(uint256 _length, uint256 alignExp) pure returns(uint256)
func (_Flow *FlowSession) NextAlign(_length *big.Int, alignExp *big.Int) (*big.Int, error) {
	return _Flow.Contract.NextAlign(&_Flow.CallOpts, _length, alignExp)
}

// NextAlign is a free data retrieval call binding the contract method 0x555430a1.
//
// Solidity: function nextAlign(uint256 _length, uint256 alignExp) pure returns(uint256)
func (_Flow *FlowCallerSession) NextAlign(_length *big.Int, alignExp *big.Int) (*big.Int, error) {
	return _Flow.Contract.NextAlign(&_Flow.CallOpts, _length, alignExp)
}

// NextPow2 is a free data retrieval call binding the contract method 0x3d75d9c2.
//
// Solidity: function nextPow2(uint256 _length) pure returns(uint256)
func (_Flow *FlowCaller) NextPow2(opts *bind.CallOpts, _length *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "nextPow2", _length)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextPow2 is a free data retrieval call binding the contract method 0x3d75d9c2.
//
// Solidity: function nextPow2(uint256 _length) pure returns(uint256)
func (_Flow *FlowSession) NextPow2(_length *big.Int) (*big.Int, error) {
	return _Flow.Contract.NextPow2(&_Flow.CallOpts, _length)
}

// NextPow2 is a free data retrieval call binding the contract method 0x3d75d9c2.
//
// Solidity: function nextPow2(uint256 _length) pure returns(uint256)
func (_Flow *FlowCallerSession) NextPow2(_length *big.Int) (*big.Int, error) {
	return _Flow.Contract.NextPow2(&_Flow.CallOpts, _length)
}

// NumSubmissions is a free data retrieval call binding the contract method 0x77e19824.
//
// Solidity: function numSubmissions() view returns(uint256)
func (_Flow *FlowCaller) NumSubmissions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "numSubmissions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumSubmissions is a free data retrieval call binding the contract method 0x77e19824.
//
// Solidity: function numSubmissions() view returns(uint256)
func (_Flow *FlowSession) NumSubmissions() (*big.Int, error) {
	return _Flow.Contract.NumSubmissions(&_Flow.CallOpts)
}

// NumSubmissions is a free data retrieval call binding the contract method 0x77e19824.
//
// Solidity: function numSubmissions() view returns(uint256)
func (_Flow *FlowCallerSession) NumSubmissions() (*big.Int, error) {
	return _Flow.Contract.NumSubmissions(&_Flow.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Flow *FlowCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Flow *FlowSession) Paused() (bool, error) {
	return _Flow.Contract.Paused(&_Flow.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Flow *FlowCallerSession) Paused() (bool, error) {
	return _Flow.Contract.Paused(&_Flow.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_Flow *FlowCaller) Root(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_Flow *FlowSession) Root() ([32]byte, error) {
	return _Flow.Contract.Root(&_Flow.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_Flow *FlowCallerSession) Root() ([32]byte, error) {
	return _Flow.Contract.Root(&_Flow.CallOpts)
}

// RootHistory is a free data retrieval call binding the contract method 0xc7dd5221.
//
// Solidity: function rootHistory() view returns(address)
func (_Flow *FlowCaller) RootHistory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "rootHistory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RootHistory is a free data retrieval call binding the contract method 0xc7dd5221.
//
// Solidity: function rootHistory() view returns(address)
func (_Flow *FlowSession) RootHistory() (common.Address, error) {
	return _Flow.Contract.RootHistory(&_Flow.CallOpts)
}

// RootHistory is a free data retrieval call binding the contract method 0xc7dd5221.
//
// Solidity: function rootHistory() view returns(address)
func (_Flow *FlowCallerSession) RootHistory() (common.Address, error) {
	return _Flow.Contract.RootHistory(&_Flow.CallOpts)
}

// SubmissionIndex is a free data retrieval call binding the contract method 0xb8a409ac.
//
// Solidity: function submissionIndex() view returns(uint256)
func (_Flow *FlowCaller) SubmissionIndex(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "submissionIndex")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SubmissionIndex is a free data retrieval call binding the contract method 0xb8a409ac.
//
// Solidity: function submissionIndex() view returns(uint256)
func (_Flow *FlowSession) SubmissionIndex() (*big.Int, error) {
	return _Flow.Contract.SubmissionIndex(&_Flow.CallOpts)
}

// SubmissionIndex is a free data retrieval call binding the contract method 0xb8a409ac.
//
// Solidity: function submissionIndex() view returns(uint256)
func (_Flow *FlowCallerSession) SubmissionIndex() (*big.Int, error) {
	return _Flow.Contract.SubmissionIndex(&_Flow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Flow *FlowCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Flow *FlowSession) Token() (common.Address, error) {
	return _Flow.Contract.Token(&_Flow.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Flow *FlowCallerSession) Token() (common.Address, error) {
	return _Flow.Contract.Token(&_Flow.CallOpts)
}

// UnstagedHeight is a free data retrieval call binding the contract method 0x364800ec.
//
// Solidity: function unstagedHeight() view returns(uint256)
func (_Flow *FlowCaller) UnstagedHeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "unstagedHeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstagedHeight is a free data retrieval call binding the contract method 0x364800ec.
//
// Solidity: function unstagedHeight() view returns(uint256)
func (_Flow *FlowSession) UnstagedHeight() (*big.Int, error) {
	return _Flow.Contract.UnstagedHeight(&_Flow.CallOpts)
}

// UnstagedHeight is a free data retrieval call binding the contract method 0x364800ec.
//
// Solidity: function unstagedHeight() view returns(uint256)
func (_Flow *FlowCallerSession) UnstagedHeight() (*big.Int, error) {
	return _Flow.Contract.UnstagedHeight(&_Flow.CallOpts)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 height) pure returns(bytes32)
func (_Flow *FlowCaller) Zeros(opts *bind.CallOpts, height *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Flow.contract.Call(opts, &out, "zeros", height)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 height) pure returns(bytes32)
func (_Flow *FlowSession) Zeros(height *big.Int) ([32]byte, error) {
	return _Flow.Contract.Zeros(&_Flow.CallOpts, height)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 height) pure returns(bytes32)
func (_Flow *FlowCallerSession) Zeros(height *big.Int) ([32]byte, error) {
	return _Flow.Contract.Zeros(&_Flow.CallOpts, height)
}

// BatchSubmit is a paid mutator transaction binding the contract method 0x9e62a38e.
//
// Solidity: function batchSubmit((uint256,bytes,(bytes32,uint256)[])[] submissions) returns(uint256[] indexes, bytes32[] digests, uint256[] startIndexes, uint256[] lengths)
func (_Flow *FlowTransactor) BatchSubmit(opts *bind.TransactOpts, submissions []Submission) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "batchSubmit", submissions)
}

// BatchSubmit is a paid mutator transaction binding the contract method 0x9e62a38e.
//
// Solidity: function batchSubmit((uint256,bytes,(bytes32,uint256)[])[] submissions) returns(uint256[] indexes, bytes32[] digests, uint256[] startIndexes, uint256[] lengths)
func (_Flow *FlowSession) BatchSubmit(submissions []Submission) (*types.Transaction, error) {
	return _Flow.Contract.BatchSubmit(&_Flow.TransactOpts, submissions)
}

// BatchSubmit is a paid mutator transaction binding the contract method 0x9e62a38e.
//
// Solidity: function batchSubmit((uint256,bytes,(bytes32,uint256)[])[] submissions) returns(uint256[] indexes, bytes32[] digests, uint256[] startIndexes, uint256[] lengths)
func (_Flow *FlowTransactorSession) BatchSubmit(submissions []Submission) (*types.Transaction, error) {
	return _Flow.Contract.BatchSubmit(&_Flow.TransactOpts, submissions)
}

// CommitRoot is a paid mutator transaction binding the contract method 0xd34353c9.
//
// Solidity: function commitRoot() returns()
func (_Flow *FlowTransactor) CommitRoot(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "commitRoot")
}

// CommitRoot is a paid mutator transaction binding the contract method 0xd34353c9.
//
// Solidity: function commitRoot() returns()
func (_Flow *FlowSession) CommitRoot() (*types.Transaction, error) {
	return _Flow.Contract.CommitRoot(&_Flow.TransactOpts)
}

// CommitRoot is a paid mutator transaction binding the contract method 0xd34353c9.
//
// Solidity: function commitRoot() returns()
func (_Flow *FlowTransactorSession) CommitRoot() (*types.Transaction, error) {
	return _Flow.Contract.CommitRoot(&_Flow.TransactOpts)
}

// MakeContext is a paid mutator transaction binding the contract method 0x38d45e10.
//
// Solidity: function makeContext() returns()
func (_Flow *FlowTransactor) MakeContext(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "makeContext")
}

// MakeContext is a paid mutator transaction binding the contract method 0x38d45e10.
//
// Solidity: function makeContext() returns()
func (_Flow *FlowSession) MakeContext() (*types.Transaction, error) {
	return _Flow.Contract.MakeContext(&_Flow.TransactOpts)
}

// MakeContext is a paid mutator transaction binding the contract method 0x38d45e10.
//
// Solidity: function makeContext() returns()
func (_Flow *FlowTransactorSession) MakeContext() (*types.Transaction, error) {
	return _Flow.Contract.MakeContext(&_Flow.TransactOpts)
}

// MakeContextWithResult is a paid mutator transaction binding the contract method 0xb464b53e.
//
// Solidity: function makeContextWithResult() returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowTransactor) MakeContextWithResult(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "makeContextWithResult")
}

// MakeContextWithResult is a paid mutator transaction binding the contract method 0xb464b53e.
//
// Solidity: function makeContextWithResult() returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowSession) MakeContextWithResult() (*types.Transaction, error) {
	return _Flow.Contract.MakeContextWithResult(&_Flow.TransactOpts)
}

// MakeContextWithResult is a paid mutator transaction binding the contract method 0xb464b53e.
//
// Solidity: function makeContextWithResult() returns((uint256,uint256,bytes32,uint256,bytes32,bytes32))
func (_Flow *FlowTransactorSession) MakeContextWithResult() (*types.Transaction, error) {
	return _Flow.Contract.MakeContextWithResult(&_Flow.TransactOpts)
}

// QueryContextAtPosition is a paid mutator transaction binding the contract method 0x31bae174.
//
// Solidity: function queryContextAtPosition(uint128 targetPosition) returns((uint128,uint128,bytes32) range)
func (_Flow *FlowTransactor) QueryContextAtPosition(opts *bind.TransactOpts, targetPosition *big.Int) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "queryContextAtPosition", targetPosition)
}

// QueryContextAtPosition is a paid mutator transaction binding the contract method 0x31bae174.
//
// Solidity: function queryContextAtPosition(uint128 targetPosition) returns((uint128,uint128,bytes32) range)
func (_Flow *FlowSession) QueryContextAtPosition(targetPosition *big.Int) (*types.Transaction, error) {
	return _Flow.Contract.QueryContextAtPosition(&_Flow.TransactOpts, targetPosition)
}

// QueryContextAtPosition is a paid mutator transaction binding the contract method 0x31bae174.
//
// Solidity: function queryContextAtPosition(uint128 targetPosition) returns((uint128,uint128,bytes32) range)
func (_Flow *FlowTransactorSession) QueryContextAtPosition(targetPosition *big.Int) (*types.Transaction, error) {
	return _Flow.Contract.QueryContextAtPosition(&_Flow.TransactOpts, targetPosition)
}

// Submit is a paid mutator transaction binding the contract method 0xef3e12dc.
//
// Solidity: function submit((uint256,bytes,(bytes32,uint256)[]) submission) returns(uint256, bytes32, uint256, uint256)
func (_Flow *FlowTransactor) Submit(opts *bind.TransactOpts, submission Submission) (*types.Transaction, error) {
	return _Flow.contract.Transact(opts, "submit", submission)
}

// Submit is a paid mutator transaction binding the contract method 0xef3e12dc.
//
// Solidity: function submit((uint256,bytes,(bytes32,uint256)[]) submission) returns(uint256, bytes32, uint256, uint256)
func (_Flow *FlowSession) Submit(submission Submission) (*types.Transaction, error) {
	return _Flow.Contract.Submit(&_Flow.TransactOpts, submission)
}

// Submit is a paid mutator transaction binding the contract method 0xef3e12dc.
//
// Solidity: function submit((uint256,bytes,(bytes32,uint256)[]) submission) returns(uint256, bytes32, uint256, uint256)
func (_Flow *FlowTransactorSession) Submit(submission Submission) (*types.Transaction, error) {
	return _Flow.Contract.Submit(&_Flow.TransactOpts, submission)
}

// FlowNewEpochIterator is returned from FilterNewEpoch and is used to iterate over the raw logs and unpacked data for NewEpoch events raised by the Flow contract.
type FlowNewEpochIterator struct {
	Event *FlowNewEpoch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowNewEpochIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowNewEpoch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowNewEpoch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowNewEpochIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowNewEpochIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowNewEpoch represents a NewEpoch event raised by the Flow contract.
type FlowNewEpoch struct {
	Sender          common.Address
	Index           *big.Int
	StartMerkleRoot [32]byte
	SubmissionIndex *big.Int
	FlowLength      *big.Int
	Context         [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewEpoch is a free log retrieval operation binding the contract event 0xbc8a3fd82465d43f1709e44ed882f7e1af0147274196ef1ec009f5d52ff4e993.
//
// Solidity: event NewEpoch(address indexed sender, uint256 indexed index, bytes32 startMerkleRoot, uint256 submissionIndex, uint256 flowLength, bytes32 context)
func (_Flow *FlowFilterer) FilterNewEpoch(opts *bind.FilterOpts, sender []common.Address, index []*big.Int) (*FlowNewEpochIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Flow.contract.FilterLogs(opts, "NewEpoch", senderRule, indexRule)
	if err != nil {
		return nil, err
	}
	return &FlowNewEpochIterator{contract: _Flow.contract, event: "NewEpoch", logs: logs, sub: sub}, nil
}

// WatchNewEpoch is a free log subscription operation binding the contract event 0xbc8a3fd82465d43f1709e44ed882f7e1af0147274196ef1ec009f5d52ff4e993.
//
// Solidity: event NewEpoch(address indexed sender, uint256 indexed index, bytes32 startMerkleRoot, uint256 submissionIndex, uint256 flowLength, bytes32 context)
func (_Flow *FlowFilterer) WatchNewEpoch(opts *bind.WatchOpts, sink chan<- *FlowNewEpoch, sender []common.Address, index []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}

	logs, sub, err := _Flow.contract.WatchLogs(opts, "NewEpoch", senderRule, indexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowNewEpoch)
				if err := _Flow.contract.UnpackLog(event, "NewEpoch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewEpoch is a log parse operation binding the contract event 0xbc8a3fd82465d43f1709e44ed882f7e1af0147274196ef1ec009f5d52ff4e993.
//
// Solidity: event NewEpoch(address indexed sender, uint256 indexed index, bytes32 startMerkleRoot, uint256 submissionIndex, uint256 flowLength, bytes32 context)
func (_Flow *FlowFilterer) ParseNewEpoch(log types.Log) (*FlowNewEpoch, error) {
	event := new(FlowNewEpoch)
	if err := _Flow.contract.UnpackLog(event, "NewEpoch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Flow contract.
type FlowPausedIterator struct {
	Event *FlowPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowPaused represents a Paused event raised by the Flow contract.
type FlowPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Flow *FlowFilterer) FilterPaused(opts *bind.FilterOpts) (*FlowPausedIterator, error) {

	logs, sub, err := _Flow.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &FlowPausedIterator{contract: _Flow.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Flow *FlowFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *FlowPaused) (event.Subscription, error) {

	logs, sub, err := _Flow.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowPaused)
				if err := _Flow.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Flow *FlowFilterer) ParsePaused(log types.Log) (*FlowPaused, error) {
	event := new(FlowPaused)
	if err := _Flow.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowSubmitIterator is returned from FilterSubmit and is used to iterate over the raw logs and unpacked data for Submit events raised by the Flow contract.
type FlowSubmitIterator struct {
	Event *FlowSubmit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowSubmitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowSubmit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowSubmit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowSubmitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowSubmitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowSubmit represents a Submit event raised by the Flow contract.
type FlowSubmit struct {
	Sender          common.Address
	Identity        [32]byte
	SubmissionIndex *big.Int
	StartPos        *big.Int
	Length          *big.Int
	Submission      Submission
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterSubmit is a free log retrieval operation binding the contract event 0x167ce04d2aa1981994d3a31695da0d785373335b1078cec239a1a3a2c7675555.
//
// Solidity: event Submit(address indexed sender, bytes32 indexed identity, uint256 submissionIndex, uint256 startPos, uint256 length, (uint256,bytes,(bytes32,uint256)[]) submission)
func (_Flow *FlowFilterer) FilterSubmit(opts *bind.FilterOpts, sender []common.Address, identity [][32]byte) (*FlowSubmitIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _Flow.contract.FilterLogs(opts, "Submit", senderRule, identityRule)
	if err != nil {
		return nil, err
	}
	return &FlowSubmitIterator{contract: _Flow.contract, event: "Submit", logs: logs, sub: sub}, nil
}

// WatchSubmit is a free log subscription operation binding the contract event 0x167ce04d2aa1981994d3a31695da0d785373335b1078cec239a1a3a2c7675555.
//
// Solidity: event Submit(address indexed sender, bytes32 indexed identity, uint256 submissionIndex, uint256 startPos, uint256 length, (uint256,bytes,(bytes32,uint256)[]) submission)
func (_Flow *FlowFilterer) WatchSubmit(opts *bind.WatchOpts, sink chan<- *FlowSubmit, sender []common.Address, identity [][32]byte) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var identityRule []interface{}
	for _, identityItem := range identity {
		identityRule = append(identityRule, identityItem)
	}

	logs, sub, err := _Flow.contract.WatchLogs(opts, "Submit", senderRule, identityRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowSubmit)
				if err := _Flow.contract.UnpackLog(event, "Submit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSubmit is a log parse operation binding the contract event 0x167ce04d2aa1981994d3a31695da0d785373335b1078cec239a1a3a2c7675555.
//
// Solidity: event Submit(address indexed sender, bytes32 indexed identity, uint256 submissionIndex, uint256 startPos, uint256 length, (uint256,bytes,(bytes32,uint256)[]) submission)
func (_Flow *FlowFilterer) ParseSubmit(log types.Log) (*FlowSubmit, error) {
	event := new(FlowSubmit)
	if err := _Flow.contract.UnpackLog(event, "Submit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlowUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Flow contract.
type FlowUnpausedIterator struct {
	Event *FlowUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlowUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlowUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlowUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlowUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlowUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlowUnpaused represents a Unpaused event raised by the Flow contract.
type FlowUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Flow *FlowFilterer) FilterUnpaused(opts *bind.FilterOpts) (*FlowUnpausedIterator, error) {

	logs, sub, err := _Flow.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &FlowUnpausedIterator{contract: _Flow.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Flow *FlowFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *FlowUnpaused) (event.Subscription, error) {

	logs, sub, err := _Flow.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlowUnpaused)
				if err := _Flow.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Flow *FlowFilterer) ParseUnpaused(log types.Log) (*FlowUnpaused, error) {
	event := new(FlowUnpaused)
	if err := _Flow.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
