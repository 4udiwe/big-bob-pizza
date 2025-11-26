package outbox

import "github.com/google/uuid"

type Event struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}
