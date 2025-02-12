package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type StockData struct {
	Status     string  `json:"status"`
	From       string  `json:"from"`
	Symbol     string  `json:"symbol"`
	Open       float64 `json:"open"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Close      float64 `json:"close"`
	Volume     float64 `json:"volume"`
	AfterHours float64 `json:"afterHours"`
	PreMarket  float64 `json:"preMarket"`
}

type Symbol struct {
	Id              primitive.ObjectID `bson:"_id,omitempty"`
	Symbol          string             `bson:"symbol"`
	LastFetchedDate primitive.DateTime `bson:"lastFetchedDate"`
}

type Price struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Symbol     string             `bson:"symbol"`
	Date       primitive.DateTime `bson:"date"`
	Open       float64            `bson:"open"`
	High       float64            `bson:"high"`
	Low        float64            `bson:"low"`
	Close      float64            `bson:"close"`
	Volume     float64            `bson:"volume"`
	AfterHours float64            `bson:"afterHours"`
	PreMarket  float64            `bson:"preMarket"`
	HAOpen     float64            `bson:"haOpen,omitempty"`
	HAClose    float64            `bson:"haClose,omitempty"`
}
