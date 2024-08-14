package batcher

import (
	"context"
	"errors"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/gammazero/workerpool"
)

type Orchestrator struct {
	Config
	TimeoutConfig

	Queue          disperser.BlobStore
	MinibatchStore MinibatchStore
	Dispatcher     disperser.Dispatcher
	EncoderClient  disperser.EncoderClient

	ChainState            core.IndexedChainState
	AssignmentCoordinator core.AssignmentCoordinator
	Aggregator            core.SignatureAggregator
	EncodingStreamer      *EncodingStreamer
	Transactor            core.Transactor
	TransactionManager    TxnManager
	Metrics               *Metrics
	HeartbeatChan         chan time.Time

	ethClient common.EthClient
	finalizer Finalizer
	logger    logging.Logger

	MiniBatcher    *Minibatcher
	BatchConfirmer *BatchConfirmer
}

func NewOrchestrator(
	config Config,
	timeoutConfig TimeoutConfig,
	queue disperser.BlobStore,
	minibatchStore MinibatchStore,
	dispatcher disperser.Dispatcher,
	chainState core.IndexedChainState,
	assignmentCoordinator core.AssignmentCoordinator,
	encoderClient disperser.EncoderClient,
	aggregator core.SignatureAggregator,
	ethClient common.EthClient,
	finalizer Finalizer,
	transactor core.Transactor,
	txnManager TxnManager,
	logger logging.Logger,
	metrics *Metrics,
	heartbeatChan chan time.Time,
) (*Orchestrator, error) {
	batchTrigger := NewEncodedSizeNotifier(
		make(chan struct{}, 1),
		uint64(config.BatchSizeMBLimit)*1024*1024, // convert to bytes
	)
	streamerConfig := StreamerConfig{
		SRSOrder:                 config.SRSOrder,
		EncodingRequestTimeout:   config.PullInterval,
		EncodingQueueLimit:       config.EncodingRequestQueueSize,
		TargetNumChunks:          config.TargetNumChunks,
		MaxBlobsToFetchFromStore: config.MaxBlobsToFetchFromStore,
		FinalizationBlockDelay:   config.FinalizationBlockDelay,
		ChainStateTimeout:        timeoutConfig.ChainStateTimeout,
	}
	encodingWorkerPool := workerpool.New(config.NumConnections)
	encodingStreamer, err := NewEncodingStreamer(streamerConfig, queue, chainState, encoderClient, assignmentCoordinator, batchTrigger, encodingWorkerPool, metrics.EncodingStreamerMetrics, metrics, logger)
	if err != nil {
		return nil, err
	}

	miniBatcher, err := NewMinibatcher(config.MinibatcherConfig, queue, minibatchStore, dispatcher, encodingStreamer, encodingWorkerPool, logger, metrics)
	if err != nil {
		return nil, err
	}

	batchConfirmerConfig := BatchConfirmerConfig{
		PullInterval:                 config.PullInterval,
		DispersalTimeout:             timeoutConfig.DispersalTimeout,
		DispersalStatusCheckInterval: config.DispersalStatusCheckInterval,
		AttestationTimeout:           timeoutConfig.AttestationTimeout,
		SRSOrder:                     config.SRSOrder,
		NumConnections:               config.MinibatcherConfig.MaxNumConnections,
		MaxNumRetriesPerBlob:         config.MinibatcherConfig.MaxNumRetriesPerBlob,
	}
	batchConfirmer, err := NewBatchConfirmer(batchConfirmerConfig, queue, minibatchStore, dispatcher, chainState, assignmentCoordinator, encodingStreamer, aggregator, ethClient, transactor, txnManager, miniBatcher, logger, metrics)
	if err != nil {
		return nil, err
	}

	return &Orchestrator{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		Queue:          queue,
		MinibatchStore: minibatchStore,
		Dispatcher:     dispatcher,
		EncoderClient:  encoderClient,

		ChainState:            chainState,
		AssignmentCoordinator: assignmentCoordinator,
		Aggregator:            aggregator,
		EncodingStreamer:      encodingStreamer,
		Transactor:            transactor,
		TransactionManager:    txnManager,
		Metrics:               metrics,

		ethClient:      ethClient,
		finalizer:      finalizer,
		logger:         logger.With("component", "Orchestrator"),
		HeartbeatChan:  heartbeatChan,
		MiniBatcher:    miniBatcher,
		BatchConfirmer: batchConfirmer,
	}, nil
}

func (o *Orchestrator) Start(ctx context.Context) error {
	err := o.ChainState.Start(ctx)
	if err != nil {
		return err
	}

	// Wait for few seconds for indexer to index blockchain
	// This won't be needed when we switch to using Graph node
	time.Sleep(indexerWarmupDelay)
	err = o.EncodingStreamer.Start(ctx)
	if err != nil {
		return err
	}
	batchTrigger := o.EncodingStreamer.EncodedSizeNotifier

	err = o.BatchConfirmer.Start(ctx)
	if err != nil {
		return err
	}

	go func() {
		receiptChan := o.TransactionManager.ReceiptChan()
		for {
			select {
			case <-ctx.Done():
				return
			case receiptOrErr := <-receiptChan:
				o.logger.Info("received response from transaction manager", "receipt", receiptOrErr.Receipt, "err", receiptOrErr.Err)
				err := o.BatchConfirmer.ProcessConfirmedBatch(ctx, receiptOrErr)
				if err != nil {
					o.logger.Error("failed to process confirmed batch", "err", err)
				}
			}
		}
	}()
	o.TransactionManager.Start(ctx)

	o.finalizer.Start(ctx)

	go func() {
		ticker := time.NewTicker(o.PullInterval)
		defer ticker.Stop()
		cancelFuncs := make([]context.CancelFunc, 0)
		for {
			select {
			case <-ctx.Done():
				for _, cancel := range cancelFuncs {
					cancel()
				}
				return
			case <-ticker.C:
				o.signalLiveness()
				cancel, err := o.MiniBatcher.HandleSingleMinibatch(ctx)
				if err != nil {
					if errors.Is(err, errNoEncodedResults) {
						o.logger.Warn("no encoded results to construct minibatch")
					} else {
						o.logger.Error("failed to process minibatch", "err", err)
					}
				}
				if cancel != nil {
					cancelFuncs = append(cancelFuncs, cancel)
				}
			case <-batchTrigger.Notify:
				ticker.Stop()
				o.signalLiveness()
				cancel, err := o.MiniBatcher.HandleSingleMinibatch(ctx)
				if err != nil {
					if errors.Is(err, errNoEncodedResults) {
						o.logger.Warn("no encoded results to construct minibatch")
					} else {
						o.logger.Error("failed to process minibatch", "err", err)
					}
				}
				if cancel != nil {
					cancelFuncs = append(cancelFuncs, cancel)
				}
				ticker.Reset(o.PullInterval)
			}
		}
	}()
	return nil
}

func (o *Orchestrator) signalLiveness() {
	select {
	case o.HeartbeatChan <- time.Now():
		o.logger.Info("Heartbeat signal sent")
	default:
		// This case happens if there's no receiver ready to consume the heartbeat signal.
		// It prevents the goroutine from blocking if the channel is full or not being listened to.
		o.logger.Warn("Heartbeat signal skipped, no receiver on the channel")
	}
}
