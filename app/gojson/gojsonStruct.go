package main

type TestStruct struct {
	Info          []TestStruct_sub2 `bson:"info" json:"info"`
	Result        string            `bson:"result" json:"result"`
	ResultMessage string            `bson:"result_message" json:"result_message"`
}

type TestStruct_sub2 struct {
	Amount           int64             `bson:"amount" json:"amount"`
	AuthorizeExpire  string            `bson:"authorize_expire" json:"authorize_expire"`
	DisburseDate     interface{}       `bson:"disburse_date" json:"disburse_date"`
	FeeType          string            `bson:"fee_type" json:"fee_type"`
	InfoCustomerJSON string            `bson:"info_customer_json" json:"info_customer_json"`
	Installment      int64             `bson:"installment" json:"installment"`
	OrderID          string            `bson:"order_id" json:"order_id"`
	ProductName      string            `bson:"product_name" json:"product_name"`
	Refundlist       []TestStruct_sub1 `bson:"refundlist" json:"refundlist"`
	ReserveDate      string            `bson:"reserve_date" json:"reserve_date"`
	SpanappID        string            `bson:"spanapp_id" json:"spanapp_id"`
	Store            interface{}       `bson:"store" json:"store"`
	TransactingDate  string            `bson:"transacting_date" json:"transacting_date"`
	TransactionState string            `bson:"transaction_state" json:"transaction_state"`
}

type TestStruct_sub1 struct {
	FinalAmount  int64  `bson:"final_amount" json:"final_amount"`
	RefundAmount int64  `bson:"refund_amount" json:"refund_amount"`
	RefundID     string `bson:"refund_id" json:"refund_id"`
	RefundTime   string `bson:"refund_time" json:"refund_time"`
}
