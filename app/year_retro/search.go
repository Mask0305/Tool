package year_retro

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/tyr-tech-team/hawk/config"
	"github.com/tyr-tech-team/hawk/infra/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BuyOrderSearch(filter bson.M) []*BuyOrder {
	client, _ := mongodb.NewDial(config.Mongo{
		Host:       "mongo1:27017,mongo2:27017,mongo3:27017",
		User:       "",
		Password:   "",
		Database:   "eagle",
		ReplicaSet: "rs0",
	})

	c := client.Database("eagle").Collection("buyOrder")

	opts := options.Find()

	//opts.SetSort(bson.M{"info.brand": 1})

	raw, err := c.Find(context.Background(), filter, opts)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	data := make([]*BuyOrder, 0)

	if err := raw.All(context.Background(), &data); err != nil {
		spew.Dump(err)
		return nil
	}
	return data

}

func MemberDeliverySearch(filter bson.M) []*Delivery {
	client, _ := mongodb.NewDial(config.Mongo{
		Host:       "mongo1:27017,mongo2:27017,mongo3:27017",
		User:       "",
		Password:   "",
		Database:   "eagle",
		ReplicaSet: "rs0",
	})

	c := client.Database("eagle").Collection("memberDelivery")

	opts := options.Find()

	//opts.SetSort(bson.M{"info.brand": 1})

	raw, err := c.Find(context.Background(), filter, opts)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	data := make([]*Delivery, 0)

	if err := raw.All(context.Background(), &data); err != nil {
		spew.Dump(err)
		return nil
	}
	return data

}

func MemberSearch(filter bson.M) []*Member {
	client, _ := mongodb.NewDial(config.Mongo{
		Host:       "mongo1:27017,mongo2:27027,mongo3:27037",
		User:       "",
		Password:   "",
		Database:   "eagle",
		ReplicaSet: "rs0",
	})

	c := client.Database("eagle").Collection("member")

	opts := options.Find()

	//opts.SetSort(bson.M{"info.brand": 1})

	raw, err := c.Find(context.Background(), filter, opts)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	data := make([]*Member, 0)

	if err := raw.All(context.Background(), &data); err != nil {
		spew.Dump(err)
		return nil
	}
	return data

}
func SellOrderSearch(filter bson.M) []*Delivery {
	client, _ := mongodb.NewDial(config.Mongo{
		Host:       "mongo1:27017,mongo2:27017,mongo3:27017",
		User:       "",
		Password:   "",
		Database:   "eagle",
		ReplicaSet: "rs0",
	})

	c := client.Database("eagle").Collection("sellorder")

	opts := options.Find()

	//opts.SetSort(bson.M{"info.brand": 1})

	raw, err := c.Find(context.Background(), filter, opts)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	data := make([]*Sellorder, 0)

	if err := raw.All(context.Background(), &data); err != nil {
		spew.Dump(err)
		return nil
	}
	result := make([]*Delivery, len(data))

	for i, v := range data {
		result[i] = &Delivery{
			No: v.Buyer.No,
			Type: func() string {
				switch v.Logistics.Type {
				case 1:
					return "宅配"
				case 2:
					return "超商取貨"
				case 3:
					return "面交"
				}
				return ""
			}(),
			Address: v.Logistics.Address,
			Price:   v.OrderDetail.OrderTotalAmount,
		}
	}

	return result

}
