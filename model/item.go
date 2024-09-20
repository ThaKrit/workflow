package model

type Item struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title    string `json:"title" gorm:"type:varchar(255);not null"`
	Amount   int    `json:"amount" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Status   string `json:"status" gorm:"type:varchar(50);not null"`
	OwnerID  uint   `json:"owner_id" gorm:"not null"`
}

type ItemStatus string

const (
	ItemPendingStatus ItemStatus = "PENDING"
	// ItemApprovedStatus ItemStatus = "APPROVED"
	// ItemRejectedStatus ItemStatus = "REJECTED"
)

type RequestFindItem struct {
	Statuses string `json:"statuses"`
}
