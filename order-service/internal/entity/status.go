package entity

type StatusName string

const (
	StatusCreated    StatusName = "created"
	StatusPaid       StatusName = "paid"
	StatusPrepearing StatusName = "prepearing"
	StatusDelivering StatusName = "delivering"
	StatusCompleted  StatusName = "completed"
	StatusCancelled  StatusName = "cancelled"
)

type OrderStatus struct {
	ID   int
	Name StatusName
}
