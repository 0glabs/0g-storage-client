package transfer

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

// smallFileSizeThreshold is the maximum file size to upload without log entry available on storage node.
const smallFileSizeThreshold = int64(256 * 1024)

type UploadOption struct {
	Tags  []byte // for kv operations
	Force bool   // for kv to upload same file
}

type Uploader struct {
	flow   *contract.FlowContract
	client *node.ZeroGStorageClient
}

func NewUploader(flow *contract.FlowContract, client *node.Client) *Uploader {
	return &Uploader{
		flow:   flow,
		client: client.ZeroGStorage(),
	}
}

func NewUploaderLight(client *node.Client) *Uploader {
	return &Uploader{
		client: client.ZeroGStorage(),
	}
}

func (uploader *Uploader) Upload(data core.IterableData, option ...UploadOption) error {
	var opt UploadOption
	if len(option) > 0 {
		opt = option[0]
	}

	logrus.WithFields(logrus.Fields{
		"size":     data.Size(),
		"chunks":   data.NumChunks(),
		"segments": data.NumSegments(),
	}).Info("Data prepared to upload")

	// Calculate file merkle root.
	tree, err := core.MerkleTree(data)
	if err != nil {
		return errors.WithMessage(err, "Failed to create data merkle tree")
	}
	logrus.WithField("root", tree.Root()).Info("Data merkle root calculated")

	info, err := uploader.client.GetFileInfo(tree.Root())
	if err != nil {
		return errors.WithMessage(err, "Failed to get data info from storage node")
	}

	logrus.WithField("info", info).Debug("Log entry retrieved from storage node")

	// In case that user interact with blockchain via Metamask
	if uploader.flow == nil && info == nil {
		return errors.New("log entry not available on storage node")
	}

	// already finalized
	if info != nil && info.Finalized {
		if !opt.Force {
			return errors.New("Data already exists on ZeroGStorage network")
		}

		// Allow to upload duplicated file for KV scenario
		if err = uploader.uploadDuplicatedFile(data, opt.Tags, tree.Root()); err != nil {
			return errors.WithMessage(err, "Failed to upload duplicated data")
		}

		return nil
	}

	// Log entry unavailable on storage node yet.
	segNum := uint64(0)
	if info == nil {
		// Append log on blockchain
		if _, err = uploader.submitLogEntry(data, opt.Tags); err != nil {
			return errors.WithMessage(err, "Failed to submit log entry")
		}

		// For small data, could upload file to storage node immediately.
		// Otherwise, need to wait for log entry available on storage node,
		// which requires transaction confirmed on blockchain.
		if data.Size() <= smallFileSizeThreshold {
			logrus.Info("Upload small data immediately")
		} else {
			// Wait for storage node to retrieve log entry from blockchain
			if err = uploader.waitForLogEntry(tree.Root(), false); err != nil {
				return errors.WithMessage(err, "Failed to check if log entry available on storage node")
			}
			info, err = uploader.client.GetFileInfo(tree.Root())
			if err != nil {
				return errors.WithMessage(err, "Failed to get file info from storage node after waitForLogEntry.")
			}
			segNum = info.UploadedSegNum
		}
	}

	// Upload file to storage node
	if err = uploader.uploadFile(data, tree, segNum); err != nil {
		return errors.WithMessage(err, "Failed to upload file")
	}

	// Wait for transaction finality
	if err = uploader.waitForLogEntry(tree.Root(), true); err != nil {
		return errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
	}

	return nil
}

func (uploader *Uploader) submitLogEntry(data core.IterableData, tags []byte) (*types.Receipt, error) {
	// Construct submission
	flow := core.NewFlow(data, tags)
	submission, err := flow.CreateSubmission()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create flow submission")
	}

	// Submit log entry to smart contract.
	opts, err := uploader.flow.CreateTransactOpts()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	tx, err := uploader.flow.Submit(opts, *submission)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to send transaction to append log entry")
	}

	logrus.WithField("hash", tx.Hash().Hex()).Info("Succeeded to send transaction to append log entry")

	// Wait for successful execution
	return uploader.flow.WaitForReceipt(tx.Hash(), true)
}

// Wait for log entry ready on storage node.
func (uploader *Uploader) waitForLogEntry(root common.Hash, finalityRequired bool) error {
	logrus.WithFields(logrus.Fields{
		"root":     root,
		"finality": finalityRequired,
	}).Info("Wait for log entry on storage node")

	for {
		time.Sleep(time.Second)

		info, err := uploader.client.GetFileInfo(root)
		if err != nil {
			return errors.WithMessage(err, "Failed to get file info from storage node")
		}

		// log entry unavailable yet
		if info == nil {
			continue
		}

		if finalityRequired && !info.Finalized {
			continue
		}

		break
	}

	return nil
}

// TODO error tolerance
func (uploader *Uploader) uploadFile(data core.IterableData, tree *merkle.Tree, segIndex uint64) error {
	logrus.WithField("segIndex", segIndex).Info("Begin to upload file")

	iter := data.Iterate(int64(segIndex*core.DefaultSegmentSize), core.DefaultSegmentSize, true)

	for {
		ok, err := iter.Next()
		if err != nil {
			return errors.WithMessage(err, "Failed to read segment")
		}

		if !ok {
			break
		}

		segment := iter.Current()
		proof := tree.ProofAt(int(segIndex))

		// Skip upload rear padding data
		numChunks := data.NumChunks()
		startIndex := segIndex * core.DefaultSegmentMaxChunks
		allDataUploaded := false
		if startIndex >= numChunks {
			// file real data already uploaded
			break
		} else if startIndex+uint64(len(segment))/core.DefaultChunkSize >= numChunks {
			// last segment has real data
			expectedLen := core.DefaultChunkSize * int(numChunks-startIndex)
			segment = segment[:expectedLen]
			allDataUploaded = true
		}

		segWithProof := node.SegmentWithProof{
			Root:     tree.Root(),
			Data:     segment,
			Index:    segIndex,
			Proof:    proof,
			FileSize: uint64(data.Size()),
		}

		if _, err = uploader.client.UploadSegment(segWithProof); err != nil {
			return errors.WithMessage(err, "Failed to upload segment")
		}

		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			chunkIndex := segIndex * core.DefaultSegmentMaxChunks
			logrus.WithFields(logrus.Fields{
				"total":      data.NumSegments(),
				"index":      segIndex,
				"chunkStart": chunkIndex,
				"chunkEnd":   chunkIndex + uint64(len(segment))/core.DefaultChunkSize,
				"root":       core.SegmentRoot(segment),
			}).Debug("Segment uploaded")
		}

		if allDataUploaded {
			break
		}

		segIndex++
	}

	logrus.Info("Completed to upload file")

	return nil
}
