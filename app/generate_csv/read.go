package generate_csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func ReadCsv() []string {

	csvFile, err := os.Open("fix.csv")
	if err != nil {
		spew.Dump(err)
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	csvData := make([]*Data, 0)
	for _, line := range csvLines {
		csvData = append(csvData, &Data{
			ItemID: line[0],
		})

	}
	result := ""
	for i, v := range csvData {
		if i == 0 {
			continue
		}

		result += fmt.Sprintf(`"%v",`, v.ItemID)
		if i+1 == len(csvData) {
			result += fmt.Sprintf(`"%v"`, v.ItemID)
		}
	}

	mongoQ := fmt.Sprintf(`db.item.updateMany( { id:{ "$in":[ %v, ], } }, { "$set":{ "info.quantity":0 } })`, result)

	fmt.Println(mongoQ)

	return nil
}

type Data struct {
	ItemID string
}
