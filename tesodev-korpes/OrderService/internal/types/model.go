package types

import (
	"time"
)

type Order struct {
	Id               string    `bson:"_id" json:"id"`
	CustomerName     string    `bson:"customer_name" json:"customer_name"`
	CustomerLastName string    `bson:"customer_last_name" json:"customer_last_name"`
	CustomerId       string    `bson:"customer_id" json:"customer_id"`
	OrderName        string    `bson:"order_name" json:"order_name"`
	ShipmentStatus   string    `bson:"shipment_status" json:"shipment_status"`
	PaymentMethod    string    `bson:"payment_method" json:"payment_method"`
	OrderTotal       int       `bson:"order_total" json:"order_total"`
	Price            int       `bson:"price" json:"price"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}
