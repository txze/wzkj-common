package define

import "time"

type SubOrderDoc struct {
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UserId     string `json:"user_id"`
	ProviderId int    `json:"provider_id"`

	SubOrderNo string `json:"sub_order_no"`
	OrderNo    string `json:"order_no"`
	Status     int    `json:"status"`

	TotalAmount    int `json:"total_amount"`
	PaymentAmount  int `json:"payment_amount"`
	DiscountAmount int `json:"discount_amount"`

	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`

	PaymentTime *time.Time `json:"payment_time"`
	PaymentType string     `json:"payment_type"`
}

func (s *SubOrderDoc) IndexName() string {
	return "sub_order_query"
}
