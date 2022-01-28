package generate_csv

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
