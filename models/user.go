package models

type User struct {
	ID         int64  `bson:"_id"` // Telegram ID
	Lastname   string `bson:"Lastname"`
	Name       string `bson:"Name"`
	SecondName string `bson:"SecondName"`
	DocNumber  string `bson:"DocNumber"`
}
