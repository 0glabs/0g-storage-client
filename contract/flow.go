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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"startMerkleRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"submissionIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"flowLength\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"context\",\"type\":\"bytes32\"}],\"name\":\"NewEpoch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"identity\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"submissionIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startPos\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tags\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"internalType\":\"structSubmissionNode[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structSubmission\",\"name\":\"submission\",\"type\":\"tuple\"}],\"name\":\"Submit\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getContext\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"mineStart\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"flowRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"flowLength\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockDigest\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"internalType\":\"structMineContext\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"}],\"name\":\"getEpochRange\",\"outputs\":[{\"components\":[{\"internalType\":\"uint128\",\"name\":\"start\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"end\",\"type\":\"uint128\"}],\"internalType\":\"structEpochRange\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"makeContext\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSubmissions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"tags\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"internalType\":\"structSubmissionNode[]\",\"name\":\"nodes\",\"type\":\"tuple[]\"}],\"internalType\":\"structSubmission\",\"name\":\"submission\",\"type\":\"tuple\"}],\"name\":\"submit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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
