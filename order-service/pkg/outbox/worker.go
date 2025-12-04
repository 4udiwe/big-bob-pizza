package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Worker реализует outbox-паттерн:
// периодически читает события из таблицы outbox, публикует их в Kafka
// и помечает как обработанные или неудавшиеся.
type Worker struct {
	// repo — абстракция над хранилищем outbox-событий (обычно таблица БД).
	repo Repository
	// publisher — абстракция над транспортом (Kafka-паблишер и т.п.).
	publisher Publisher
	// topic — Kafka-топик, в который будут публиковаться события.
	topic string

	// batchLimit — сколько событий максимум за один проход processBatch.
	batchLimit int
	// interval — как часто забирать новые pending-события.
	interval time.Duration

	// requeBatchLimit — сколько failed-событий пытаться перевыставить за раз.
	requeBatchLimit int
	// requeFailedInterval — период, с которым будут перевыставляться failed-события.
	requeFailedInterval time.Duration
}

// NewWorker конструирует Worker с заданным репозиторием, паблишером и настройками батчей/интервалов.
// Worker сам по себе ничего не делает, пока не будет вызван Run.
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

// Run запускает основной цикл воркера в отдельной горутине и немедленно возвращает управление.
// Внутри горутины воркер периодически:
//   - читает pending-события из outbox (processBatch);
//   - перевыставляет failed-события (requeueFailed).
//
// Остановка:
//   - когда ctx.Done() будет закрыт, цикл завершится и воркер корректно остановится.
func (w *Worker) Run(ctx context.Context) {
	go func() {
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
	}()
}

// processBatch забирает из репозитория pending-события и пытается отправить каждое в Kafka.
// Успешные события помечаются как processed, провалившиеся — как failed с текстом ошибки.
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

// requeueFailed просит репозиторий перевыставить ограниченное число failed-событий обратно в pending.
// Конкретная логика (например, увеличение retry-счётчика) реализуется в Repository.
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
