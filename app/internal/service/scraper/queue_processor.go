package scraper

import (
	"example-wikipedia-scraper/internal/dto"
	"example-wikipedia-scraper/internal/interfaces"
	"sync"
	"time"
)

// QueueProcessor przetwarza kolejkę ofert (SRP)
type QueueProcessor struct {
	pageQueue      chan *dto.PageDTO
	batchSaver     *BatchSaver
	batchSize      int
	flushInterval  time.Duration
	processingChan chan struct{}
	logger         interfaces.LoggerInterface
	wg             sync.WaitGroup
}

func NewQueueProcessor(
	pageQueue chan *dto.PageDTO,
	batchSaver *BatchSaver,
	batchSize int,
	flushInterval time.Duration,
	logger interfaces.LoggerInterface,
) *QueueProcessor {
	return &QueueProcessor{
		pageQueue:      pageQueue,
		batchSaver:     batchSaver,
		batchSize:      batchSize,
		flushInterval:  flushInterval,
		processingChan: make(chan struct{}, 1),
		logger:         logger,
	}
}

func (qp *QueueProcessor) Start() {
	qp.wg.Add(1)
	go qp.processBatches()
}

func (qp *QueueProcessor) Stop() {
	close(qp.pageQueue)
	qp.wg.Wait()
}

func (qp *QueueProcessor) TriggerFlush() {
	select {
	case qp.processingChan <- struct{}{}:
	default:
	}
}

func (qp *QueueProcessor) processBatches() {
	defer qp.wg.Done()

	batch := make([]*dto.PageDTO, 0, qp.batchSize)
	ticker := time.NewTicker(qp.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case page, ok := <-qp.pageQueue:
			if !ok {
				qp.batchSaver.ProcessBatch(batch)
				return
			}

			batch = append(batch, page)
			if len(batch) >= qp.batchSize {
				batch = qp.batchSaver.ProcessBatch(batch)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				batch = qp.batchSaver.ProcessBatch(batch)
			}

		case <-qp.processingChan:
			batch = qp.batchSaver.ProcessBatch(batch)
		}
	}
}
