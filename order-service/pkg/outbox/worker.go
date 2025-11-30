package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	repo      Repository
	publisher Publisher
	topic     string

	batchLimit int
	interval   time.Duration

	requeBatchLimit     int
	requeFailedInterval time.Duration
}

func NewWorker(repo Repository, publisher Publisher, topic string, batchLimit, requeBatchLimit int, interval, requeInterval time.Duration) *Worker {
	return &Worker{
		repo:                repo,
		publisher:           publisher,
		topic:               topic,
		batchLimit:          batchLimit,
		requeBatchLimit:     requeBatchLimit,
		interval:            interval,
		requeFailedInterval: requeInterval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	requeTicker := time.NewTicker(w.requeFailedInterval)
	defer requeTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Info("OutboxWorker: shutting down")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		case <-requeTicker.C:
			w.requeueFailed(ctx)
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
	events, err := w.repo.FetchPending(ctx, w.batchLimit)
	if err != nil {
		logrus.Errorf("OutboxWorker: failed to fetch pending events: %v", err)
		return
	}

	if len(events) == 0 {
		logrus.Debug("OutboxWorker: no pending events")
		return
	}

	processedIDs := make([]uuid.UUID, 0, len(events))

	for _, e := range events {

		err := w.publisher.Publish(ctx, w.topic, e.EventType, e.Payload)
		if err != nil {
			logrus.Errorf("OutboxWorker: failed to publish event %v: %v", e.ID, err)
			if errMark := w.repo.MarkFailed(ctx, e.ID, fmt.Sprintf("%v", err)); errMark != nil {
				logrus.Errorf("OutboxWorker: failed to mark event %v as failed: %v", e.ID, errMark)
			}
			continue
		}

		logrus.Infof("OutboxWorker: successfully published event %v", e.ID)
		processedIDs = append(processedIDs, e.ID)
	}

	if len(processedIDs) > 0 {
		if err := w.repo.MarkProcessed(ctx, processedIDs); err != nil {
			logrus.Errorf("OutboxWorker: failed to mark events as processed: %v", err)
		}
	}
}

func (w *Worker) requeueFailed(ctx context.Context) {
	events, err := w.repo.RequeueFailed(ctx, w.requeBatchLimit)
	if err != nil {
		logrus.Errorf("OutboxWorker: failed to requeue failed events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	logrus.Infof("OutboxWorker: requeued %d failed events", len(events))
}
