package datamodels

import "time"

type Order struct {
	ID         int       `json:"id" sql:"id"`
	UserID     int       `json:"user_id" sql:"user_id"`
	SellerID   int       `json:"seller_id" sql:"seller_id"`
	ProductID  int       `json:"product_id" sql:"product_id"`
	OrderNum   int       `json:"order_num" sql:"order_num"`
	TotalPrice float64   `json:"total_price" sql:"total_price"`
	Status     int       `json:"status" sql:"status"`
	CreateTime time.Time `json:"create_time" sql:"create_time"`
}

const (
	OrderWait = iota
	// 1
	OrderSuccess
	// 2
	OrderFailed
)
