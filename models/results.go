package models

type Results struct {
	ID          int64    `bson:"_id"` // owner's Telegram ID
	LastUpdated int64    `bson:"LastUpdated"`
	List        []Result `bson:"List"`
}

type Result struct {
	Subject string `bson:"Subject"`
	Points  int    `bson:"Points"`
	Mark    string `bson:"Mark"`
}
