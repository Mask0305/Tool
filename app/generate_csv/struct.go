package generate_csv

import "time"

type Item struct {
	ID               string    `bson:"id"`
	Info             Info      `bson:"info"`
	Catalog          Catalog   `bson:"catalog"`
	GrossProfitPrice int64     `bson:"-"`
	GrossProfitRate  float64   `bson:"-"`
	CreatedAt        time.Time `bson:"createdAt"`
}

type Catalog struct {
	// 品名
	Info CatalogInfo `bson:"info"`
	// 品牌
	Brand Brand `bson:"brand"`
	// 款式
	Style      Style `bson:"style"`
	TotalViews int64 `bson:"totalViews"`
}

type CatalogInfo struct {
	// 品名
	Name string
}

type Info struct {
	Price    Price  `bson:"price"`
	SaleMode uint32 `bson:"saleMode"`
}

type Price struct {
	Sale int64 `bson:"sale"`
	Buy  int64 `bson:"buy"`
}

type Brand struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type Style struct {
	ID string `bson:"id"`
	// 款式細項
	Name string `bson:"name"`
	// 款式大項
	Head string `bson:"head"`
}
