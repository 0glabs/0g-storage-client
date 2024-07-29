package transfer

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sort"
	"time"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// smallFileSizeThreshold is the maximum file size to upload without log entry available on storage node.
const smallFileSizeThreshold = int64(256 * 1024)

// defaultTaskSize is the default number of data segments to upload in a single upload RPC request
const defaultTaskSize = uint(10)

var dataAlreadyExistsError = "Invalid params: root; data: already uploaded and finalized"
var segmentAlreadyExistsError = "segment has already been uploaded or is being uploaded"

func isDuplicateError(msg string) bool {
	return msg == dataAlreadyExistsError || msg == segmentAlreadyExistsError
}

// UploadOption upload option for a file
type UploadOption struct {
	Tags             []byte // transaction tags
	FinalityRequired bool   // wait for file finalized on uploaded nodes or not
	TaskSize         uint   // number of segment to upload in single rpc request
	ExpectedReplica  uint   // expected number of replications
	SkipTx           bool   // skip sending transaction on chain, this can set to true only if the data has already settled on chain before
}

// Uploader uploader to upload file to 0g storage, send on-chain transactions and transfer data to storage nodes.
type Uploader struct {
	flow    *contract.FlowContract // flow contract instance
	clients []*node.ZgsClient      // 0g storage clients
	logger  *logrus.Logger         // logger
}

func getShardConfigs(ctx context.Context, clients []*node.ZgsClient) ([]*shard.ShardConfig, error) {
	shardConfigs := make([]*shard.ShardConfig, 0)
	for _, client := range clients {
		shardConfig, err := client.GetShardConfig(ctx)
		if err != nil {
			return nil, err
		}
		if !shardConfig.IsValid() {
			return nil, errors.New("NumShard is zero")
		}
		shardConfigs = append(shardConfigs, &shardConfig)
	}
	return shardConfigs, nil
}

// NewUploader Initialize a new uploader.
func NewUploader(flow *contract.FlowContract, clients []*node.ZgsClient, opts ...zg_common.LogOption) (*Uploader, error) {
	if len(clients) == 0 {
		return nil, errors.New("storage node not specified")
	}
	uploader := &Uploader{
		clients: clients,
		logger:  zg_common.NewLogger(opts...),
		flow:    flow,
	}
	return uploader, nil
}

func (uploader *Uploader) checkLogExistance(ctx context.Context, root common.Hash) (bool, error) {
	for _, client := range uploader.clients {
		info, err := client.GetFileInfo(ctx, root)
		if err != nil {
			return false, errors.WithMessage(err, fmt.Sprintf("Failed to get file info from storage node %v", client.URL()))
		}
		// log entry available
		if info != nil {
			return true, nil
		}
	}
	return false, nil
}

// BatchUpload submit multiple data to 0g storage contract batchly in single on-chain transaction, then transfer the data to the storage nodes.
func (uploader *Uploader) BatchUpload(ctx context.Context, datas []core.IterableData, waitForLogEntry bool, option ...[]UploadOption) (common.Hash, []common.Hash, error) {
	stageTimer := time.Now()

	n := len(datas)
	if n == 0 {
		return common.Hash{}, nil, errors.New("empty datas")
	}
	var opts []UploadOption
	if len(option) > 0 {
		opts = option[0]
	} else {
		opts = make([]UploadOption, n)
	}
	if len(opts) != n {
		return common.Hash{}, nil, errors.New("datas and tags length mismatch")
	}
	uploader.logger.WithFields(logrus.Fields{
		"dataNum": n,
	}).Info("Prepare to upload batchly")

	trees := make([]*merkle.Tree, n)
	toSubmitDatas := make([]core.IterableData, 0)
	toSubmitTags := make([][]byte, 0)
	dataRoots := make([]common.Hash, n)
	var lastTreeToSubmit *merkle.Tree
	for i := 0; i < n; i++ {
		data := datas[i]
		opt := opts[i]

		uploader.logger.WithFields(logrus.Fields{
			"size":     data.Size(),
			"chunks":   data.NumChunks(),
			"segments": data.NumSegments(),
		}).Info("Data prepared to upload")

		// Calculate file merkle root.
		tree, err := core.MerkleTree(data)
		if err != nil {
			return common.Hash{}, nil, errors.WithMessage(err, "Failed to create data merkle tree")
		}
		uploader.logger.WithField("root", tree.Root()).Info("Data merkle root calculated")
		trees[i] = tree
		dataRoots[i] = trees[i].Root()

		// Check existance
		exist, err := uploader.checkLogExistance(ctx, trees[i].Root())
		if err != nil {
			return common.Hash{}, nil, errors.WithMessage(err, "Failed to check if skipped log entry available on storage node")
		}
		if !opts[i].SkipTx || !exist {
			toSubmitDatas = append(toSubmitDatas, data)
			toSubmitTags = append(toSubmitTags, opt.Tags)
		}
		lastTreeToSubmit = trees[i]
	}

	// Append log on blockchain
	var txHash common.Hash
	var receipt *types.Receipt
	if len(toSubmitDatas) > 0 {
		var err error
		if txHash, receipt, err = uploader.SubmitLogEntry(ctx, toSubmitDatas, toSubmitTags, waitForLogEntry); err != nil {
			return common.Hash{}, nil, errors.WithMessage(err, "Failed to submit log entry")
		}
		if waitForLogEntry {
			// Wait for storage node to retrieve log entry from blockchain
			if err := uploader.waitForLogEntry(ctx, lastTreeToSubmit.Root(), false, receipt); err != nil {
				return common.Hash{}, nil, errors.WithMessage(err, "Failed to check if log entry available on storage node")
			}
		}
	}

	for i := 0; i < n; i++ {
		// Upload file to storage node
		if err := uploader.uploadFile(ctx, datas[i], trees[i], opts[i].ExpectedReplica, opts[i].TaskSize); err != nil {
			return common.Hash{}, nil, errors.WithMessage(err, "Failed to upload file")
		}

		if waitForLogEntry {
			// Wait for transaction finality
			if err := uploader.waitForLogEntry(ctx, trees[i].Root(), opts[i].FinalityRequired, receipt); err != nil {
				return common.Hash{}, nil, errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
			}
		}
	}

	uploader.logger.WithField("duration", time.Since(stageTimer)).Info("batch upload took")

	return txHash, dataRoots, nil
}

// Upload submit data to 0g storage contract, then transfer the data to the storage nodes.
func (uploader *Uploader) Upload(ctx context.Context, data core.IterableData, option ...UploadOption) error {
	stageTimer := time.Now()

	var opt UploadOption
	if len(option) > 0 {
		opt = option[0]
	}

	uploader.logger.WithFields(logrus.Fields{
		"size":     data.Size(),
		"chunks":   data.NumChunks(),
		"segments": data.NumSegments(),
	}).Info("Data prepared to upload")

	// Calculate file merkle root.
	tree, err := core.MerkleTree(data)
	if err != nil {
		return errors.WithMessage(err, "Failed to create data merkle tree")
	}
	uploader.logger.WithField("root", tree.Root()).Info("Data merkle root calculated")

	// Check existance
	exist, err := uploader.checkLogExistance(ctx, tree.Root())
	if err != nil {
		return errors.WithMessage(err, "Failed to check if skipped log entry available on storage node")
	}
	// Append log on blockchain
	if !opt.SkipTx || !exist {
		var receipt *types.Receipt

		if _, receipt, err = uploader.SubmitLogEntry(ctx, []core.IterableData{data}, [][]byte{opt.Tags}, true); err != nil {
			return errors.WithMessage(err, "Failed to submit log entry")
		}

		// For small data, could upload file to storage node immediately.
		// Otherwise, need to wait for log entry available on storage node,
		// which requires transaction confirmed on blockchain.
		if data.Size() <= smallFileSizeThreshold {
			uploader.logger.Info("Upload small data immediately")
		} else {
			// Wait for storage node to retrieve log entry from blockchain
			if err = uploader.waitForLogEntry(ctx, tree.Root(), false, receipt); err != nil {
				return errors.WithMessage(err, "Failed to check if log entry available on storage node")
			}
		}
	}

	// Upload file to storage node
	if err = uploader.uploadFile(ctx, data, tree, opt.ExpectedReplica, opt.TaskSize); err != nil {
		return errors.WithMessage(err, "Failed to upload file")
	}

	// Wait for transaction finality
	if err = uploader.waitForLogEntry(ctx, tree.Root(), opt.FinalityRequired, nil); err != nil {
		return errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
	}

	uploader.logger.WithField("duration", time.Since(stageTimer)).Info("upload took")

	return nil
}

// SubmitLogEntry submit the data to 0g storage contract by sending a transaction
func (uploader *Uploader) SubmitLogEntry(ctx context.Context, datas []core.IterableData, tags [][]byte, waitForReceipt bool) (common.Hash, *types.Receipt, error) {
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
	opts, err := uploader.flow.CreateTransactOpts(ctx)
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to create opts to send transaction")
	}

	var tx *types.Transaction
	if len(datas) == 1 {
		opts.Value = submissions[0].Fee()
		tx, err = uploader.flow.Submit(opts, submissions[0])
	} else {
		opts.Value = big.NewInt(0)
		for _, v := range submissions {
			opts.Value = new(big.Int).Add(opts.Value, v.Fee())
		}
		tx, err = uploader.flow.BatchSubmit(opts, submissions)
	}
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to send transaction to append log entry")
	}

	uploader.logger.WithField("hash", tx.Hash().Hex()).Info("Succeeded to send transaction to append log entry")

	if waitForReceipt {
		// Wait for successful execution
		receipt, err := uploader.flow.WaitForReceipt(ctx, tx.Hash(), true)
		return tx.Hash(), receipt, err
	}
	return tx.Hash(), nil, nil
}

// Wait for log entry ready on storage node.
func (uploader *Uploader) waitForLogEntry(ctx context.Context, root common.Hash, finalityRequired bool, receipt *types.Receipt) error {
	uploader.logger.WithFields(logrus.Fields{
		"root":     root,
		"finality": finalityRequired,
	}).Info("Wait for log entry on storage node")

	reminder := util.NewReminder(uploader.logger, time.Minute)

	for {
		time.Sleep(time.Second)

		ok := true
		for _, client := range uploader.clients {
			info, err := client.GetFileInfo(ctx, root)
			if err != nil {
				return errors.WithMessage(err, fmt.Sprintf("Failed to get file info from storage node %v", client.URL()))
			}
			// log entry unavailable yet
			if info == nil {
				fields := logrus.Fields{}
				if receipt != nil {
					if status, err := client.GetStatus(ctx); err == nil {
						fields["txBlockNumber"] = receipt.BlockNumber
						fields["zgsNodeSyncHeight"] = status.LogSyncHeight
					}
				}

				reminder.Remind("Log entry is unavailable yet", fields)
				ok = false
				break
			}

			if finalityRequired && !info.Finalized {
				reminder.Remind("Log entry is available, but not finalized yet", logrus.Fields{
					"cached":           info.IsCached,
					"uploadedSegments": info.UploadedSegNum,
				})
				ok = false
				break
			}
		}

		if ok {
			break
		}
	}

	return nil
}

func (uploader *Uploader) newSegmentUploader(ctx context.Context, data core.IterableData, tree *merkle.Tree, expectedReplica uint, taskSize uint) (*segmentUploader, error) {
	numSegments := data.NumSegments()
	shardConfigs, err := getShardConfigs(ctx, uploader.clients)
	if err != nil {
		return nil, err
	}
	if !shard.CheckReplica(shardConfigs, expectedReplica) {
		return nil, fmt.Errorf("selected nodes cannot cover all shards")
	}
	clientTasks := make([][]*uploadTask, 0)
	for clientIndex, shardConfig := range shardConfigs {
		// skip finalized nodes
		info, _ := uploader.clients[clientIndex].GetFileInfo(ctx, tree.Root())
		if info != nil && info.Finalized {
			continue
		}
		// create upload tasks
		segIndex := shardConfig.ShardId
		tasks := make([]*uploadTask, 0)
		for ; segIndex < numSegments; segIndex += shardConfig.NumShard * uint64(taskSize) {
			tasks = append(tasks, &uploadTask{
				clientIndex: clientIndex,
				segIndex:    segIndex,
				numShard:    shardConfig.NumShard,
			})
		}
		clientTasks = append(clientTasks, tasks)
	}
	sort.SliceStable(clientTasks, func(i, j int) bool {
		return len(clientTasks[i]) > len(clientTasks[j])
	})
	tasks := make([]*uploadTask, 0)
	if len(clientTasks) > 0 {
		for taskIndex := 0; taskIndex < len(clientTasks[0]); taskIndex += 1 {
			for i := 0; i < len(clientTasks) && taskIndex < len(clientTasks[i]); i += 1 {
				tasks = append(tasks, clientTasks[i][taskIndex])
			}
		}
	}

	return &segmentUploader{
		data:     data,
		tree:     tree,
		clients:  uploader.clients,
		tasks:    tasks,
		taskSize: taskSize,
		logger:   uploader.logger,
	}, nil
}

func (uploader *Uploader) uploadFile(ctx context.Context, data core.IterableData, tree *merkle.Tree, expectedReplica uint, taskSize uint) error {
	stageTimer := time.Now()

	if taskSize == 0 {
		taskSize = defaultTaskSize
	}

	uploader.logger.WithFields(logrus.Fields{
		"segNum":  data.NumSegments(),
		"nodeNum": len(uploader.clients),
	}).Info("Begin to upload file")

	segmentUploader, err := uploader.newSegmentUploader(ctx, data, tree, expectedReplica, taskSize)
	if err != nil {
		return err
	}

	err = parallel.Serial(ctx, segmentUploader, len(segmentUploader.tasks), min(runtime.GOMAXPROCS(0), len(uploader.clients)*5), 0)
	if err != nil {
		return err
	}

	uploader.logger.WithFields(logrus.Fields{
		"duration": time.Since(stageTimer),
		"segNum":   data.NumSegments(),
	}).Info("Completed to upload file")

	return nil
}
