package routes

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type BalanceData struct {
	Id      bson.ObjectId `bson:"_id"`
	Account string        `bson:"account"`
	Usdt    string        `bson:"usdt"`
	Token   string        `bson:"token"`
	Date    time.Time     `bson:"date"`
}

type VolumeData struct {
	Id      bson.ObjectId `bson:"_id"`
	Account string        `bson:"account"`
	Volume  string        `bson:"volume"`
	Date    time.Time     `bson:"date"`
}
type OrderData struct {
	Id      bson.ObjectId `bson:"_id"`
	OrderID string        `bson:"order_id"`
	Time    string        `bson:"time"`
}
