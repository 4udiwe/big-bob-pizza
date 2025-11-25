package outbox

const (
	StatusPending   = "pending"
	StatusFailed    = "failed"
	StatusProcessed = "processed"
)

type Status struct {
	ID   int
	Name string
}
