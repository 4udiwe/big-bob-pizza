package kafka

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	EventID    uuid.UUID       `json:"eventId"`
	EventType  string          `json:"eventType"`
	OccurredAt time.Time       `json:"occuredAt"`
	Data       json.RawMessage `json:"data"`
}
