package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/openweb3/web3go"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
)

type FlowContract struct {
	*blockchain.Contract
	*Flow
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

	return &FlowContract{contract, flow}, nil
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
