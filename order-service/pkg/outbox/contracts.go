package outbox

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FetchPending(ctx context.Context, limit int) ([]Event, error)
	MarkProcessed(ctx context.Context, ids []uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errorText string) error
	RequeueFailed(ctx context.Context, limit int) ([]Event, error)
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}
