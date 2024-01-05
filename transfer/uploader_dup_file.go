package transfer

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/blockchain"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

var submissionEventHash = common.HexToHash("0x398e4f14f8588468d3654c03dc3f266e5af46083542d34db23fb04953067194b")

// uploadDuplicatedFile uploads file to storage node that already exists by root.
// In this case, user only need to submit transaction on blockchain, and wait for
// file finality on storage node.
func (uploader *Uploader) uploadDuplicatedFile(data core.IterableData, tags []byte, root common.Hash) error {
	// submit transaction on blockchain
	receipt, err := uploader.submitLogEntry(data, tags)
	if err != nil {
		return errors.WithMessage(err, "Failed to submit log entry")
	}

	// parse submission from event log
	var submission *contract.FlowSubmit
	for _, v := range receipt.Logs {
		if v.Topics[0] == submissionEventHash {
			log := blockchain.ConvertToGethLog(v)

			if submission, err = uploader.flow.ParseSubmit(*log); err != nil {
				return err
			}

			break
		}
	}

	// wait for finality from storage node
	txSeq := submission.SubmissionIndex.Uint64()
	info, err := uploader.waitForFileFinalityByTxSeq(txSeq)
	if err != nil {
		return errors.WithMessagef(err, "Failed to wait for finality for tx %v", txSeq)
	}

	if info.Tx.DataMerkleRoot != root {
		return errors.New("Merkle root mismatch, maybe transaction reverted")
	}

	return nil
}

func (uploader *Uploader) waitForFileFinalityByTxSeq(txSeq uint64) (*node.FileInfo, error) {
	logrus.WithField("txSeq", txSeq).Info("Wait for finality on storage node")

	for {
		time.Sleep(time.Second)

		info, err := uploader.client.GetFileInfoByTxSeq(txSeq)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to get file info from storage node")
		}

		if info != nil && info.Finalized {
			return info, nil
		}
	}
}
