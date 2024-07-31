package contract

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openweb3/web3go"
)

type FlowContract struct {
	*blockchain.Contract
	*Flow
	clientWithSigner *web3go.Client
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

func (submission Submission) Fee(pricePerSector *big.Int) *big.Int {
	var sectors int64
	for _, node := range submission.Nodes {
		sectors += 1 << node.Height.Int64()
	}

	return big.NewInt(0).Mul(big.NewInt(sectors), pricePerSector)
}
