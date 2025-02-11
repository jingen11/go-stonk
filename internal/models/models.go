package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockData struct {
	MetaData        Meta                  `json:"Meta Data"`
	TimeSeriesDaily map[string]DailyPrice `json:"Time Series (Daily)"`
}

type Meta struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type DailyPrice struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Symbol struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Symbol          string             `bson:"symbol"`
	LastFetchedDate primitive.DateTime `bson:"lastFetchedDate"`
}

type Price struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	Symbol string             `bson:"symbol"`
	Date   primitive.DateTime `bson:"date"`
	Open   string             `bson:"open"`
	High   string             `bson:"high"`
	Low    string             `bson:"low"`
	Close  string             `bson:"close"`
	Volume string             `bson:"volume"`
}
