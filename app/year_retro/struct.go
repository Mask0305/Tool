package year_retro

import "time"

type BuyOrder struct {
	Memo string `bson:"memo"`
}

type Item struct {
	ID               string  `bson:"id"`
	Info             Info    `bson:"info"`
	GrossProfitPrice int64   `bson:"-"`
	GrossProfitRate  float64 `bson:"-"`
}

type Price struct {
	Sale int64 `bson:"sale"`
	Buy  int64 `bson:"buy"`
}

type Info struct {
	Name  string `bson:"name"`
	Brand string `bson:"brand"`
	Price Price  `bson:"price"`
	Style Style  `bson:"style"`
}

type Style struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
}

type Delivery struct {
	No      string  `bson:"no"`
	Address Address `bson:"address"`
	Type    string  `bson:"type"`
	Price   int64   `bson:"-"`
}

type Address struct {
	City  string `bson:"city"`
	Addrs string `bson:"addrs"`
}

type Sellorder struct {
	Buyer       Buyer       `bson:"buyer"`
	Logistics   Logistics   `bson:"logistics"`
	OrderDetail OrderDetail `bson:"orderDetail"`
}

type OrderDetail struct {
	OrderTotalAmount int64 `bson:"orderTotalAmount"`
}

type Buyer struct {
	No   string `bson:"no"`
	Name string `bson:"name"`
}

type Member struct {
	No       string    `bson:"no"`
	Birthday time.Time `bson:"birthday"`
	IDcard   IDcard    `bson:"idCard"`
}

type IDcard struct {
	No string `bson:"no"`
}

type Logistics struct {
	Type    LogisticsType `bson:"type"`
	Address Address       `bson:"address"`
}

// LogisticsType - 配送方式
type LogisticsType int32

const (
	// LOGISTICSTYPE_NONE  - Default
	LOGISTICSTYPE_NONE LogisticsType = iota

	// LOGISTICSTYPE_HOME - 宅配到府
	LOGISTICSTYPE_HOME

	// LOGISTICSTYPE_STORE - 超商取貨
	LOGISTICSTYPE_STORE

	// LOGISTICSTYPE_PERSON - 面交
	LOGISTICSTYPE_PERSON
)
