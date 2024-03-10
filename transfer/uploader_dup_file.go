package transfer

import (
	"time"

	"github.com/0glabs/0g-storage-client/common/blockchain"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const SubmitEventHash = "0x167ce04d2aa1981994d3a31695da0d785373335b1078cec239a1a3a2c7675555"

// uploadDuplicatedFile uploads file to storage node that already exists by root.
// In this case, user only need to submit transaction on blockchain, and wait for
// file finality on storage node.
func (uploader *Uploader) uploadDuplicatedFile(data core.IterableData, tags []byte, root common.Hash) error {
	// submit transaction on blockchain
	_, receipt, err := uploader.SubmitLogEntry([]core.IterableData{data}, [][]byte{tags}, true)
	if err != nil {
		return errors.WithMessage(err, "Failed to submit log entry")
	}

	// parse submission from event log
	var submission *contract.FlowSubmit
	for _, v := range receipt.Logs {
		if v.Topics[0] == common.HexToHash(SubmitEventHash) {
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

		info, err := uploader.clients[0].GetFileInfoByTxSeq(txSeq)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to get file info from storage node")
		}

		if info != nil && info.Finalized {
			return info, nil
		}
	}
}
