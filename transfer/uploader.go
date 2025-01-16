package transfer

import (
	"context"
	"fmt"
	"math/big"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer/dir"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/web3go"
	"github.com/openweb3/web3go/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// defaultTaskSize is the default number of data segments to upload in a single upload RPC request
const defaultTaskSize = uint(10)
const defaultBatchSize = uint(10)

var dataAlreadyExistsError = "Invalid params: root; data: already uploaded and finalized"
var segmentAlreadyExistsError = "segment has already been uploaded or is being uploaded"
var tooManyDataError = "too many data writing"
var tooManyDataRetries = 12

func isDuplicateError(msg string) bool {
	return strings.Contains(msg, dataAlreadyExistsError) || strings.Contains(msg, segmentAlreadyExistsError)
}

func isTooManyDataError(msg string) bool {
	return strings.Contains(msg, tooManyDataError)
}

var submitLogEntryRetries = 12
var specifiedBlockError = "Specified block header does not exist"

func isRetriableSubmitLogEntryError(msg string) bool {
	return strings.Contains(msg, specifiedBlockError)
}

type FinalityRequirement uint

const (
	FileFinalized     FinalityRequirement = iota // wait for file finalization
	TransactionPacked                            // wait for transaction packed
)

// UploadOption upload option for a file
type UploadOption struct {
	Tags             []byte              // transaction tags
	FinalityRequired FinalityRequirement // finality setting
	TaskSize         uint                // number of segment to upload in single rpc request
	ExpectedReplica  uint                // expected number of replications
	SkipTx           bool                // skip sending transaction on chain, this can set to true only if the data has already settled on chain before
	Fee              *big.Int            // fee in neuron
	Nonce            *big.Int            // nonce for transaction
}

// BatchUploadOption upload option for a batching
type BatchUploadOption struct {
	Fee         *big.Int       // fee in neuron
	Nonce       *big.Int       // nonce for transaction
	TaskSize    uint           // number of files to upload simutanously
	DataOptions []UploadOption // upload option for single file, nonce and fee are ignored
}

// Uploader uploader to upload file to 0g storage, send on-chain transactions and transfer data to storage nodes.
type Uploader struct {
	flow     *contract.FlowContract // flow contract instance
	market   *contract.Market       // market contract instance
	clients  []*node.ZgsClient      // 0g storage clients
	routines int                    // number of go routines for uploading
	logger   *logrus.Logger         // logger
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
func NewUploader(ctx context.Context, w3Client *web3go.Client, clients []*node.ZgsClient, opts ...zg_common.LogOption) (*Uploader, error) {
	if len(clients) == 0 {
		return nil, errors.New("Storage node not specified")
	}

	status, err := clients[0].GetStatus(context.Background())
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to get status from storage node %v", clients[0].URL())
	}

	chainId, err := w3Client.Eth.ChainId()
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get chain ID from blockchain node")
	}

	if chainId != nil && *chainId != status.NetworkIdentity.ChainId {
		return nil, errors.Errorf("Chain ID mismatch, blockchain = %v, storage node = %v", *chainId, status.NetworkIdentity.ChainId)
	}

	flow, err := contract.NewFlowContract(status.NetworkIdentity.FlowContractAddress, w3Client)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create flow contract")
	}

	market, err := flow.GetMarketContract(ctx)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to get market contract from flow contract %v", status.NetworkIdentity.FlowContractAddress)
	}

	uploader := &Uploader{
		clients: clients,
		logger:  zg_common.NewLogger(opts...),
		flow:    flow,
		market:  market,
	}

	return uploader, nil
}

func checkLogExistence(ctx context.Context, clients []*node.ZgsClient, root common.Hash) (*node.FileInfo, error) {
	var info *node.FileInfo
	var err error
	for _, client := range clients {
		info, err = client.GetFileInfo(ctx, root)
		if err != nil {
			return nil, err
		}
		// log entry available
		if info != nil {
			return info, nil
		}
	}
	return info, nil
}

func (uploader *Uploader) WithRoutines(routines int) *Uploader {
	uploader.routines = routines
	return uploader
}

// SplitableUpload submit data to 0g storage contract and large data will be splited to reduce padding cost.
func (uploader *Uploader) SplitableUpload(ctx context.Context, data core.IterableData, fragmentSize int64, option ...UploadOption) ([]common.Hash, []common.Hash, error) {
	if fragmentSize < core.DefaultChunkSize {
		fragmentSize = core.DefaultChunkSize
	}
	// align size of fragment to 2 power
	fragmentSize = int64(core.NextPow2(uint64(fragmentSize)))
	uploader.logger.Infof("fragment size: %v", fragmentSize)

	txHashes := make([]common.Hash, 0)
	rootHashes := make([]common.Hash, 0)
	if data.Size() <= fragmentSize {
		txHash, rootHash, err := uploader.Upload(ctx, data, option...)
		if err != nil {
			return txHashes, rootHashes, err
		}
		txHashes = append(txHashes, txHash)
		rootHashes = append(rootHashes, rootHash)
	} else {
		fragments := data.Split(fragmentSize)
		uploader.logger.Infof("splitted origin file into %v fragments, %v bytes each.", len(fragments), fragmentSize)
		var opt UploadOption
		if len(option) > 0 {
			opt = option[0]
		}
		for l := 0; l < len(fragments); l += int(defaultBatchSize) {
			r := min(l+int(defaultBatchSize), len(fragments))
			uploader.logger.Infof("batch submitting fragments %v to %v...", l, r)
			opts := BatchUploadOption{
				Fee:         nil,
				Nonce:       nil,
				DataOptions: make([]UploadOption, 0),
			}
			for i := l; i < r; i += 1 {
				opts.DataOptions = append(opts.DataOptions, opt)
			}
			txHash, roots, err := uploader.BatchUpload(ctx, fragments[l:r], opts)
			if err != nil {
				return txHashes, rootHashes, err
			}
			txHashes = append(txHashes, txHash)
			rootHashes = append(rootHashes, roots...)
		}
	}
	return txHashes, rootHashes, nil
}

// BatchUpload submit multiple data to 0g storage contract batchly in single on-chain transaction, then transfer the data to the storage nodes.
// The nonce for upload transaction will be the first non-nil nonce in given upload options, the protocol fee is the sum of fees in upload options.
func (uploader *Uploader) BatchUpload(ctx context.Context, datas []core.IterableData, option ...BatchUploadOption) (common.Hash, []common.Hash, error) {
	stageTimer := time.Now()

	n := len(datas)
	if n == 0 {
		return common.Hash{}, nil, errors.New("empty datas")
	}
	var opts BatchUploadOption
	if len(option) > 0 {
		opts = option[0]
	} else {
		opts = BatchUploadOption{
			Fee:         nil,
			Nonce:       nil,
			DataOptions: make([]UploadOption, n),
		}
	}
	opts.TaskSize = max(opts.TaskSize, 1)
	if len(opts.DataOptions) != n {
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

	var wg sync.WaitGroup
	errs := make(chan error, opts.TaskSize)
	fileInfos := make([]*node.FileInfo, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			data := datas[i]
			uploader.logger.WithFields(logrus.Fields{
				"size":     data.Size(),
				"chunks":   data.NumChunks(),
				"segments": data.NumSegments(),
			}).Info("Data prepared to upload")

			// Calculate file merkle root.
			tree, err := core.MerkleTree(data)
			if err != nil {
				errs <- errors.WithMessage(err, "Failed to create data merkle tree")
				return
			}
			uploader.logger.WithField("root", tree.Root()).Info("Data merkle root calculated")
			trees[i] = tree
			dataRoots[i] = trees[i].Root()

			// Check existence
			info, err := checkLogExistence(ctx, uploader.clients, trees[i].Root())
			if err != nil {
				errs <- errors.WithMessage(err, "Failed to check if skipped log entry available on storage node")
				return
			}
			fileInfos[i] = info
		}(i)
		if (i+1)%int(opts.TaskSize) == 0 || i == n-1 {
			wg.Wait()
			close(errs)
			for e := range errs {
				if e != nil {
					return common.Hash{}, nil, e
				}
			}
			errs = make(chan error, opts.TaskSize)
		}
	}
	for i := 0; i < n; i += 1 {
		opt := opts.DataOptions[i]
		if !opt.SkipTx || fileInfos[i] == nil {
			toSubmitDatas = append(toSubmitDatas, datas[i])
			toSubmitTags = append(toSubmitTags, opt.Tags)
			lastTreeToSubmit = trees[i]
		}
	}

	// Append log on blockchain
	var txHash common.Hash
	var receipt *types.Receipt
	if len(toSubmitDatas) > 0 {
		var err error
		if txHash, receipt, err = uploader.SubmitLogEntry(ctx, toSubmitDatas, toSubmitTags, opts.Nonce, opts.Fee); err != nil {
			return txHash, nil, errors.WithMessage(err, "Failed to submit log entry")
		}
		// Wait for storage node to retrieve log entry from blockchain
		if _, err := uploader.waitForLogEntry(ctx, lastTreeToSubmit.Root(), TransactionPacked, receipt); err != nil {
			return txHash, nil, errors.WithMessage(err, "Failed to check if log entry available on storage node")
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			info := fileInfos[i]
			if info == nil {
				var err error
				info, err = uploader.waitForLogEntry(ctx, trees[i].Root(), TransactionPacked, receipt)
				if err != nil {
					errs <- errors.WithMessage(err, "Failed to get file info from storage node")
					return
				}
			}
			// Upload file to storage node
			if err := uploader.uploadFile(ctx, info, datas[i], trees[i], opts.DataOptions[i].ExpectedReplica, opts.DataOptions[i].TaskSize); err != nil {
				errs <- errors.WithMessage(err, "Failed to upload file")
				return
			}

			// Wait for transaction finality
			if _, err := uploader.waitForLogEntry(ctx, trees[i].Root(), opts.DataOptions[i].FinalityRequired, receipt); err != nil {
				errs <- errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
				return
			}
			errs <- nil
		}(i)
		if (i+1)%int(opts.TaskSize) == 0 || i == n-1 {
			wg.Wait()
			close(errs)
			for e := range errs {
				if e != nil {
					return txHash, nil, e
				}
			}
			errs = make(chan error, opts.TaskSize)
		}
	}

	uploader.logger.WithField("duration", time.Since(stageTimer)).Info("batch upload took")

	return txHash, dataRoots, nil
}

// Upload submit data to 0g storage contract, then transfer the data to the storage nodes.
// returns the submission transaction hash and the hash will be zero if transaction is skipped.
func (uploader *Uploader) Upload(ctx context.Context, data core.IterableData, option ...UploadOption) (common.Hash, common.Hash, error) {
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
		return common.Hash{}, common.Hash{}, errors.WithMessage(err, "Failed to create data merkle tree")
	}
	uploader.logger.WithField("root", tree.Root()).Info("Data merkle root calculated")

	// Check existence
	info, err := checkLogExistence(ctx, uploader.clients, tree.Root())
	if err != nil {
		return common.Hash{}, tree.Root(), errors.WithMessage(err, "Failed to check if skipped log entry available on storage node")
	}
	txHash := common.Hash{}
	// Append log on blockchain
	if !opt.SkipTx || info == nil {
		var receipt *types.Receipt

		txHash, receipt, err = uploader.SubmitLogEntry(ctx, []core.IterableData{data}, [][]byte{opt.Tags}, opt.Nonce, opt.Fee)
		if err != nil {
			return txHash, tree.Root(), errors.WithMessage(err, "Failed to submit log entry")
		}

		// Wait for storage node to retrieve log entry from blockchain
		info, err = uploader.waitForLogEntry(ctx, tree.Root(), TransactionPacked, receipt)
		if err != nil {
			return txHash, tree.Root(), errors.WithMessage(err, "Failed to check if log entry available on storage node")
		}
	}
	// Upload file to storage node
	if err := uploader.uploadFile(ctx, info, data, tree, opt.ExpectedReplica, opt.TaskSize); err != nil {
		return txHash, tree.Root(), errors.WithMessage(err, "Failed to upload file")
	}

	// Wait for transaction finality
	if _, err = uploader.waitForLogEntry(ctx, tree.Root(), opt.FinalityRequired, nil); err != nil {
		return txHash, tree.Root(), errors.WithMessage(err, "Failed to wait for transaction finality on storage node")
	}

	uploader.logger.WithField("duration", time.Since(stageTimer)).Info("upload took")

	return txHash, tree.Root(), nil
}

func (uploader *Uploader) UploadDir(ctx context.Context, folder string, option ...UploadOption) (txnHash, rootHash common.Hash, _ error) {
	// Build the file tree representation of the directory.
	root, err := dir.BuildFileTree(folder)
	if err != nil {
		return txnHash, rootHash, errors.WithMessage(err, "failed to build file tree")
	}

	tdata, err := root.MarshalBinary()
	if err != nil {
		return txnHash, rootHash, errors.WithMessage(err, "failed to encode file tree")
	}

	// Create an in-memory data object from the encoded file tree.
	iterdata, err := core.NewDataInMemory(tdata)
	if err != nil {
		return txnHash, rootHash, errors.WithMessage(err, "failed to create `IterableData` in memory")
	}

	// Generate the Merkle tree from the in-memory data.
	mtree, err := core.MerkleTree(iterdata)
	if err != nil {
		return txnHash, rootHash, errors.WithMessage(err, "failed to create merkle tree")
	}
	rootHash = mtree.Root()

	// Flattening the file tree to get the list of files and their relative paths.
	_, relPaths := root.Flatten(func(n *dir.FsNode) bool {
		return n.Type == dir.FileTypeFile && n.Size > 0
	})

	logrus.Infof("Total %d files to be uploaded", len(relPaths))

	// Upload each file to the storage network.
	for i := range relPaths {
		path := filepath.Join(folder, relPaths[i])
		txhash, _, err := uploader.UploadFile(ctx, path, option...)
		if err != nil {
			return txnHash, rootHash, errors.WithMessagef(err, "failed to upload file %s", path)
		}

		logrus.WithFields(logrus.Fields{
			"txnHash": txhash,
			"path":    path,
		}).Info("File uploaded successfully")
	}

	// Finally, upload the directory metadata
	txnHash, _, err = uploader.Upload(ctx, iterdata, option...)
	if err != nil {
		err = errors.WithMessage(err, "failed to upload directory metadata")
	}

	return txnHash, rootHash, err
}

func (uploader *Uploader) UploadFile(ctx context.Context, path string, option ...UploadOption) (txnHash common.Hash, rootHash common.Hash, err error) {
	file, err := core.Open(path)
	if err != nil {
		err = errors.WithMessagef(err, "failed to open file %s", path)
		return
	}
	defer file.Close()

	return uploader.Upload(ctx, file, option...)
}

// SubmitLogEntry submit the data to 0g storage contract by sending a transaction
func (uploader *Uploader) SubmitLogEntry(ctx context.Context, datas []core.IterableData, tags [][]byte, nonce *big.Int, fee *big.Int) (common.Hash, *types.Receipt, error) {
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
	if nonce != nil {
		opts.Nonce = nonce
	}

	var tx *types.Transaction
	pricePerSector, err := uploader.market.PricePerSector(&bind.CallOpts{Context: ctx})
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to read price per sector")
	}
	if len(datas) == 1 {
		if fee != nil {
			opts.Value = fee
		} else {
			opts.Value = submissions[0].Fee(pricePerSector)
		}
		uploader.logger.WithField("fee(neuron)", opts.Value).Info("submit with fee")
		for attempt := 0; attempt < submitLogEntryRetries; attempt++ {
			tx, err = uploader.flow.Submit(opts, submissions[0])
			if err == nil || !isRetriableSubmitLogEntryError(err.Error()) || attempt >= submitLogEntryRetries-1 {
				break
			}
			uploader.logger.WithFields(logrus.Fields{
				"error":   err,
				"attempt": attempt,
			}).Warn("Failed to submit, retrying...")
			time.Sleep(10 * time.Second)
		}
	} else {
		if fee != nil {
			opts.Value = fee
		} else {
			opts.Value = big.NewInt(0)
			for _, v := range submissions {
				opts.Value = new(big.Int).Add(opts.Value, v.Fee(pricePerSector))
			}
		}
		uploader.logger.WithField("fee(neuron)", opts.Value).Info("batch submit with fee")
		for attempt := 0; attempt < submitLogEntryRetries; attempt++ {
			tx, err = uploader.flow.BatchSubmit(opts, submissions)
			if err == nil || !isRetriableSubmitLogEntryError(err.Error()) || attempt >= submitLogEntryRetries-1 {
				break
			}
			uploader.logger.WithFields(logrus.Fields{
				"error":   err,
				"attempt": attempt,
			}).Warn("Failed to submit, retrying...")
			time.Sleep(10 * time.Second)
		}
	}
	if err != nil {
		return common.Hash{}, nil, errors.WithMessage(err, "Failed to send transaction to append log entry")
	}

	uploader.logger.WithField("hash", tx.Hash().Hex()).Info("Succeeded to send transaction to append log entry")

	// Wait for successful execution
	receipt, err := uploader.flow.WaitForReceipt(ctx, tx.Hash(), true)
	return tx.Hash(), receipt, err
}

// Wait for log entry ready on storage node.
func (uploader *Uploader) waitForLogEntry(ctx context.Context, root common.Hash, finalityRequired FinalityRequirement, receipt *types.Receipt) (*node.FileInfo, error) {
	uploader.logger.WithFields(logrus.Fields{
		"root":     root,
		"finality": finalityRequired,
	}).Info("Wait for log entry on storage node")

	reminder := util.NewReminder(uploader.logger, time.Minute)

	var info *node.FileInfo
	var err error

	for {
		time.Sleep(time.Second)

		ok := true
		for _, client := range uploader.clients {
			info, err = client.GetFileInfo(ctx, root)
			if err != nil {
				return nil, err
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

			if finalityRequired <= FileFinalized && !info.Finalized {
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

	return info, nil
}

func (uploader *Uploader) newSegmentUploader(ctx context.Context, info *node.FileInfo, data core.IterableData, tree *merkle.Tree, expectedReplica uint, taskSize uint) (*segmentUploader, error) {
	shardConfigs, err := getShardConfigs(ctx, uploader.clients)
	if err != nil {
		return nil, err
	}
	if !shard.CheckReplica(shardConfigs, expectedReplica) {
		return nil, fmt.Errorf("selected nodes cannot cover all shards")
	}
	// compute index in flow
	startSegmentIndex, endSegmentIndex := core.SegmentRange(info.Tx.StartEntryIndex, info.Tx.Size)
	clientTasks := make([][]*uploadTask, 0)
	for clientIndex, shardConfig := range shardConfigs {
		// skip finalized nodes
		info, _ := uploader.clients[clientIndex].GetFileInfo(ctx, tree.Root())
		if info != nil && info.Finalized {
			continue
		}
		// create upload tasks
		// segIndex % NumShard = shardId (in flow)
		segIndex := shardConfig.NextSegmentIndex(startSegmentIndex)
		tasks := make([]*uploadTask, 0)
		for ; segIndex <= endSegmentIndex; segIndex += shardConfig.NumShard * uint64(taskSize) {
			tasks = append(tasks, &uploadTask{
				clientIndex: clientIndex,
				segIndex:    segIndex - startSegmentIndex,
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
		txSeq:    info.Tx.Seq,
		clients:  uploader.clients,
		tasks:    tasks,
		taskSize: taskSize,
		logger:   uploader.logger,
	}, nil
}

func (uploader *Uploader) uploadFile(ctx context.Context, info *node.FileInfo, data core.IterableData, tree *merkle.Tree, expectedReplica uint, taskSize uint) error {
	stageTimer := time.Now()

	if taskSize == 0 {
		taskSize = defaultTaskSize
	}

	uploader.logger.WithFields(logrus.Fields{
		"segNum":  data.NumSegments(),
		"nodeNum": len(uploader.clients),
	}).Info("Begin to upload file")

	segmentUploader, err := uploader.newSegmentUploader(ctx, info, data, tree, expectedReplica, taskSize)
	if err != nil {
		return err
	}

	opt := parallel.SerialOption{
		Routines: uploader.routines,
	}
	err = parallel.Serial(ctx, segmentUploader, len(segmentUploader.tasks), opt)
	if err != nil {
		return err
	}

	uploader.logger.WithFields(logrus.Fields{
		"duration": time.Since(stageTimer),
		"segNum":   data.NumSegments(),
	}).Info("Completed to upload file")

	return nil
}

// FileSegmentsWithProof wraps segments with proof and file info
type FileSegmentsWithProof struct {
	*node.FileInfo
	Segments []node.SegmentWithProof
}

type FileSegmentUploader struct {
	clients []*node.ZgsClient // 0g storage clients
	logger  *logrus.Logger    // logger
}

func NewFileSegementUploader(clients []*node.ZgsClient, opts ...zg_common.LogOption) *FileSegmentUploader {
	return &FileSegmentUploader{
		clients: clients,
		logger:  zg_common.NewLogger(opts...),
	}
}

// Upload uploads file segments with proof to the storage nodes parallelly.
// Note: only `ExpectedReplica` and `TaskSize` are used from UploadOption.
func (uploader *FileSegmentUploader) Upload(ctx context.Context, fileSeg FileSegmentsWithProof, option ...UploadOption) error {
	var opt UploadOption
	if len(option) > 0 {
		opt = option[0]
	}

	if opt.TaskSize == 0 {
		opt.TaskSize = defaultTaskSize
	}

	stageTimer := time.Now()
	if uploader.logger.IsLevelEnabled(logrus.DebugLevel) {
		uploader.logger.WithFields(logrus.Fields{
			"uploadOption": opt,
			"fileSegments": fileSeg,
		}).Debug("Begin to upload file segments with proof")
	}

	fsUploader, err := uploader.newFileSegmentUploader(ctx, fileSeg, opt.ExpectedReplica, opt.TaskSize)
	if err != nil {
		return err
	}

	sopt := parallel.SerialOption{
		Routines: min(runtime.GOMAXPROCS(0), len(uploader.clients)*5),
	}
	err = parallel.Serial(ctx, fsUploader, len(fsUploader.tasks), sopt)
	if err != nil {
		return err
	}

	if uploader.logger.IsLevelEnabled(logrus.DebugLevel) {
		uploader.logger.WithFields(logrus.Fields{
			"total":    len(fileSeg.Segments),
			"duration": time.Since(stageTimer),
		}).Debug("Completed to upload file segments with proof")
	}

	return nil
}

func (uploader *FileSegmentUploader) newFileSegmentUploader(
	ctx context.Context, fileSeg FileSegmentsWithProof, expectedReplica uint, taskSize uint) (*fileSegmentUploader, error) {

	//  get shard configurations
	shardConfigs, err := getShardConfigs(ctx, uploader.clients)
	if err != nil {
		return nil, err
	}

	// validate replica requirements
	if !shard.CheckReplica(shardConfigs, expectedReplica) {
		return nil, fmt.Errorf("selected nodes cannot cover all shards")
	}

	// create upload tasks for each segment
	clientTasks := make([][]*uploadTask, len(uploader.clients))
	for i, segment := range fileSeg.Segments {
		startSegmentIndex, endSegmentIndex := core.SegmentRange(fileSeg.Tx.StartEntryIndex, fileSeg.Tx.Size)
		segmentIndex := startSegmentIndex + segment.Index

		if segmentIndex > endSegmentIndex {
			return nil, errors.New("segment index out of range")
		}

		// assign segment to shard configurations
		for clientIndex, shardConfig := range shardConfigs {
			// skip nodes that do not cover the segment
			if !shardConfig.HasSegment(segmentIndex) {
				continue
			}

			// skip finalized nodes
			nodeInfo, _ := uploader.clients[clientIndex].GetFileInfo(ctx, segment.Root)
			if nodeInfo != nil && nodeInfo.Finalized {
				continue
			}

			// add task for the client to upload the segment
			clientTasks[clientIndex] = append(clientTasks[clientIndex], &uploadTask{
				clientIndex: clientIndex,
				segIndex:    uint64(i),
			})
		}
	}
	sort.SliceStable(clientTasks, func(i, j int) bool {
		return len(clientTasks[i]) > len(clientTasks[j])
	})

	// group tasks by task size
	uploadTasks := make([][]*uploadTask, 0, len(clientTasks))
	for _, tasks := range clientTasks {
		// split tasks into batches of taskSize
		for len(tasks) > int(taskSize) {
			uploadTasks = append(uploadTasks, tasks[:taskSize])
			tasks = tasks[taskSize:]
		}

		if len(tasks) > 0 {
			uploadTasks = append(uploadTasks, tasks)
		}
	}
	util.Shuffle(uploadTasks)

	return &fileSegmentUploader{
		FileSegmentsWithProof: fileSeg,
		clients:               uploader.clients,
		tasks:                 uploadTasks,
		logger:                uploader.logger,
	}, nil
}
