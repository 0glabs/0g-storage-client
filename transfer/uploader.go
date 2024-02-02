package transfer

import (
	"sync"
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
const routines = 20

type UploadOption struct {
	Tags     []byte // for kv operations
	Force    bool   // for kv to upload same file
	Disperse bool   // disperse files to different nodes
}

type Uploader struct {
	flow    *contract.FlowContract
	clients []*node.ZeroGStorageClient
}

func NewUploader(flow *contract.FlowContract, clients []*node.Client) *Uploader {
	uploader := NewUploaderLight(clients)
	uploader.flow = flow
	return uploader
}

func NewUploaderLight(clients []*node.Client) *Uploader {
	if len(clients) == 0 {
		panic("storage node not specified")
	}
	zgClients := make([]*node.ZeroGStorageClient, 0)
	for _, client := range clients {
		zgClients = append(zgClients, client.ZeroGStorage())
	}
	return &Uploader{
		clients: zgClients,
	}
}

// upload data(batchly in 1 blockchain transaction if there are more than one files)
func (uploader *Uploader) BatchUpload(datas []core.IterableData, waitForLogEntry bool, option ...[]UploadOption) (common.Hash, error) {
	stageTimer := time.Now()

	n := len(datas)
	if n == 0 {
		return common.Hash{}, errors.New("empty datas")
	}
	var opts []UploadOption
	if len(option) > 0 {
		opts = option[0]
	} else {
		opts = make([]UploadOption, n)
	}
	if len(opts) != n {
		return common.Hash{}, errors.New("datas and tags length mismatch")
	}
	logrus.WithFields(logrus.Fields{
		"dataNum": n,
	}).Info("Prepare to upload batchly")

	trees := make([]*merkle.Tree, n)
	toSubmitDatas := make([]core.IterableData, 0)
	toSubmitTags := make([][]byte, 0)
	var lastTreeToSubmit *merkle.Tree
	for i := 0; i < n; i++ {
		data := datas[i]
		opt := opts[i]

		logrus.WithFields(logrus.Fields{
			"size":     data.Size(),
			"chunks":   data.NumChunks(),
			"segments": data.NumSegments(),
		}).Info("Data prepared to upload")

		// Calculate file merkle root.
		tree, err := core.MerkleTree(data)
		if err != nil {
			return common.Hash{}, errors.WithMessage(err, "Failed to create data merkle tree")
		}
		logrus.WithField("root", tree.Root()).Info("Data merkle root calculated")
		trees[i] = tree

		toSubmitDatas = append(toSubmitDatas, data)
		toSubmitTags = append(toSubmitTags, opt.Tags)
		lastTreeToSubmit = trees[i]
	}

	// Append log on blockchain
	var txHash common.Hash
	if len(toSubmitDatas) > 0 {
		var err error
		if txHash, _, err = uploader.submitLogEntry(toSubmitDatas, toSubmitTags, waitForLogEntry); err != nil {
			return common.Hash{}, errors.WithMessage(err, "Failed to submit log entry")
		}
		if waitForLogEntry {
			// Wait for storage node to retrieve log entry from blockchain
			if err := uploader.waitForLogEntry(lastTreeToSubmit.Root(), false); err != nil {
				return common.Hash{}, errors.WithMessage(err, "Failed to check if log entry available on storage node")
			}
		}
	}

	for i := 0; i < n; i++ {
		// Upload file to storage node
		if err := uploader.uploadFile(datas[i], trees[i], 0, opts[i].Disperse); err != nil {
			return common.Hash{}, errors.WithMessage(err, "Failed to upload file")
		}

		if waitForLogEntry {
			// Wait for transaction finality
			if err := uploader.waitForLogEntry(trees[i].Root(), !opts[i].Disperse); err != nil {
				return common.Hash{}, errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
			}
		}
	}

	logrus.WithField("duration", time.Since(stageTimer)).Info("batch upload took")

	return txHash, nil
}

func (uploader *Uploader) Upload(data core.IterableData, option ...UploadOption) error {
	stageTimer := time.Now()

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

	info, err := uploader.clients[0].GetFileInfo(tree.Root())
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
		if _, _, err = uploader.submitLogEntry([]core.IterableData{data}, [][]byte{opt.Tags}, true); err != nil {
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
			info, err = uploader.clients[0].GetFileInfo(tree.Root())
			if err != nil {
				return errors.WithMessage(err, "Failed to get file info from storage node after waitForLogEntry.")
			}
			segNum = info.UploadedSegNum
		}
	}

	// Upload file to storage node
	if err = uploader.uploadFile(data, tree, segNum, opt.Disperse); err != nil {
		return errors.WithMessage(err, "Failed to upload file")
	}

	// Wait for transaction finality
	if err = uploader.waitForLogEntry(tree.Root(), !opt.Disperse); err != nil {
		return errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
	}

	logrus.WithField("duration", time.Since(stageTimer)).Info("upload took")

	return nil
}

func (uploader *Uploader) submitLogEntry(datas []core.IterableData, tags [][]byte, waitForReceipt bool) (common.Hash, *types.Receipt, error) {
	// Construct submission
	submissions := make([]contract.Submission, len(datas))
	for i := 0; i < len(datas); i++ {
		flow := core.NewFlow(datas[i], tags[i])
		submission, err := flow.CreateSubmission()
		if err != nil {
			return common.Hash{}, nil, errors.WithMessage(err, "Failed to create flow submission")
		}
		submissions[i] = *submission
	}

	// Submit log entry to smart contract.
	opts, err := uploader.flow.CreateTransactOpts()
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	var tx *types.Transaction
	if len(datas) == 1 {
		tx, err = uploader.flow.Submit(opts, submissions[0])
	} else {
		tx, err = uploader.flow.BatchSubmit(opts, submissions)
	}
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to send transaction to append log entry")
	}

	logrus.WithField("hash", tx.Hash().Hex()).Info("Succeeded to send transaction to append log entry")

	if waitForReceipt {
		// Wait for successful execution
		receipt, err := uploader.flow.WaitForReceipt(tx.Hash(), true)
		return tx.Hash(), receipt, err
	}
	return tx.Hash(), nil, nil
}

// Wait for log entry ready on storage node.
func (uploader *Uploader) waitForLogEntry(root common.Hash, finalityRequired bool) error {
	logrus.WithFields(logrus.Fields{
		"root":     root,
		"finality": finalityRequired,
	}).Info("Wait for log entry on storage node")

	for {
		time.Sleep(time.Second)

		info, err := uploader.clients[0].GetFileInfo(root)
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

func (uploader *Uploader) uploadSegment(tree *merkle.Tree, segIndex uint64, data core.IterableData, segment []byte, disperse bool) error {
	proof := tree.ProofAt(int(segIndex))
	dataSize := uint64(data.Size())
	numSegments := data.NumSegments()

	segWithProof := node.SegmentWithProof{
		Root:     tree.Root(),
		Data:     segment,
		Index:    segIndex,
		Proof:    proof,
		FileSize: dataSize,
	}

	if !disperse {
		if _, err := uploader.clients[0].UploadSegment(segWithProof); err != nil {
			return errors.WithMessage(err, "Failed to upload segment")
		}
	} else {
		clientIndex := segIndex % uint64(len(uploader.clients))
		ok := false
		// retry
		for i := 0; i < len(uploader.clients); i++ {
			logrus.WithFields(logrus.Fields{
				"total":       numSegments,
				"index":       segIndex,
				"clientIndex": clientIndex,
			}).Debug("Uploading segment to node..")
			if _, err := uploader.clients[clientIndex].UploadSegment(segWithProof); err != nil {
				logrus.WithFields(logrus.Fields{
					"total":       numSegments,
					"index":       segIndex,
					"clientIndex": clientIndex,
					"error":       err,
				}).Warn("Failed to upload segment to node, try next node..")
				clientIndex = (clientIndex + 1) % uint64(len(uploader.clients))
			} else {
				ok = true
				break
			}
		}
		if !ok {
			if _, err := uploader.clients[clientIndex].UploadSegment(segWithProof); err != nil {
				return errors.WithMessage(err, "Failed to upload segment")
			}
		}
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		chunkIndex := segIndex * core.DefaultSegmentMaxChunks
		logrus.WithFields(logrus.Fields{
			"total":      numSegments,
			"index":      segIndex,
			"chunkStart": chunkIndex,
			"chunkEnd":   chunkIndex + uint64(len(segment))/core.DefaultChunkSize,
			"root":       core.SegmentRoot(segment),
		}).Debug("Segment uploaded")
	}

	return nil
}

// TODO error tolerance
func (uploader *Uploader) uploadFile(data core.IterableData, tree *merkle.Tree, segIndex uint64, disperse bool) error {
	stageTimer := time.Now()

	logrus.WithFields(logrus.Fields{
		"segIndex": segIndex,
		"disperse": disperse,
		"nodeNum":  len(uploader.clients),
	}).Info("Begin to upload file")

	iter := data.Iterate(int64(segIndex*core.DefaultSegmentSize), core.DefaultSegmentSize, true)
	numChunks := data.NumChunks()

	var wg sync.WaitGroup
	errs := make(chan error, routines)

	for {
		ok, err := iter.Next()
		if err != nil {
			return errors.WithMessage(err, "Failed to read segment")
		}

		if !ok {
			break
		}

		segment := iter.Current()

		// Skip upload rear padding data

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

		if !disperse {
			err := uploader.uploadSegment(tree, segIndex, data, segment, disperse)
			if err != nil {
				return err
			}
		} else {
			wg.Add(1)
			go func(tree *merkle.Tree, segIndex uint64, data core.IterableData, segment []byte) {
				defer wg.Done()
				errs <- uploader.uploadSegment(tree, segIndex, data, segment, true)
			}(tree, segIndex, data, append([]byte(nil), segment...))
			// if segIndex%uint64(len(uploader.clients)) == uint64(len(uploader.clients))-1 || allDataUploaded {
			if segIndex%uint64(routines) == 0 || allDataUploaded {
				wg.Wait()
				close(errs)
				for e := range errs {
					if e != nil {
						return e
					}
				}
				errs = make(chan error, routines)
			}
		}

		if allDataUploaded {
			break
		}

		segIndex++
	}

	logrus.Info("Completed to upload file")

	logrus.WithField("duration", time.Since(stageTimer)).Info("upload file took")

	return nil
}
