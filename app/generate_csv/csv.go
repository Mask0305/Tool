package generate_csv

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateCSV() {

	end, _ := time.Parse(time.RFC3339, "2022-05-25T19:01:00+08:00")

	filter := bson.M{

		"createdAt": bson.M{
			"$lte": end,
		},

		// 上架、排程上架
		//"status": bson.M{
		//"$in": []int32{1, 5},
		//},

		// 有庫存
		"info.quantity": bson.M{
			"$gt": 0,
		},

		"catalogId": "",

		//"onSoldAt": bson.M{"$lt": end},
	}

	itemList := CsvSearch(filter)

	Generate(itemList)

}
