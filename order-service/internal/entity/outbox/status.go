package outbox

type StatusName string

const (
	StatusPending   StatusName = "pending"
	StatusFailed    StatusName = "failed"
	StatusProcessed StatusName = "processed"
)

type Status struct {
	ID   int
	Name StatusName
}
