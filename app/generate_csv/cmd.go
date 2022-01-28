package generate_csv

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tyr-tech-team/hawk/config"
	"github.com/tyr-tech-team/hawk/infra/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CsvSearch(filter bson.M) []*Item {
	client, _ := mongodb.NewDial(config.Mongo{
		Host:       "mongo1:27017,mongo2:27017,mongo3:27017",
		User:       "",
		Password:   "",
		Database:   "eagle",
		ReplicaSet: "rs0",
	})

	c := client.Database("eagle").Collection("item")

	opts := options.Find()

	opts.SetSort(
		bson.D{
			bson.E{Key: "info.brand", Value: 1},
			bson.E{Key: "info.price.sale", Value: 1},
		},
	)

	raw, err := c.Find(context.Background(), filter, opts)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	data := make([]*Item, 0)

	if err := raw.All(context.Background(), &data); err != nil {
		spew.Dump(err)
		return nil
	}
	return data

}

func Generate(items []*Item) {

	bytesBuffer := &bytes.Buffer{}

	writer := csv.NewWriter(bytesBuffer)
	writer.Write([]string{"序號", "品牌", "品名", "款式", "商品ID", "銷售金額", "賣方實收金額", "毛利金額", "毛利率"})

	num := 1
	for _, v := range items {
		decCost := decimal.NewFromInt(v.Info.Price.Buy)
		decSale := decimal.NewFromInt(v.Info.Price.Sale)
		// 毛利
		decProfit := decSale.Sub(decCost)

		v.GrossProfitPrice = decProfit.IntPart()

		calRate := decSale.Mul(decimal.NewFromFloat(0.01))
		f, _ := calRate.Float64()
		if f < float64(1) {
			spew.Dump(v.ID)
		}
		decRate := decProfit.Div(calRate)
		v.GrossProfitRate, _ = decRate.Round(2).Float64()

		// CSV
		writer.Write([]string{
			cast.ToString(num),
			// 品牌
			v.Info.Brand,
			// 品名
			v.Info.Name,
			// 款式
			func() string {
				return fmt.Sprintf("%s-%s", v.Info.Style.Type, v.Info.Style.Name)
			}(),
			// 商品ID
			v.ID,
			// 銷售金額
			cast.ToString(v.Info.Price.Sale),
			// 賣方實收金額
			cast.ToString(v.Info.Price.Buy),
			// 毛利金額
			cast.ToString(v.GrossProfitPrice),
			// 毛利率
			func() string {
				return cast.ToString(v.GrossProfitRate) + "%"
			}(),
		})
		num++
	}

	writer.Flush()

	err := ioutil.WriteFile("0128上架與排程中商品資料.CSV", bytesBuffer.Bytes(), 0644)
	if err != nil {
		spew.Dump(err)
	}
}
