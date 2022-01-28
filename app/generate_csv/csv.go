package generate_csv

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GenerateCSV() {

	end, _ := time.Parse(time.RFC3339, "2022-01-28T10:00:00+08:00")

	filter := bson.M{

		//"createdAt": bson.M{
		//"$lte": end,
		//},

		"status": bson.M{
			"$in": []int32{1, 5},
		},

		"onSoldAt": bson.M{"$lt": end},
	}

	itemList := CsvSearch(filter)

	Generate(itemList)

}
